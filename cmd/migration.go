// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	matchFirstCap = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	matchAllCap   = regexp.MustCompile(`([a-z\d])([A-Z])`)
)

var templateContent = `-- +migrate Up

-- +migrate StatementBegin
-- Use this block to create functions
-- +migrate StatementEnd

-- Migrate Down
`

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// migrationCmd represents the migrate command
var migrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "Create a database migration file",
	Long:  `Create a database migration file`,
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, _ []string) {
		cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		timestamp := time.Now().Format("20060102150405")
		name := toSnakeCase(args[0])

		var baseDir string
		if viper.GetBool("seed") {
			baseDir = "db/seed"
		} else {
			baseDir = "db/migrations"
		}

		filename := path.Join(baseDir, timestamp+"_"+name+".sql")
		err := os.WriteFile(filename, []byte(templateContent), 0o644)
		if err != nil {
			return fmt.Errorf("failed to create migration file: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrationCmd)

	migrationCmd.Flags().Bool("seed", false, "Create a seed file")
}
