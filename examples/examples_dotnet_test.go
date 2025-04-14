// Copyright 2016-2025, Pulumi Corporation.  All rights reserved.
//go:build dotnet || all
// +build dotnet all

package examples

import (
	"testing"
)

func TestDotnet(t *testing.T) {
	t.Parallel()
	LanguageTests(t, "dotnet")
}
