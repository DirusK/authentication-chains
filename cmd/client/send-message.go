/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package client

import (
	"fmt"

	"github.com/DirusK/utils/printer"
	"github.com/spf13/cobra"

	"authentication-chains/cmd/helpers"
	"authentication-chains/internal/client"
)

// sendMessageCmd represents the sendMessage command
var sendMessageCmd = &cobra.Command{
	Use:   "send-message [message]",
	Short: "Send message to the node",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configPath := fmt.Sprintf(defaultConfigPath, cfgName)

		nodeClient, err := client.New(helpers.Ctx, configPath)
		if err != nil {
			return
		}

		response, err := nodeClient.SendMessage([]byte(args[0]))
		if err != nil {
			return
		}

		printer.Infot(helpers.TagCLI, "Received response", "message", string(response.Data), "block_hash", fmt.Sprintf("%x", response.BlockHash))
	},
}

func init() {
	ClientCmd.AddCommand(sendMessageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sendMessageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendMessageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
