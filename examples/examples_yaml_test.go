// Copyright 2016-2025, Pulumi Corporation.  All rights reserved.
//go:build yaml || all
// +build yaml all

package examples

import (
	"testing"
)

func TestYAML(t *testing.T) {
	t.Parallel()
	LanguageTests(t, "yaml")
}
