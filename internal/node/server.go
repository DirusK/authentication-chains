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

var _ types.NodeServer = (*Node)(nil)

func (n *Node) GetStatus(ctx context.Context, _ *types.StatusRequest) (*types.StatusResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "get status")
	defer logger.FinishTrace()

	logger.Debugw("received status request")

	lastBlock := n.chain.GetLastBlock()

	return &types.StatusResponse{
		Peer: &types.Peer{
			Name:          n.cfg.Name,
			Level:         n.cfg.Level,
			DeviceId:      n.deviceID,
			ClusterHeadId: n.getClusterHeadDeviceID(),
			GrpcAddress:   n.cfg.GRPC.Address,
		},
		LastBlockIndex: lastBlock.Index,
	}, nil
}

func (n *Node) GetBlock(ctx context.Context, request *types.BlockRequest) (*types.BlockResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "get block")
	defer logger.FinishTrace()

	logger.Debugw("received get block request", "index", request.Index)

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

	logger.Debugw("received get blocks request", "from", request.From, "to", request.To)

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

	logger.Debugw("received get peers request", "level", request.Level)

	if request.Level != n.cfg.Level && request.Level != n.cfg.Level-1 {
		return nil, fmt.Errorf("level %d is not supported", request.Level)
	}

	var peers []*types.Peer

	switch {
	case request.Level == n.cfg.Level:
		if n.clusterNodes != nil {
			peers = n.clusterNodes.ToProto()
		}
	case request.Level == n.cfg.Level-1:
		if n.childrenNodes != nil {
			peers = n.childrenNodes.ToProto()
		}
	}

	return &types.PeersResponse{
		Peers: peers,
	}, nil

}

func (n *Node) SendBlock(ctx context.Context, request *types.BlockValidationRequest) (*types.BlockValidationResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "send block")
	logger = logger.WithFields("block_hash", fmt.Sprintf("%x", request.Block.Hash))
	defer logger.FinishTrace()

	logger.Debugw("received send block request")

	response := &types.BlockValidationResponse{}

	switch {
	case request.Block.Dar == nil:
		return response, ErrBlockHasNoDAR

	// if block from children node -> validate and add auth entry
	case bytes.Equal(request.Block.Dar.ClusterHeadId, n.deviceID):
		if err := n.validateBlock(ctx, request.Block); err != nil {
			if errors.Is(err, ErrBlockValidation) {
				response.IsValid = false
				logger.Debugw("block is invalid")
				return response, nil
			}

			logger.Errorf("validate block %x: %s", request.Block.Hash, err)
			return response, err
		}

		if err := n.addAuthenticationEntry(ctx, request.Block, n.cfg.Level-1); err != nil {
			logger.Errorf("add authentication entry: %s", err)
			return response, err
		}

		// if block from cluster node -> validate and add to auth table and chain
	default:
		if err := n.addBlock(ctx, request.Block); err != nil {
			if errors.Is(err, ErrBlockValidation) {
				response.IsValid = false
				logger.Debugw("block is invalid")
				return response, nil
			}

			return response, err
		}
	}

	response.IsValid = true

	logger.Debug("block is valid")

	return response, nil
}

func (n *Node) SendDAR(ctx context.Context, request *types.DeviceAuthenticationRequest) (*types.DeviceAuthenticationResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "broadcast dar")
	defer logger.FinishTrace()

	logger.Debugw("received send dar request", "device_id", string(request.DeviceId))

	if _, err := n.getAuthenticationEntry(ctx, request.DeviceId); err == nil {
		return nil, errors.New("device is already registered in authentication table")
	}

	if err := cipher.VerifyDAR(request); err != nil {
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
	logger = logger.WithFields("sender_id", string(message.SenderId))
	defer logger.FinishTrace()

	logger.Debug("received send message request")

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

	pubKey, err := cipher.DeserializePublicKey(message.SenderId)
	if err != nil {
		return nil, err
	}

	data, err := cipher.EncryptContent(pubKey, respContent)
	if err != nil {
		return nil, err
	}

	logger.Debugw("response message is sent", "message", string(respContent.Data))

	return types.NewMessage(n.deviceID, message.SenderId, data), nil
}

func (n *Node) RegisterNode(ctx context.Context, request *types.NodeRegistrationRequest) (*types.NodeRegistrationResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "register node")
	defer logger.FinishTrace()

	logger.Debugw("received register node request", "node", request.Node.Name)

	client, err := initClient(ctx, request.Node.GrpcAddress)
	if err != nil {
		logger.Errorf("init cluster head client: %s", err)
		return nil, err
	}

	if err = n.addPeer(ctx, NewPeer(
		request.Node.Name,
		request.Node.DeviceId,
		request.Node.ClusterHeadId,
		request.Node.GrpcAddress,
		request.Node.Level,
		client,
	)); err != nil {
		return nil, err
	}

	peers := n.childrenNodes.ToProto()
	for i, peer := range peers {
		if peer.GrpcAddress == request.Node.GrpcAddress {
			peers = append(peers[:i], peers[i+1:]...)
			break
		}
	}

	return &types.NodeRegistrationResponse{
		GenesisHash: n.authBlockHash,
		Peers:       peers,
	}, nil
}

func (n *Node) VerifyDevice(ctx context.Context, request *types.VerifyDeviceRequest) (*types.VerifyDeviceResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "verify device")
	logger = logger.WithFields("device_id", string(request.DeviceId))
	defer logger.FinishTrace()

	logger.Debugw("received verify device request")

	if err := n.verifyAuthentication(ctx, request.DeviceId, request.BlockHash); err != nil {
		return &types.VerifyDeviceResponse{IsVerified: false}, err
	}

	logger.Debugw("device is verified")

	return &types.VerifyDeviceResponse{IsVerified: true}, nil
}

func (n *Node) GetAuthenticationTable(
	ctx context.Context,
	_ *types.AuthenticationTableRequest,
) (*types.AuthenticationTableResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "get authentication table")
	defer logger.FinishTrace()

	logger.Debugw("received get authentication table request")

	table, err := n.getAuthenticationTable(ctx)
	if err != nil {
		return nil, err
	}

	return &types.AuthenticationTableResponse{Table: table}, nil
}
