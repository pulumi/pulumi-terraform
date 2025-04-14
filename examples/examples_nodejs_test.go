// Copyright 2016-2025, Pulumi Corporation.  All rights reserved.
//go:build nodejs || all
// +build nodejs all

package examples

import (
	"testing"
)

func TestNodejs(t *testing.T) {
	t.Parallel()
	LanguageTests(t, "nodejs")
}
