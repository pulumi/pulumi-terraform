// Copyright 2016-2025, Pulumi Corporation.  All rights reserved.
//go:build java || all
// +build java all

package examples

import (
	"testing"
)

func TestJava(t *testing.T) {
	t.Parallel()
	LanguageTests(t, "java")
}
