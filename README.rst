Gazelle Jsonnet
===============

.. All external links are here
.. _Developer Certificate of Origin: https://cla.vmware.com/dco
.. _CONTRIBUTING.md: CONTRIBUTING.md
.. _Apache 2 license: LICENSE.txt

.. role:: direc(code)
.. role:: value(code)
.. End of directives

This implements the ``jsonnet`` language for Gazelle.

Setup
-----

Running Gazelle with Bazel
~~~~~~~~~~~~~~~~~~~~~~~~~~

To use Gazelle in a new project, add the ``bazel_gazelle`` repository and its
dependencies to your WORKSPACE file and call ``gazelle_dependencies``. It
should look like this:

.. code:: bzl

    load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

    http_archive(
        name = "io_bazel_rules_go",
        urls = [
            "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/rules_go/releases/download/v0.20.1/rules_go-v0.20.1.tar.gz",
            "https://github.com/bazelbuild/rules_go/releases/download/v0.20.1/rules_go-v0.20.1.tar.gz",
        ],
        sha256 = "842ec0e6b4fbfdd3de6150b61af92901eeb73681fd4d185746644c338f51d4c0",
    )

    http_archive(
        name = "bazel_gazelle",
        urls = [
            "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.19.0/bazel-gazelle-v0.19.0.tar.gz",
            "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.19.0/bazel-gazelle-v0.19.0.tar.gz",
        ],
        sha256 = "41bff2a0b32b02f20c227d234aa25ef3783998e5453f7eade929704dcff7cd4b",
    )

    load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

    go_rules_dependencies()

    go_register_toolchains()

    load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

    gazelle_dependencies()

    load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

    git_repository(
        name = "jsonnet_gazelle",
        branch = "master",
        remote = "https://github.com/bitnami/jsonnet-gazelle",
    )

Add the code below to the BUILD or BUILD.bazel file in the root directory of
your repository to build a gazelle binary including the Jsonnet language.

.. code:: bzl

    load("@bazel_gazelle//:def.bzl", "DEFAULT_LANGUAGES", "gazelle", "gazelle_binary")

    gazelle_binary(
        name = "gazelle_jsonnet_binary",
        languages = DEFAULT_LANGUAGES + [
            "@jsonnet_gazelle//language/jsonnet:go_default_library",
        ],
        visibility = ["//visibility:public"],
    )

    load("@io_bazel_rules_jsonnet//jsonnet:jsonnet.bzl", "jsonnet_library")

    gazelle(
        name = "gazelle",
        gazelle = "//:gazelle_jsonnet_binary",
    )

After adding this code, you can run Gazelle with Bazel.

.. code::

  $ bazel run //:gazelle

This will generate new BUILD.bazel files for your project. You can run the same
command in the future to update existing BUILD.bazel files to include new source
files.

You can pass additional arguments to Gazelle after a ``--`` argument.

.. code::

  $ bazel run //:gazelle -- -jsonnet_ignore_folders=scripts

Directives
~~~~~~~~~~

Gazelle can be configured with *directives*, which are written as top-level
comments in build files. Most options that can be set on the command line
can also be set using directives. Some options can only be set with
directives.

Directive comments have the form ``# gazelle:key value``.

Example
^^^^^^^

.. code:: bzl

  load("@io_bazel_rules_jsonnet//jsonnet:jsonnet.bzl", "jsonnet_library")

  # gazelle:jsonnet_ignore_folders scripts

  gazelle(
      name = "gazelle_jsonnet",
      gazelle = "//:gazelle_jsonnet_binary",
  )

Directives apply in the directory where they are set *and* in subdirectories.
This means, for example, if you set ``# gazelle:jsonnet_ignore_folders`` in the build file
in your project's root directory, it affects your whole project. If you
set it in a subdirectory, it only affects rules in that subtree.

The following directives are recognized:

+-----------------------------------------------------+--------------------------------------+
| **Directive**                                       | **Default value**                    |
+=====================================================+======================================+
| :direc:`# gazelle:jsonnet_ignore_folders`           | none                                 |
+-----------------------------------------------------+--------------------------------------+
| Comma-separated list of folders that should not be processed. If not specified, Gazelle    |
| will process all the folders.                                                              |
+-----------------------------------------------------+--------------------------------------+

Contributing
------------

The jsonnet-lang-for-gazelle project team welcomes contributions from the community. Before you start working with jsonnet-lang-for-gazelle, please
read our `Developer Certificate of Origin`_. All contributions to this repository must be
signed as described on that page. Your signature certifies that you wrote the patch or have the right to pass it on
as an open-source patch. For more detailed information, refer to `CONTRIBUTING.md`_.

License
-------

jsonnet-lang-for-gazelle is available under the `Apache 2 license`_.
