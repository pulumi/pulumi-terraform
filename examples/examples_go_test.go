// Copyright 2016-2025, Pulumi Corporation.  All rights reserved.
//go:build go || all
// +build go all

package examples

import (
	"testing"
)

func TestGo(t *testing.T) {
	t.Parallel()
	LanguageTests(t, "go")
}
