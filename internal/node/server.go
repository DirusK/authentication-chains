/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"authentication-chains/internal/cipher"
	"authentication-chains/internal/types"
)

func (n *Node) GetStatus(ctx context.Context, request *types.StatusRequest) (*types.StatusResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "get status")
	defer logger.FinishTrace()

	lastBlock := n.chain.GetLastBlock()

	var clusterHeadID []byte
	if n.clusterHead != nil {
		clusterHeadID = n.clusterHead.DeviceID
	}

	return &types.StatusResponse{
		Peer: &types.Peer{
			Name:          n.cfg.Name,
			Level:         n.cfg.Level,
			DeviceId:      n.cipher.SerializePublicKey(),
			ClusterHeadId: clusterHeadID,
			GrpcAddress:   n.cfg.GRPC.Address,
		},
		LastBlockIndex: lastBlock.Index,
		LastBlockHash:  lastBlock.Hash,
	}, nil
}

func (n *Node) GetBlock(ctx context.Context, request *types.BlockRequest) (*types.BlockResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "get block")
	defer logger.FinishTrace()

	block, err := n.chain.GetBlock(request.Index)
	if err != nil {
		return nil, err
	}

	return &types.BlockResponse{
		Block: block,
	}, nil
}

func (n *Node) GetBlocks(ctx context.Context, request *types.BlocksRequest) (*types.BlocksResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "get blocks")
	defer logger.FinishTrace()

	blocks, err := n.chain.GetAllBlocks(request.From, request.To)
	if err != nil {
		return nil, err
	}

	return &types.BlocksResponse{
		Blocks: blocks,
	}, nil
}

func (n *Node) GetPeers(ctx context.Context, request *types.PeersRequest) (*types.PeersResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "get peers")
	defer logger.FinishTrace()

	if request.Level != n.cfg.Level {
		return nil, fmt.Errorf("level %d is not supported", request.Level)
	}

	peers := make([]*types.Peer, 0, len(n.clusterNodes.GetAll()))

	for _, node := range n.clusterNodes.GetAll() {
		peers = append(peers, &types.Peer{
			Name:          node.Name,
			Level:         node.Level,
			DeviceId:      node.DeviceID,
			ClusterHeadId: node.ClusterHeadID,
			GrpcAddress:   node.GRPCAddress,
		})
	}

	return &types.PeersResponse{
		Peers: peers,
	}, nil

}

func (n *Node) SendBlock(ctx context.Context, request *types.BlockValidationRequest) (*types.BlockValidationResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "send block")
	defer logger.FinishTrace()

	response := &types.BlockValidationResponse{}

	switch {
	case request.Block.Dar == nil:
		return response, ErrBlockHasNoDAR

	// if block from children node -> validate and add auth entry
	case bytes.Equal(request.Block.Dar.ClusterHeadId, n.cipher.SerializePublicKey()):
		if err := n.validateBlock(ctx, request.Block); err != nil {
			if errors.Is(err, ErrBlockValidation) {
				response.IsValid = false
				return response, nil
			}

			logger.Errorf("validate block %x: %s", request.Block.Hash, err)
			return response, err
		}

		if err := n.addAuthenticationEntry(request.Block, n.cfg.Level); err != nil {
			logger.Errorf("add authentication entry: %s", err)
			return response, err
		}

		// if block from cluster node -> validate and add to auth table and chain
	default:
		if err := n.addBlock(ctx, request.Block); err != nil {
			if errors.Is(err, ErrBlockValidation) {
				response.IsValid = false
				return response, nil
			}

			return response, err
		}
	}

	response.IsValid = true

	return response, nil
}

func (n *Node) SendDAR(ctx context.Context, request *types.DeviceAuthenticationRequest) (*types.DeviceAuthenticationResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "send dar")
	defer logger.FinishTrace()

	return n.sendDAR(ctx, request)
}

func (n *Node) BroadcastDAR(ctx context.Context, request *types.DeviceAuthenticationRequest) (*types.DeviceAuthenticationResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "broadcast dar")
	defer logger.FinishTrace()

	n.workerPool.

	response := &types.DeviceAuthenticationResponse{}

	if err := cipher.VerifyDAR(request); err != nil {
		if errors.Is(err, cipher.ErrDARVerification) {
			response.IsValid = false
			return response, nil
		}

		return response, err
	}

	n.chain.AddToMemPool(request)

	response.IsValid = true

	return response, nil
}

func (n *Node) RegisterNode(ctx context.Context, request *types.NodeRegistrationRequest) (*types.NodeRegistrationResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (n *Node) SendMessage(ctx context.Context, message *types.Message) (*types.Message, error) {
	// TODO implement me
	panic("implement me")
}
