Architecture of Gazelle Jsonnet
===============================

.. All external links are here.

.. _buildifier: https://github.com/bazelbuild/buildtools/tree/master/buildifier
.. _full list of directives: README.rst#Directives
.. _running gazelle with bazel: README.rst#RunningGazelleWithBazel

.. Actual content is below

Gazelle is a tool that was designed to generate and update Bazel build files for Go projects.
It is intended to simplify the maintenance of Bazel Go projects as much as possible.

This document describes how Gazelle for Jsonnet projects works.

.. contents::

Overview
--------

Gazelle for Jsonnet projects generates and updates build files according the algorithm outlined
below. Each of the steps here is described in more detail in the sections below.

* Build a configuration from command line arguments and special comments
  in the top-level build file. See Configuration_.

* For each directory in the repository:

  * Read the build file if one is present.

  * If the build file should be updated (based on configuration):

    * Scan the native jsonnet files and collect metadata needed to generate rules
      for the directory. See `Scanning native files`_.

    * Generate new rules from the build metadata collected earlier. See
      `Generating rules`_.

    * Merge the new rules into the directory's build file. Delete any rules
      which are now empty. See `Merging and deleting rules`_.

  * Add the library rules in the directory's build file to a global table,
    indexed by import path.

Configuration
-------------

Gazelle for Jsonnet stores configuration information in ``jsonnetConfig`` objects. These objects
contain settings that affect the behavior of the program.
For example:

* A list of allowed imports (non-native extensions, likely for ``importstr`` directives).
* A list of native imports (``.jsonnet``, ``.libsonnet``).
* A list of folders to ignore (the ``BUILD.bazel`` won't be modified).

``jsonnetConfig`` objects apply to individual directories. Each directory inherits
the ``jsonnetConfig`` from its parent. Values in a ``jsonnetConfig`` may be modified within
a directory using *directives* written in the directory's build file. A
directive is a special comment formatted like this:

::

  # gazelle:key value

Here are a few examples. See the `full list of directives`_.

* ``# gazelle:jsonnet_ignore_folders`` - sets a list of folders to ignore.

Scanning native files
---------------------

The information needed to render a jsonnet file is encoded in the file content.
A jsonnet file might import other jsonnet files and any kind of text file
using ``importstr`` directives.

Therefore, jsonnet files are quite flexible. This tool will not take arbitrary
imports into account but they can be defined using gazelle directives. See `Generating rules`_.

Generating rules
----------------

Once build metadata has been extracted from the sources in a directory,
Gazelle generates rules for using those sources.

We may generate the following rules:

* ``jsonnet_library`` are generated for each of the jsonnet files found.
* ``filegroup`` are generated for each of the non-native files that have been
  allowed as so. See `Configuration`_.

Rules are named using word characters only. Non-word characters are replaced with ``_``.

Example
^^^^^^^

::

    foo-bar.k.jsonnet => foo_bar_k

At this point, Gazelle does not have enough information to generate expressions
``deps`` attributes in ``jsonnet_library``. We only have a ``FilePath`` map for
each imported file extracted from the jsonnet files. These imports are stored
temporarily in a special ``_gazelle_imports`` private attribute in each rule.
Later, the imports are converted to Bazel labels and this attribute is replaced
with ``deps``.


Building and running Gazelle
----------------------------

It is recommend to run Gazelle through Bazel so all developers use the same
version set in the ``WORKSPACE`` file.

Developers should add a gazelle rule to include the ``jsonnet`` language within the
gazelle binary in the ``WORKSPACE`` file. See `running gazelle with bazel`_.

This is the most convenient way to run Gazelle, and it's what we recommend to
users.
