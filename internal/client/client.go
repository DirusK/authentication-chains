/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package client

import (
	"authentication-chains/internal/cipher"
	cfg "authentication-chains/internal/config"
	"authentication-chains/internal/types"
)

type Client struct {
	config cfg.Client
	cipher cipher.Cipher
	client types.NodeClient
}

func New(configPath string) (*Client, error) {
	confi
	c, err := cipher.FromHexPrivateKey()
	if err != nil {
		return nil, err
	}

	return &Client{
		config: config,
		cipher: c,
	}, nil
}
