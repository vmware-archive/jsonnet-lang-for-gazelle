// Copyright 2019 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package jsonnet

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func (*Lang) Fix(c *config.Config, f *rule.File) {}
