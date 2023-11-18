/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

type (
	// Config is a node configuration.
	Config struct {
		Name    string  `yaml:"name" validate:"required"`
		Cluster Cluster `yaml:"cluster"`
		Storage Storage `yaml:"storage"`
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

	// Server is a node server configuration.
	Server struct {
		GRPCAddress string `yaml:"grpc-address" validate:"required"`
	}
)
