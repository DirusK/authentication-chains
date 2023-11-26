/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"github.com/DirusK/utils/config"
	"github.com/DirusK/utils/printer"
	"gopkg.in/yaml.v3"

	"authentication-chains/internal/cipher"
	cfg "authentication-chains/internal/config"
	"authentication-chains/internal/node"
	"authentication-chains/internal/types"
)

const tag = "client"

type Client struct {
	config cfg.Client
	ctx    context.Context
	cipher cipher.Cipher
	client types.NodeClient
	peer   *node.Peer
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

	status, err := client.GetStatus(ctx, &types.StatusRequest{})
	if err != nil {
		printer.Errort(tag, err, "Failed to get status from node", "address", cfg.GRPC.Address)
		return nil, err
	}

	return &Client{
		ctx:    ctx,
		config: cfg,
		cipher: c,
		client: client,
		peer:   node.NewPeer(status.Peer.Name, status.Peer.DeviceId, status.Peer.ClusterHeadId, status.Peer.GrpcAddress, status.Peer.Level, client),
	}, nil
}

func (c *Client) SendDAR() (string, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.config.GRPC.Timeout)
	defer cancel()

	dar := &types.DeviceAuthenticationRequest{
		DeviceId:      c.cipher.SerializePublicKey(),
		ClusterHeadId: c.peer.ClusterHeadID,
		Signature:     nil,
	}

	if err := c.cipher.SignDAR(dar); err != nil {
		printer.Errort(tag, err, "Failed to sign DAR")
		return "", err
	}

	printer.Infot(tag, "Creating DAR",
		"device_id", fmt.Sprintf("\n%s\n", dar.DeviceId),
		"cluster_head_id", fmt.Sprintf("\n%s\n", dar.ClusterHeadId),
		"signature", fmt.Sprintf("%x", dar.Signature),
	)

	printer.Infot(tag, "Sending DAR", "node", c.peer.Name, "address", c.peer.GRPCAddress, "level", c.peer.Level)

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

	printer.Infot(tag, "Getting blocks",
		"node", c.peer.Name,
		"address", c.peer.GRPCAddress, ""+
			"level", c.peer.Level,
		"from", from,
		"to", to,
	)

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

func (c *Client) GetAuthenticationTable() (*types.AuthenticationTableResponse, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.config.GRPC.Timeout)
	defer cancel()

	printer.Infot(tag, "Getting authentication table",
		"node", c.peer.Name,
		"address", c.peer.GRPCAddress,
		"level", c.peer.Level,
	)

	response, err := c.client.GetAuthenticationTable(ctx, &types.AuthenticationTableRequest{})
	if err != nil {
		printer.Errort(tag, err, "Failed to get authentication table")
		return nil, err
	}

	return response, nil
}

func (c *Client) SendMessage(data []byte) (*types.Content, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.config.GRPC.Timeout)
	defer cancel()

	printer.Infot(tag, "Sending message",
		"node", c.peer.Name,
		"address", c.peer.GRPCAddress,
		"level", c.peer.Level,
	)

	hash, err := hex.DecodeString(c.config.BlockHash)
	if err != nil {
		printer.Errort(tag, err, "Failed to decode block hash")
		return nil, err
	}

	content := &types.Content{
		Data:      data,
		BlockHash: hash,
	}

	pubKey, err := cipher.DeserializePublicKey(c.peer.DeviceID)
	if err != nil {
		printer.Errort(tag, err, "Failed to deserialize peer public key")
		return nil, err
	}

	encryptedMessage, err := cipher.EncryptContent(pubKey, content)
	if err != nil {
		printer.Errort(tag, err, "Failed to encrypt message")
		return nil, err
	}

	response, err := c.client.SendMessage(ctx, &types.Message{
		SenderId:   c.cipher.SerializePublicKey(),
		ReceiverId: c.peer.DeviceID,
		Data:       encryptedMessage,
	})
	if err != nil {
		printer.Errort(tag, err, "Failed to send message")
		return nil, err
	}

	content, err = c.cipher.DecryptContent(response.Data)
	if err != nil {
		printer.Errort(tag, err, "Failed to decrypt message")
		return nil, err
	}

	if !bytes.Equal(response.ReceiverId, c.cipher.SerializePublicKey()) {
		err = errors.New("response receiver id is not equal to device id")
		printer.Errort(tag, err)
		return nil, err
	}

	if !bytes.Equal(response.SenderId, c.peer.DeviceID) {
		err = errors.New("response sender id is not equal to peer id")
		printer.Errort(tag, err)
		return nil, err
	}

	return content, nil
}
