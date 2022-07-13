// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"megpoid.dev/go/contact-form/config"
	"megpoid.dev/go/contact-form/store/sqlstore"
	"os"
	"os/signal"
	"syscall"
)

func unmarshalFunc(val any) error {
	return viper.Unmarshal(val, unmarshalDecoder)
}

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Apply the database migrations to the database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(config.WithUnmarshal(unmarshalFunc))
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		quit := make(chan os.Signal, 1)

		// Database initialization
		conn, err := sqlstore.NewConnection(cfg.SqlSettings)
		if err != nil {
			return err
		}
		defer conn.Close()

		go func() {
			err := sqlstore.RunMigrations(ctx, conn, cfg)
			if err != nil {
				log.Println(err.Error())
			}
			quit <- os.Interrupt
		}()

		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		cancel()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().Bool("rollback", false, "Rollback last migration")
	migrateCmd.Flags().Bool("redo", false, "Rollback last migration then migrate again")
	migrateCmd.Flags().Bool("reset", false, "Drop all tables and run migration")
	migrateCmd.Flags().Int("step", 1, "Steps to rollback/redo")
	migrateCmd.Flags().String("dsn", "", "Database connection string. Setting the DSN ignores the db-* settings")
	migrateCmd.Flags().String("driver", "postgres", "Database driver")
	err := viper.BindPFlags(migrateCmd.Flags())
	cobra.CheckErr(err)
}
