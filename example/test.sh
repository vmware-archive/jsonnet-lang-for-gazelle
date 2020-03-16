#!/bin/bash

bazel run //:gazelle
bazel query "filter('demo', kind('jsonnet_to_json', '//...'))" | xargs bazel build --keep_going
