// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.megpoid.dev/go-skel/pkg/cfg"
	"megpoid.dev/go/contact-form/app"
	"megpoid.dev/go/contact-form/config"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start service",
	Long:  `Starts the HTTP endpoint and other services`,
	PreRun: func(cmd *cobra.Command, _ []string) {
		cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServer()
	},
}

func runServer() error {
	if viper.GetBool("debug") {
		// Setup logger
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
		slog.SetDefault(slog.New(handler))
	} else {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	}

	// show version on console
	printVersion()

	appConfig := app.Config{}
	if err := cfg.ReadConfig(&appConfig.General); err != nil {
		return fmt.Errorf("failed to read general config: %w", err)
	}

	if err := cfg.ReadConfig(&appConfig.Server); err != nil {
		return fmt.Errorf("failed to read server config: %w", err)
	}

	if err := cfg.ReadConfig(&appConfig.Database); err != nil {
		return fmt.Errorf("failed to read database config: %w", err)
	}

	if err := cfg.ReadConfig(&appConfig.SMTP); err != nil {
		return fmt.Errorf("failed to read smtp config: %w", err)
	}

	if err := cfg.ReadConfig(&appConfig.Captcha); err != nil {
		return fmt.Errorf("failed to read captcha config: %w", err)
	}

	// setup channel to check when app is stopped
	quit := make(chan os.Signal, 1)

	// Create a new newApp
	newApp, err := app.NewApp(appConfig)
	if err != nil {
		return fmt.Errorf("cannot create newApp: %w", err)
	}
	defer newApp.Shutdown()

	// Start newApp
	if err := newApp.Start(); err != nil {
		return fmt.Errorf("cannot start newApp: %w", err)
	}

	// Wait for kill signal before attempting to gracefully stop the running service
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	return nil
}

func init() {
	rootCmd.AddCommand(serveCmd)

	generalFs := config.LoadGeneralFlags(serveCmd.Name())
	serverFs := config.LoadServerFlags(serveCmd.Name())
	databaseFs := config.LoadDatabaseFlags(serveCmd.Name())
	smtpFs := config.LoadSMTPFlags(serveCmd.Name())
	captchaFs := config.LoadCaptchaFlags(serveCmd.Name())

	serveCmd.Flags().AddFlagSet(generalFs)
	serveCmd.Flags().AddFlagSet(serverFs)
	serveCmd.Flags().AddFlagSet(databaseFs)
	serveCmd.Flags().AddFlagSet(smtpFs)
	serveCmd.Flags().AddFlagSet(captchaFs)
}
