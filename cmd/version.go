// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"megpoid.dev/go/contact-form/version"

	"github.com/spf13/cobra"
)

const versionFormatter = "GoApp version: %s, commit: %s, date: %s, clean build: %t\n"

func printVersion() {
	fmt.Printf(versionFormatter, version.Tag, version.Revision, version.LastCommit, !version.Modified)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version",
	Long:  `Prints the version and build info`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
