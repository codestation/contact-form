// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cmd

import (
	"log/slog"

	"go.megpoid.dev/go-skel/version"

	"github.com/spf13/cobra"
)

func printVersion() {
	slog.Info("contact-form",
		slog.String("version", version.Tag),
		slog.String("commit", version.Revision),
		slog.Time("date", version.LastCommit),
		slog.Bool("clean_build", !version.Modified),
	)
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
