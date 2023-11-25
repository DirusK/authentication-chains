/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"authentication-chains/internal/client"
)

// getBlocksCmd represents the get-blocks command
var getBlocksCmd = &cobra.Command{
	Use:   "get-blocks [from] [to]",
	Short: "Get blocks from the blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := fmt.Sprintf(defaultConfigPath, args[0])

		client, err := client.New(ctx, configPath)
		if err != nil {
			return
		}

		blockHash, err := client.GetBlocks(ctx)
		if err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(getBlocksCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getBlockCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getBlockCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
