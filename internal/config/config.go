/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package config

import (
	"time"

	"github.com/DirusK/utils/log"
)

type (
	// Config is a node configuration.
	Config struct {
		Node       Node       `yaml:"node" validate:"required"`
		Storage    Storage    `yaml:"storage" validate:"required"`
		GRPC       GRPC       `yaml:"grpc" validate:"required"`
		Logger     log.Config `yaml:"logger" validate:"required"`
		WorkerPool WorkerPool `yaml:"worker-pool" validate:"required"`
	}

	// Node is a node cluster configuration.
	Node struct {
		Name                   string   `yaml:"name" validate:"required"`
		Level                  uint32   `yaml:"level" validate:"required"`
		GenesisHash            string   `yaml:"genesis-hash"`
		IsClusterHead          bool     `yaml:"is-cluster-head"`
		ClusterHeadGRPCAddress string   `yaml:"cluster-head-grpc-address"`
		NodesGRPCAddresses     []string `yaml:"nodes-grpc-addresses"`
		WaitBlock              Retry    `yaml:"wait-block"`
		GRPC                   GRPC     `yaml:"grpc" validate:"required"`
	}

	// Storage is a node database configuration.
	Storage struct {
		Directory string `yaml:"directory" validate:"required,dirpath"`
	}

	// GRPC is a node server configuration.
	GRPC struct {
		Address string        `yaml:"address" validate:"required"`
		Timeout time.Duration `yaml:"timeout" validate:"required"`
	}

	// Retry is a retry configuration.
	Retry struct {
		Attempts uint          `yaml:"attempts" validate:"required"`
		Interval time.Duration `yaml:"interval" validate:"required"`
	}

	// Scheduler is a scheduler configuration.
	Scheduler struct {
		Enabled   bool          `yaml:"enabled" validate:"required"`
		Interval  time.Duration `yaml:"interval" validate:"required"`
		Immediate bool          `yaml:"immediate" validate:"required"`
	}

	// WorkerPool defines configuration for workers.
	WorkerPool struct {
		MaxWorkers  int `yaml:"max-workers" valid:"required"`
		MaxCapacity int `yaml:"max-capacity" valid:"required"`
	}
)
