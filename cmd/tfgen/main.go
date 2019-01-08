// Copyright 2016-2018, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/util/cmdutil"
	"github.com/pulumi/pulumi/pkg/util/contract"
	"github.com/spf13/cobra"

	"github.com/pulumi/pulumi-terraform/pkg/tfbridge"
	"github.com/pulumi/pulumi-terraform/pkg/tfgen"
)

func main() {
	var logToStderr bool
	var quiet bool
	var verbose int
	var pkg string
	var version string
	var options tfgen.GenerateOptions
	cmd := &cobra.Command{
		Use:   os.Args[0] + " <LANGUAGE>",
		Args:  cmdutil.SpecificArgs([]string{"language"}),
		Short: "The Pulumi TFGen compiler generates Pulumi package metadata from a Terraform provider",
		Long: "The Pulumi TFGen compiler generates Pulumi package metadata from a Terraform provider.\n" +
			"\n" +
			"The tool will load the provider from your $PATH, inspect its contents dynamically,\n" +
			"and generate all of the Pulumi metadata necessary to consume the resources.\n" +
			"\n" +
			"<LANGUAGE> indicates which language/runtime to target; the current supported set of\n" +
			"languages is " + fmt.Sprintf("%v", tfgen.SupportedLanguages) + ".\n" +
			"\n" +
			"Note that there is no custom Pulumi provider code required, because the generated\n" +
			"provider plugin is metadata-driven and thus works against all Terraform providers.\n",
		Run: cmdutil.RunFunc(func(cmd *cobra.Command, args []string) error {
			var info *tfbridge.MarshallableProviderInfo
			if err := json.NewDecoder(os.Stdin).Decode(&info); err != nil {
				return errors.Wrap(err, "could not decode provider schema")
			}

			return tfgen.Generate(args[0], pkg, version, *info.Unmarshal(), &options)
		}),
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			glog.Flush()
		},
	}

	cmd.PersistentFlags().BoolVar(
		&logToStderr, "logtostderr", false, "Log to stderr instead of to files")
	cmd.PersistentFlags().StringVarP(
		&options.OutputDir, "out", "o", "", "Save generated package metadata to this directory")
	cmd.PersistentFlags().StringVar(
		&options.OverlaysDir, "overlays", "", "Use the target directory for overlays rather than the default of overlays/")
	cmd.PersistentFlags().StringVar(
		&pkg, "package", "", "The name of the generated package")
	cmd.PersistentFlags().BoolVarP(
		&quiet, "quiet", "q", false, "Suppress non-error output progress messages")
	cmd.PersistentFlags().IntVarP(
		&verbose, "verbose", "v", 0, "Enable verbose logging (e.g., v=3); anything >3 is very verbose")
	cmd.PersistentFlags().StringVar(
		&version, "version", "", "The version of the generated package")

	if err := cmd.Execute(); err != nil {
		_, fmterr := fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		contract.IgnoreError(fmterr)
		os.Exit(-1)
	}
}
