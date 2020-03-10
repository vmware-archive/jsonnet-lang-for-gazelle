#!/usr/bin/env bash

set -euo pipefail

GCLOUD_AUTH=${GCLOUD_AUTH:-}
FLAGS=${FLAGS:-}

if [ ! -z "${GCLOUD_AUTH}" ]; then
  export GOOGLE_APPLICATION_CREDENTIALS=/tmp/gcreds.json
  echo "${GCLOUD_AUTH}" | base64 -d >"${GOOGLE_APPLICATION_CREDENTIALS}"
  FLAGS="--google_credentials=${GOOGLE_APPLICATION_CREDENTIALS}"
fi

bazel test //... ${FLAGS} "$@" \
    --experimental_inmemory_jdeps_files \
    --experimental_inmemory_dotd_files \
    --experimental_remote_download_outputs=minimal
