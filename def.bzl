"""Jsonnet library. It is a wrapper on top of filegroup"""


def jsonnet_library(name, srcs, visibility, deps = []): native.filegroup(
    name = name,
    srcs = native.glob(srcs) + deps,
    visibility = visibility,
)

# Usage:


# WORKSPACE
#
# load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
#
# git_repository(
#     name = "jsonnet_gazelle",
#     commit = "master",
#     remote = "https://github.com/jotadrilo/jsonnet-gazelle",
# )


# BUILD.bazel
#
# load("@bazel_gazelle//:def.bzl", "DEFAULT_LANGUAGES", "gazelle", "gazelle_binary")
#
# gazelle_binary(
#     name = "gazelle_jsonnet_binary",
#     languages = DEFAULT_LANGUAGES + [
#         "@jsonnet_gazelle//language/jsonnet:go_default_library",
#     ],
#     visibility = ["//visibility:public"],
# )
#
# # gazelle:jsonnet_allowed_imports json,groovy,txt
# # gazelle:jsonnet_native_imports libsonnet
#
# gazelle(
#     name = "gazelle_jsonnet",
#     gazelle = "//:gazelle_jsonnet_binary",
# )
