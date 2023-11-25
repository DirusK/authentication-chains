/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package config

type (
	Client struct {
		Name      string `yaml:"name" validate:"required"`
		BlockHash string `yaml:"block-hash" validate:"required"`
		GRPC      GRPC   `yaml:"grpc" validate:"required"`
		Keys      Keys   `yaml:"keys" validate:"required"`
	}

	Keys struct {
		PublicKey  string `yaml:"public-key" validate:"required"`
		PrivateKey string `yaml:"private-key" validate:"required"`
	}
)
