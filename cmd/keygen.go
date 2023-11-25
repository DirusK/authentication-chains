/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"authentication-chains/internal/cipher"
	cfg "authentication-chains/internal/config"
)

const (
	defaultConfigPath  = "configs/clients/%s.yaml"
	defaultGRPCAddress = "localhost:50051"
	defaultGRPCTimeout = 15 * time.Second
)

// keygenCmd represents the keygen command
var keygenCmd = &cobra.Command{
	Use:   "keygen [name]",
	Short: "Generate a new key pair",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cipher, err := cipher.New(nil)
		if err != nil {
			return err
		}

		config := cfg.Client{
			Name: args[0],
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
			return err
		}

		if err := os.WriteFile(fmt.Sprintf(defaultConfigPath, args[0]), data, 0644); err != nil {
			return err
		}

		return nil
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
	// keygenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
