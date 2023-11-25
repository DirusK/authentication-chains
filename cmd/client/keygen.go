/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package client

import (
	"fmt"
	"os"

	"github.com/DirusK/utils/printer"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"authentication-chains/cmd/helpers"
	"authentication-chains/internal/cipher"
	cfg "authentication-chains/internal/config"
)

// keygenCmd represents the keygen command
var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate a new key pair and saves it to the config file",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cipher, err := cipher.New(nil)
		if err != nil {
			printer.Errort(helpers.TagCLI, err, "Failed to generate a new key pair")
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
				PublicKey:  cipher.ToStringPublicKey(),
				PrivateKey: cipher.ToStringPrivateKey(),
			},
		}

		data, err := yaml.Marshal(config)
		if err != nil {
			printer.Errort(helpers.TagCLI, err, "Failed to marshal config")
			return
		}

		if err := os.WriteFile(fmt.Sprintf(defaultConfigPath, cfgName), data, 0644); err != nil {
			printer.Errort(helpers.TagCLI, err, "Failed to write config file")
			return
		}
	},
}

func init() {
	ClientCmd.AddCommand(keygenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// keygenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
