// Copyright 2016-2025, Pulumi Corporation.  All rights reserved.
//go:build python || all
// +build python all

package examples

import (
	"testing"
)

func TestPython(t *testing.T) {
	t.Parallel()
	LanguageTests(t, "python")
}
