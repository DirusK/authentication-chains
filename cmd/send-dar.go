/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"authentication-chains/internal/client"
)

// sendDar represents the send-dar command
var sendDar = &cobra.Command{
	Use:   "send-dar [name]",
	Short: "Send device authentication request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configPath := fmt.Sprintf(defaultConfigPath, args[0])

		client, err := client.New(ctx, configPath)
		if err != nil {
			return
		}

		blockHash, err := client.SendDAR()
		if err != nil {
			return
		}

		if err := client.SaveBlockHash(configPath, blockHash); err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(sendDar)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sendDarCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendDarCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
