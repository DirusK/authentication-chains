/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/avast/retry-go"
	"github.com/nutsdb/nutsdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	"authentication-chains/internal/cipher"
	"authentication-chains/internal/types"
)

// addAuthenticationEntry registers a device in authentication table.
func (n *Node) addAuthenticationEntry(block *types.Block, level uint32) error {
	if n.cfg.Level < level {
		return fmt.Errorf("can't add entry from upper blockchain: node level %d < entry level %d", n.cfg.Level, level)
	}

	entry := &types.AuthenticationEntry{
		DeviceId:      block.Dar.DeviceId,
		ClusterHeadId: block.Dar.ClusterHeadId,
		BlockHash:     block.Hash,
		BlockIndex:    block.Index,
	}

	data, err := proto.Marshal(entry)
	if err != nil {
		return err
	}

	if err = n.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put(bucketAuthTableLevel(level), entry.DeviceId, data, types.InfinityTTL)
	}); err != nil {
		return err
	}

	return nil
}

// verifyAuthentication verifies the authentication of the device by authentication table.
func (n *Node) verifyAuthentication(deviceID, blockHash []byte) error {
	var entry types.AuthenticationEntry

	if err := n.db.View(func(tx *nutsdb.Tx) error {
		for i := n.cfg.Level; i >= 0; i-- {
			data, err := tx.Get(bucketAuthTableLevel(i), deviceID)
			if data != nil {
				if err = proto.Unmarshal(data.Value, &entry); err != nil {
					return err
				}

				break
			}
		}

		return nil
	}); err != nil {
		return err
	}

	if !bytes.Equal(entry.BlockHash, blockHash) {
		return fmt.Errorf("%w: block hash mismatch", ErrVerification)
	}

	return nil
}

// initClient initializes a new client.
func initClient(ctx context.Context, address string) (types.NodeClient, error) {
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()

	return types.NewNodeClient(conn), nil
}

func (n *Node) createDAR() (*types.DeviceAuthenticationRequest, error) {
	dar := &types.DeviceAuthenticationRequest{
		DeviceId:      n.cipher.SerializePublicKey(),
		ClusterHeadId: n.clusterHead.DeviceID,
	}

	if err := n.cipher.SignDAR(dar); err != nil {
		return nil, err
	}

	return dar, nil
}

func (n *Node) initMaster(ctx context.Context) error {
	ctx, logger := n.logger.StartTrace(ctx, "init master")
	defer logger.FinishTrace()

	dar, err := n.createDAR()
	if err != nil {
		logger.Errorf("create dar: %s", err)
		return err
	}

	n.chain.SetGenesisHash([]byte(n.cfg.GenesisHash))
	n.chain.AddToMemPool(dar)

	return nil
}

// initPeers - initializes cluster head and peers.
func (n *Node) initPeers(ctx context.Context) error {
	ctx, logger := n.logger.StartTrace(ctx, "init peers")
	defer logger.FinishTrace()

	client, err := initClient(ctx, n.cfg.ClusterHeadGRPCAddress)
	if err != nil {
		logger.Errorf("init cluster head client: %s", err)
		return err
	}

	status, err := client.GetStatus(ctx, &types.StatusRequest{})
	if err != nil {
		logger.Errorf("get cluster head status: %s", err)
		return err
	}

	n.clusterHead = NewPeer(
		status.Peer.Name,
		status.Peer.DeviceId,
		status.Peer.ClusterHeadId,
		status.Peer.GrpcAddress,
		status.Peer.Level,
		client,
	)

	dar, err := n.createDAR()
	if err != nil {
		logger.Errorf("create dar: %s", err)
		return err
	}

	registerResponse, err := n.clusterHead.Client.RegisterNode(ctx, &types.NodeRegistrationRequest{
		Name:    n.cfg.Name,
		Address: n.cfg.GRPC.Address,
		Dar:     dar,
	})
	if err != nil {
		logger.Errorf("register node: %s", err)
		return ErrInvalidDAR
	}

	if registerResponse.IsRegistered {
		n.chain.SetGenesisHash(registerResponse.GenesisHash)
		n.chain.AddToMemPool(dar)
	} else {
		logger.Errorf("node is not registered")
		return ErrInvalidDAR
	}

	for _, peer := range registerResponse.GetPeers() {
		client, err = initClient(ctx, peer.GrpcAddress)
		if err != nil {
			logger.Errorf("init peer client: %s", err)
			return err
		}

		n.clusterNodes.Add(NewPeer(
			status.Peer.Name,
			status.Peer.DeviceId,
			status.Peer.ClusterHeadId,
			status.Peer.GrpcAddress,
			status.Peer.Level,
			client))
	}

	return nil
}

