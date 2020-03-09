#!/usr/bin/env bash
#
# `go mod tidy` finds all the packages transitively imported by packages in
# your module. It adds new module requirements for packages not provided by
# any known module, and it removes requirements on modules that don't provide
# any imported packages. If a module provides packages that are only imported
# by projects that haven't migrated to modules yet, the module requirement
# will be marked with an // indirect comment. It is always good practice to
# run go mod tidy before committing a go.mod file to version control.
#
# `bazel run //:gazelle -- update-repos` vendors go.mod file into the bazel
# go_repositories macro in `repositories.bzl`.

go mod tidy
bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro='repositories.bzl%go_repositories'
