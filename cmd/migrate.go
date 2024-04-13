// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.megpoid.dev/go-skel/pkg/cfg"
	"go.megpoid.dev/go-skel/pkg/migration"
	"go.megpoid.dev/go-skel/pkg/sql"
	"megpoid.dev/go/contact-form/config"
	"megpoid.dev/go/contact-form/db"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Apply the database migrations to the database`,
	PreRun: func(cmd *cobra.Command, _ []string) {
		cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetBool("debug") {
			// Setup logger
			handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
			slog.SetDefault(slog.New(handler))
		} else {
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
		}

		databaseSettings := config.DatabaseSettings{}
		if err := cfg.ReadConfig(&databaseSettings); err != nil {
			return fmt.Errorf("failed to read database settings: %w", err)
		}

		migrationSettings := config.MigrationSettings{}
		if err := cfg.ReadConfig(&migrationSettings); err != nil {
			return fmt.Errorf("failed to read migration settings: %w", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		quit := make(chan os.Signal, 1)

		// Database initialization
		pool, err := sql.NewConnection(sql.Config{
			DataSourceName:  databaseSettings.DataSourceName,
			MaxIdleConns:    databaseSettings.MaxIdleConns,
			MaxOpenConns:    databaseSettings.MaxOpenConns,
			ConnMaxLifetime: databaseSettings.ConnMaxLifetime,
			ConnMaxIdleTime: databaseSettings.ConnMaxIdleTime,
			QueryLimit:      databaseSettings.QueryLimit,
		})
		if err != nil {
			return err
		}
		defer pool.Close()

		var migrationErr error

		migrationConfig := migration.Options{
			TableName: "app_migrations",
			Redo:      migrationSettings.Redo,
			Reset:     migrationSettings.Reset,
			Rollback:  migrationSettings.Rollback,
			Step:      migrationSettings.Step,
			MigrationAsset: migration.AssetOptions{
				FS:   db.Assets(),
				Root: "migrations",
			},
		}

		go func() {
			defer func() {
				quit <- os.Interrupt
			}()

			migrationErr = migration.RunMigrations(ctx, pool, migrationConfig)
			if migrationErr != nil {
				slog.Error("migration failed", "error", migrationErr)
				return
			}

			if migrationSettings.Seed {
				seedAssets := migration.AssetOptions{
					FS:   db.Seeds(),
					Root: "seed",
				}
				migrationErr = migration.ApplySQLFiles(ctx, sql.NewPgxPool(pool), seedAssets)
				if migrationErr != nil {
					slog.Error("migration failed", "error", migrationErr)
					return
				}
			}
		}()

		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		cancel()

		return migrationErr
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	databaseFlags := config.LoadDatabaseFlags(migrateCmd.Name())
	migrateFlags := config.LoadMigrateFlags(migrateCmd.Name())

	migrateCmd.Flags().AddFlagSet(databaseFlags)
	migrateCmd.Flags().AddFlagSet(migrateFlags)
}
