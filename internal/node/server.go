/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"authentication-chains/internal/types"
)

var _ types.NodeServer = (*Node)(nil)

func (n *Node) GetStatus(ctx context.Context, _ *types.StatusRequest) (*types.StatusResponse, error) {
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
			DeviceId:      n.deviceID,
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

// func (n *Node) FindBlock(ctx context.Context, request *types.FindBlockRequest) (*types.FindBlockResponse, error) {
// 	ctx, logger := n.logger.StartTrace(ctx, "find block")
// 	defer logger.FinishTrace()
//
// 	block, err := n.chain.GetBlock(request.Index)
// 	if err == nil {
// 		return &types.FindBlockResponse{Block: block}, nil
// 	}
//
// 	for _, peer := range n.childrenNodes.GetAll() {
// 		response, err := peer.Client.FindBlock(ctx, request)
// 		if err == nil {
// 			return &types.FindBlockResponse{Block: response.Block}, nil
// 		}
// 	}
//
// 	return nil, ErrNotFoundBlock
// }

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
	case bytes.Equal(request.Block.Dar.ClusterHeadId, n.deviceID):
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
	ctx, logger := n.logger.StartTrace(ctx, "broadcast dar")
	defer logger.FinishTrace()

	if err := n.cipher.VerifyDAR(request); err != nil {
		return nil, err
	}

	block, err := n.mineBlock(ctx, request)
	if err != nil {
		return nil, err
	}

	return &types.DeviceAuthenticationResponse{BlockHash: block.Hash}, nil
}

func (n *Node) SendMessage(ctx context.Context, message *types.Message) (*types.Message, error) {
	ctx, logger := n.logger.StartTrace(ctx, "send message")
	defer logger.FinishTrace()

	if !bytes.Equal(message.ReceiverId, n.deviceID) {
		return nil, ErrInvalidMessageReceiver
	}

	reqContent, err := n.cipher.DecryptContent(message.Data)
	if err != nil {
		return nil, err
	}

	if err = n.verifyAuthentication(ctx, message.SenderId, reqContent.BlockHash); err != nil {
		return nil, err
	}

	respContent := &types.Content{
		BlockHash: n.authBlockHash,
		Data:      []byte("You are authenticated and message is received: " + string(reqContent.Data)),
	}

	data, err := n.cipher.EncryptContent(respContent)
	if err != nil {
		return nil, err
	}

	return types.NewMessage(n.deviceID, message.SenderId, data), nil
}

func (n *Node) RegisterNode(ctx context.Context, request *types.NodeRegistrationRequest) (*types.NodeRegistrationResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "register node")
	defer logger.FinishTrace()

	if err := n.cipher.VerifyDAR(request.Dar); err != nil {
		logger.Errorf("verify dar: %s", err)
		return nil, err
	}

	client, err := initClient(ctx, request.Node.GrpcAddress)
	if err != nil {
		logger.Errorf("init cluster head client: %s", err)
		return nil, err
	}

	n.childrenNodes.Add(NewPeer(
		request.Node.Name,
		request.Node.DeviceId,
		request.Node.ClusterHeadId,
		request.Node.GrpcAddress,
		request.Node.Level,
		client,
	))

	return &types.NodeRegistrationResponse{
		GenesisHash: n.chain.GetFirstBlock().Hash,
		Peers:       n.childrenNodes.ToProto(),
	}, nil
}

func (n *Node) VerifyDevice(ctx context.Context, request *types.VerifyDeviceRequest) (*types.VerifyDeviceResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "verify device")
	defer logger.FinishTrace()

	if err := n.verifyAuthentication(ctx, request.DeviceId, request.BlockHash); err != nil {
		return &types.VerifyDeviceResponse{IsVerified: false}, err
	}

	return &types.VerifyDeviceResponse{IsVerified: true}, nil
}
