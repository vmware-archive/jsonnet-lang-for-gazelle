FROM l.gcr.io/google/bazel:2.2.0

LABEL "name"="bazel action"
LABEL "maintainer"="Marko Mikulicic <mkm@bitnami.com>"
LABEL "version"="1.0.3"

LABEL "com.github.actions.icon"="code"
LABEL "com.github.actions.color"="green"
LABEL "com.github.actions.name"="Bazel Tools"
LABEL "com.github.actions.description"="This is an Action to run go commands."

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
