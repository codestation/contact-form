// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ApplicationKeySize must be 32 bytes to be used as a key by AES-256
const ApplicationKeySize = 32

// genkeyCmd represents the genkey command
var genkeyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "Generate a random key",
	Long:  `Generate a random key to be used as a secret for other configurations`,
	PreRun: func(cmd *cobra.Command, _ []string) {
		cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		length := viper.GetInt("length")
		if length < 8 || length > 8192 {
			return errors.New("key length must be between 8 and 8192")
		}

		key := make([]byte, length)
		_, err := rand.Read(key)
		if err != nil {
			return err
		}

		hexKey := hex.EncodeToString(key)

		outputFile := viper.GetString("output")
		if outputFile == "" {
			if viper.GetBool("quiet") {
				fmt.Printf("%s\n", hexKey)
			} else {
				fmt.Printf("Generated key: %s\n", hexKey)
			}
		} else {
			err = os.WriteFile(outputFile, []byte(hexKey+"\n"), 0o600)
			if err != nil {
				return err
			}
			if !viper.GetBool("quiet") {
				fmt.Printf("Saved generated key to file %s\n", outputFile)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(genkeyCmd)

	genkeyCmd.Flags().StringP("output", "o", "", "Insert generated key to file")
	genkeyCmd.Flags().BoolP("quiet", "q", false, "Do not print extra messages")
	genkeyCmd.Flags().IntP("length", "l", ApplicationKeySize, "Use an specific key length")
}
