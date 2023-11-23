/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"time"
)

type (
	// Config is a node configuration.
	Config struct {
		Meta    Meta    `yaml:"meta" validate:"required"`
		Cluster Cluster `yaml:"cluster"`
		Storage Storage `yaml:"storage"`
		GRPC    GRPC    `yaml:"grpc"`
	}

	// Meta is a node meta configuration.
	Meta struct {
		Name string `yaml:"name" validate:"required"`
	}

	// Cluster is a node cluster configuration.
	Cluster struct {
		Level                  uint     `yaml:"level" validate:"required"`
		GenesisHash            string   `yaml:"genesis-hash"`
		IsClusterHead          bool     `yaml:"is-cluster-head"`
		ClusterHeadGRPCAddress string   `yaml:"cluster-head-grpc-address"`
		NodesGRPCAddresses     []string `yaml:"nodes-grpc-addresses"`
	}

	// Storage is a node database configuration.
	Storage struct {
		Directory string `yaml:"directory" validate:"required,dirpath"`
	}

	// GRPC is a node server configuration.
	GRPC struct {
		Port    string        `yaml:"port" validate:"required"`
		Timeout time.Duration `yaml:"timeout" validate:"required"`
	}
)