// syncBlocks syncs blocks from the node.
func (n *Node) syncBlocks(ctx context.Context, peer *Peer, from, to uint64) error {
	ctx, logger := n.logger.StartTrace(ctx, "sync blocks from node "+peer.Name)
	defer logger.FinishTrace()

	response, err := peer.Client.GetBlocks(ctx, &types.BlocksRequest{
		From: from,
		To:   to,
	})
	if err != nil {
		logger.Errorf("get blocks from node %s: %s", peer.Name, err)
		return err
	}

	for _, block := range response.Blocks {
		if err = n.addBlock(ctx, block); err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) addBlock(ctx context.Context, block *types.Block) error {
	ctx, logger := n.logger.StartTrace(ctx, "add block")
	defer logger.FinishTrace()

	if err := n.validateBlock(ctx, block); err != nil {
		logger.Errorf("validate block %x: %s", block.Hash, err)
		return err
	}

	if err := n.chain.AddBlock(block); err != nil {
		logger.Errorf("add block %x: %s", block.Hash, err)
		return err
	}

	if err := n.addAuthenticationEntry(block, n.cfg.Level); err != nil {
		logger.Errorf("add authentication entry: %s", err)
		return err
	}

	return nil
}

func (n *Node) verifyDAR(ctx context.Context, dar *types.DeviceAuthenticationRequest) (bool, error) {
	ctx, logger := n.logger.StartTrace(ctx, "add dar")
	defer logger.FinishTrace()

	if err := cipher.VerifyDAR(dar); err != nil {
		if errors.Is(err, cipher.ErrDARVerification) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// validateBlock validates the block.
func (n *Node) validateBlock(ctx context.Context, block *types.Block) error {
	ctx, logger := n.logger.StartTrace(ctx, "validate block")
	defer logger.FinishTrace()

	var clusterHeadID []byte
	if n.clusterHead != nil {
		clusterHeadID = n.clusterHead.DeviceID
	}

	if !bytes.Equal(block.Dar.ClusterHeadId, n.cipher.SerializePublicKey()) ||
		bytes.Equal(block.Dar.ClusterHeadId, clusterHeadID) {
		return fmt.Errorf("%w: invalid cluster head", ErrBlockValidation)
	}

	hash, err := n.cipher.HashBlock(block)
	if err != nil {
		return err
	}

	if !bytes.Equal(hash, block.Hash) {
		return fmt.Errorf("%w: hash mismatch", ErrBlockValidation)
	}

	if err = cipher.VerifyDAR(block.Dar); err != nil {
		return fmt.Errorf("%w: invalid dar", ErrBlockValidation)
	}

	return nil
}

func (n *Node) sendDAR(ctx context.Context, request *types.DeviceAuthenticationRequest) (*types.DeviceAuthenticationResponse, error) {
	ctx, logger := n.logger.StartTrace(ctx, "send dar")
	defer logger.FinishTrace()

	var (
		response = &types.DeviceAuthenticationResponse{}
		err      error
	)

	response.IsVerified, err = n.verifyDAR(ctx, request)
	if err != nil {
		return nil, err
	}

	switch {
	// if dar from children node
	case bytes.Equal(request.ClusterHeadId, n.cipher.SerializePublicKey()):
		return response, nil

	default: // if dar from cluster node
		n.chain.AddToMemPool(request)
	}

	if err = retry.Do(func() error {
		lastBlock := n.chain.GetLastBlock()

		if bytes.Equal(lastBlock.Dar.DeviceId, request.DeviceId) {
			response.BlockHash = lastBlock.Hash
		} else {
			return ErrBlockHasNotMined
		}

		return nil
	},
		retry.Attempts(n.cfg.WaitBlock.Attempts),
		retry.Delay(n.cfg.WaitBlock.Interval),
		retry.Context(ctx),
		retry.LastErrorOnly(true),
	); err != nil {
		return nil, ErrWaitBlockTimeout
	}

	return response, nil
}

// bucketAuthTableLevel returns the name of the bucket that will store authentication table by level.
func bucketAuthTableLevel(level uint32) string {
	return fmt.Sprintf("%s %d", types.BucketAuthenticationTable, level)
}
