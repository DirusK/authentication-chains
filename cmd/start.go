/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package cmd

import (
	"github.com/DirusK/utils/printer"
	"github.com/spf13/cobra"

	"authentication-chains/internal/app"
	"authentication-chains/internal/config"
)

// cfgPath is a path to configuration file.
var cfgPath string

const tagNode = "NODE"

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a blockchain node",
	Run: func(cmd *cobra.Command, args []string) {
		app.New(ctx, cfgPath).Run()
		printer.Infot(tagNode, "finished successfully")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", config.DefaultPath, "Path to configuration file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
