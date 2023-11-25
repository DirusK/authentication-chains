/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/DirusK/utils/printer"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"authentication-chains/internal/cipher"
	cfg "authentication-chains/internal/config"
)

var cfgName string

const tagCLI = "CLI"

const (
	defaultConfigPath  = "configs/clients/%s.yaml"
	defaultConfigName  = "default"
	defaultGRPCAddress = "localhost:50051"
	defaultGRPCTimeout = 15 * time.Second
)

// keygenCmd represents the keygen command
var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate a new key pair and saves it to the config file",
	Run: func(cmd *cobra.Command, args []string) {
		cipher, err := cipher.New(nil)
		if err != nil {
			printer.Errort(tagCLI, err, "Failed to generate a new key pair")
			return
		}

		config := cfg.Client{
			Name:      cfgName,
			BlockHash: "",
			GRPC: cfg.GRPC{
				Address: defaultGRPCAddress,
				Timeout: defaultGRPCTimeout,
			},
			Keys: cfg.Keys{
				PublicKey:  cipher.ToHexPublicKey(),
				PrivateKey: cipher.ToHexPrivateKey(),
			},
		}

		data, err := yaml.Marshal(config)
		if err != nil {
			printer.Errort(tagCLI, err, "Failed to marshal config")
			return
		}

		if err := os.WriteFile(fmt.Sprintf(defaultConfigPath, args[0]), data, 0644); err != nil {
			printer.Errort(tagCLI, err, "Failed to write config file")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(keygenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// keygenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	keygenCmd.Flags().StringVarP(&cfgName, "name", "n", defaultConfigName, "name for the config file")
}
