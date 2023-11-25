/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package client

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DirusK/utils/printer"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"

	"authentication-chains/cmd/helpers"
	"authentication-chains/internal/client"
)

// getBlocksCmd represents the get-blocks command
var getBlocksCmd = &cobra.Command{
	Use:   "get-blocks [from] [to]",
	Short: "Get blocks from the blockchain",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		configPath := fmt.Sprintf(defaultConfigPath, cfgName)

		from, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			printer.Errort(helpers.TagCLI, err, "Failed to parse argument from")
			return
		}

		to, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			printer.Errort(helpers.TagCLI, err, "Failed to parse argument to")
			return
		}

		nodeClient, err := client.New(helpers.Ctx, configPath)
		if err != nil {
			return
		}

		blocks, err := nodeClient.GetBlocks(helpers.Ctx, from, to)
		if err != nil {
			return
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"#", "Hash", "Previous Hash", "Timestamp", "DAR"})
		t.SetOutputMirror(cmd.OutOrStdout())
		t.SetStyle(table.StyleColoredBright)
		for _, block := range blocks {
			dar := client.DeviceAuthenticationRequest{
				DeviceID:      helpers.Truncate(fmt.Sprintf("%s", block.Dar.DeviceId), 30),
				ClusterHeadID: helpers.Truncate(fmt.Sprintf("%s", block.Dar.ClusterHeadId), 30),
				Signature:     helpers.Truncate(fmt.Sprintf("%x", block.Dar.Signature), 30),
			}

			t.AppendRow(table.Row{
				block.Index,
				helpers.Truncate(fmt.Sprintf("%x", block.Hash), 20),
				helpers.Truncate(fmt.Sprintf("%x", block.PrevHash), 20),
				helpers.Truncate(time.Unix(block.Timestamp, 0).Format(time.DateTime), 20),
				helpers.Truncate(litter.Sdump(dar), 200),
			})
		}

		t.Render()
	},
}

func init() {
	ClientCmd.AddCommand(getBlocksCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getBlockCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getBlockCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
