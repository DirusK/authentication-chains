/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package client

import (
	"context"
	"fmt"
	"os"

	"github.com/DirusK/utils/config"
	"github.com/DirusK/utils/printer"
	"gopkg.in/yaml.v3"

	"authentication-chains/internal/cipher"
	cfg "authentication-chains/internal/config"
	"authentication-chains/internal/types"
)

const tag = "client"

type Client struct {
	config cfg.Client
	ctx    context.Context
	cipher cipher.Cipher
	client types.NodeClient
}

func New(ctx context.Context, configPath string) (*Client, error) {
	var cfg cfg.Client
	if err := config.LoadFromFile(configPath, &cfg); err != nil {
		printer.Errort(tag, err, "Failed to load config")
		return nil, err
	}

	c, err := cipher.FromStringPrivateKey(cfg.Keys.PrivateKey)
	if err != nil {
		printer.Errort(tag, err, "Failed to load private key")
		return nil, err
	}

	client, err := initClient(ctx, cfg.GRPC.Address)
	if err != nil {
		printer.Errort(tag, err, "Failed to grpc init client")
		return nil, err
	}

	return &Client{
		ctx:    ctx,
		config: cfg,
		cipher: c,
		client: client,
	}, nil
}

func (c *Client) SendDAR() (string, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.config.GRPC.Timeout)
	defer cancel()

	printer.Infot(tag, "Getting status", "address", c.config.GRPC.Address)

	status, err := c.client.GetStatus(ctx, &types.StatusRequest{})
	if err != nil {
		printer.Errort(tag, err, "Failed to get status from node", "address", c.config.GRPC.Address)
		return "", err
	}

	dar := &types.DeviceAuthenticationRequest{
		DeviceId:      c.cipher.SerializePublicKey(),
		ClusterHeadId: status.Peer.ClusterHeadId,
		Signature:     nil,
	}

	if err = c.cipher.SignDAR(dar); err != nil {
		printer.Errort(tag, err, "Failed to sign DAR")
		return "", err
	}

	printer.Infot(tag, "Creating DAR",
		"device_id", fmt.Sprintf("\n%s\n", dar.DeviceId),
		"cluster_head_id", fmt.Sprintf("\n%s\n", dar.ClusterHeadId),
		"signature", fmt.Sprintf("%x", dar.Signature),
	)

	printer.Infot(tag, "Sending DAR", "address", c.config.GRPC.Address, "name", status.Peer.Name)

	response, err := c.client.SendDAR(ctx, dar)
	if err != nil {
		printer.Errort(tag, err, "DAR is not verified")
		return "", err
	}

	printer.Infot(tag, "DAR is verified", "block_hash", fmt.Sprintf("%x", response.BlockHash))

	return fmt.Sprintf("%x", response.BlockHash), nil
}

func (c *Client) GetBlocks(ctx context.Context, from, to uint64) ([]*types.Block, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.config.GRPC.Timeout)
	defer cancel()

	response, err := c.client.GetBlocks(ctx, &types.BlocksRequest{
		From: from,
		To:   to,
	})
	if err != nil {
		printer.Errort(tag, err, "Failed to get blocks")
		return nil, err
	}

	return response.Blocks, nil
}

func (c *Client) SaveBlockHash(configPath, hash string) error {
	c.config.BlockHash = hash

	data, err := yaml.Marshal(c.config)
	if err != nil {
		printer.Errort(tag, err, "Failed to marshal config")
		return err
	}

	if err = os.WriteFile(configPath, data, 0644); err != nil {
		printer.Errort(tag, err, "Failed to write config")
		return err
	}

	printer.Infot(tag, "Block hash is saved and will be used in future requests")

	return nil
}
