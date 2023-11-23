/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"context"

	"authentication-chains/internal/types"
)

// server implements types.NodeServer.
type server struct {
	types.UnimplementedNodeServer
	Node *Node
}

// NewServer creates a new server instance.
func NewServer(n *Node) types.NodeServer {
	return server{Node: n}
}

func (s server) GetStatus(ctx context.Context, request *types.StatusRequest) (*types.StatusResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s server) GetBlock(ctx context.Context, request *types.BlockRequest) (*types.BlockResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s server) GetBlocks(ctx context.Context, request *types.BlocksRequest) (*types.BlocksResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s server) GetGenesisHash(ctx context.Context, request *types.GenesisHashRequest) (*types.GenesisHashResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s server) GetPeers(ctx context.Context, request *types.PeersRequest) (*types.PeersResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s server) SendMessage(ctx context.Context, message *types.Message) (*types.Message, error) {
	// TODO implement me
	panic("implement me")
}

func (s server) SendBlock(ctx context.Context, request *types.BlockValidationRequest) (*types.BlockValidationResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s server) SendDAR(ctx context.Context, request *types.DeviceAuthenticationRequest) (*types.DeviceAuthenticationResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s server) BroadcastBlock(ctx context.Context, request *types.BlockValidationRequest) (*types.BlockValidationResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s server) BroadcastDAR(ctx context.Context, request *types.DeviceAuthenticationRequest) (*types.DeviceAuthenticationResponse, error) {
	// TODO implement me
	panic("implement me")
}
