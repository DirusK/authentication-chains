/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"github.com/spf13/cobra"

	"authentication-chains/internal/config"
)

// cfgPath is a path to configuration file.
var cfgPath string

// NodeCmd represents the node command
var NodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Manage a blockchain node",
}

func init() {
	NodeCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", config.DefaultPath, "Path to configuration file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
