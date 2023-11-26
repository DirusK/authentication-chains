/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package client

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"authentication-chains/cmd/helpers"
	"authentication-chains/internal/client"
)

// getAuthTableCmd represents the getAuthTable command
var getAuthTableCmd = &cobra.Command{
	Use:   "get-auth-table",
	Short: "Get authentication table from node",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := fmt.Sprintf(defaultConfigPath, cfgName)

		nodeClient, err := client.New(helpers.Ctx, configPath)
		if err != nil {
			return
		}

		authTable, err := nodeClient.GetAuthenticationTable()
		if err != nil {
			return
		}

		for level, authTable := range authTable.Table {
			t := table.NewWriter()
			t.AppendHeader(table.Row{"Device ID", "Cluster Head ID", "Block Hash", "Block Index"})
			t.SetOutputMirror(cmd.OutOrStdout())
			t.SetStyle(table.StyleColoredBright)
			t.SetTitle("Authentication table level %d", level)
			t.SortBy([]table.SortBy{{Name: "Block Index", Mode: table.Asc}})
			t.Style().Title.Align = text.AlignCenter
			t.SetColumnConfigs([]table.ColumnConfig{
				{
					Name:        "Device ID",
					AlignHeader: text.AlignCenter,
					Align:       text.AlignLeft,
				},
				{
					Name:        "Cluster Head ID",
					AlignHeader: text.AlignCenter,
					Align:       text.AlignLeft,
					WidthMin:    33,
					WidthMax:    33,
				},
				{
					Name:        "Block Hash",
					AlignHeader: text.AlignCenter,
					Align:       text.AlignLeft,
				},
				{
					Name:        "Block Index",
					AlignHeader: text.AlignCenter,
					Align:       text.AlignCenter,
				},
			})
			for _, entry := range authTable.Entries {
				t.AppendRow(table.Row{
					helpers.Truncate(string(entry.DeviceId), 30),
					helpers.Truncate(string(entry.ClusterHeadId), 30),
					helpers.Truncate(fmt.Sprintf("%x", entry.BlockHash), 30),
					helpers.Truncate(fmt.Sprintf("%d", entry.BlockIndex), 20),
				})
			}
			t.Render()
			fmt.Println()
		}
	},
}

func init() {
	ClientCmd.AddCommand(getAuthTableCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getAuthTableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getAuthTableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
