/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package client

import (
	"time"

	"github.com/spf13/cobra"
)

var cfgName string

const (
	defaultConfigPath  = "configs/clients/%s.yaml"
	defaultConfigName  = "default"
	defaultGRPCAddress = "localhost:50051"
	defaultGRPCTimeout = 15 * time.Second
)

// ClientCmd represents the client command
var ClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Interact with the blockchain as a client",
}

func init() {
	ClientCmd.PersistentFlags().StringVarP(&cfgName, "name", "n", defaultConfigName, "name for the config file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
