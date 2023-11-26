/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"github.com/DirusK/utils/printer"
	"github.com/spf13/cobra"

	"authentication-chains/cmd/helpers"
	"authentication-chains/internal/app"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a blockchain node",
	Run: func(cmd *cobra.Command, args []string) {
		printer.Infot(helpers.TagCLI, "starting node")
		app.New(helpers.Ctx, cfgPath).Run()
		printer.Infot(helpers.TagCLI, "node stopped")
	},
}

func init() {
	NodeCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
