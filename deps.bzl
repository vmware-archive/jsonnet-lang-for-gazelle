# jsonnet dependencies. Needed by the com_github_google_go_jsonnet rule defined in repositories.bzl.
load("@com_github_google_go_jsonnet//bazel:deps.bzl", "jsonnet_go_dependencies")

def jsonnet_gazelle_dependencies():
    jsonnet_go_dependencies()
