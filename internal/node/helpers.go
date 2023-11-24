/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"bytes"
	"context"
	"fmt"

	"github.com/nutsdb/nutsdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	"authentication-chains/internal/blockchain"
	"authentication-chains/internal/cipher"
	"authentication-chains/internal/types"
)

// mineBlock mines a new block.
func (n *Node) mineBlock(ctx context.Context, dar *types.DeviceAuthenticationRequest) (*types.Block, error) {
	ctx, logger := n.logger.StartTrace(ctx, "mine block")
	defer logger.FinishTrace()

	block, err := n.chain.CreateBlock(dar)
	if err != nil {
		return nil, err
	}

	if n.clusterNodes != nil {
		peers := n.clusterNodes.GetAll()

		for _, node := range peers {
			response, err := node.Client.SendBlock(ctx, &types.BlockValidationRequest{Block: block})
			if err != nil {
				logger.Errorf("send block to node %s: %s", node.Name, err)
				continue
			}

			if !response.IsValid {
				logger.Errorf("validation by node %s: block %x is not valid", node.Name, block.Hash)
				return nil, blockchain.ErrBlockValidation
			}
		}
	}

	if n.clusterHead != nil {
		response, err := n.clusterHead.Client.SendBlock(ctx, &types.BlockValidationRequest{Block: block})
		if err != nil {
			logger.Errorf("send block to cluster head: %s", err)
			return nil, err
		}

		if !response.IsValid {
			logger.Errorf("validation by cluster head %s: block %x is not valid", n.clusterHead.Name, block.Hash)
			return nil, blockchain.ErrBlockValidation
		}
	}

	if err = n.chain.AddBlock(block); err != nil {
		return nil, err
	}

	if err = n.addAuthenticationEntry(block, n.cfg.Level); err != nil {
		return nil, err
	}

	return block, nil
}

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
func (n *Node) verifyAuthentication(ctx context.Context, deviceID, blockHash []byte) error {
	ctx, logger := n.logger.StartTrace(ctx, "verify authentication")
	defer logger.FinishTrace()

	var (
		entry types.AuthenticationEntry
		level uint32
	)

	if err := n.db.View(func(tx *nutsdb.Tx) error {
		for i := int32(n.cfg.Level); i >= 0; i-- {
			level = uint32(i)
			data, err := tx.Get(bucketAuthTableLevel(level), deviceID)
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

	switch {
	case entry.BlockHash != nil && bytes.Equal(entry.BlockHash, blockHash) && level == n.cfg.Level:
		block, err := n.chain.GetBlock(entry.BlockIndex)
		if err != nil {
			return err
		}

		if !bytes.Equal(block.Hash, blockHash) {
			return fmt.Errorf("%w: block hash mismatch", ErrVerification)
		}

	case entry.BlockHash != nil && bytes.Equal(entry.BlockHash, blockHash):
		return nil

		// for _, peer := range n.childrenNodes.GetAll() {
		// 	response, err := peer.Client.FindBlock(ctx, &types.FindBlockRequest{Index: entry.BlockIndex, Hash: blockHash})
		// 	if err != nil {
		// 		logger.Errorf("get block from node %s: %s", peer.Name, err)
		// 		continue
		// 	}
		//
		// 	if !bytes.Equal(response.Block.Hash, blockHash) {
		// 		return fmt.Errorf("%w: block hash mismatch", ErrVerification)
		// 	}
		// }

	case entry.BlockHash != nil && !bytes.Equal(entry.BlockHash, blockHash):
		return fmt.Errorf("%w: block hash mismatch", ErrVerification)

	case entry.BlockHash == nil && n.clusterHead != nil:
		verifyResponse, err := n.clusterHead.Client.VerifyDevice(ctx, &types.VerifyDeviceRequest{
			DeviceId:  deviceID,
			BlockHash: blockHash,
		})
		if err != nil {
			return err
		}

		if !verifyResponse.IsVerified {
			return fmt.Errorf("%w: device is not registered", ErrVerification)
		}

		return nil

	default:
		return fmt.Errorf("%w: device is not registered", ErrVerification)
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

func (n *Node) getClusterHeadDeviceID() []byte {
	var clusterHeadID []byte
	if n.clusterHead != nil {
		clusterHeadID = n.clusterHead.DeviceID
	}

	return clusterHeadID
}

func (n *Node) createDAR() (*types.DeviceAuthenticationRequest, error) {
	dar := &types.DeviceAuthenticationRequest{
		DeviceId:      n.deviceID,
		ClusterHeadId: n.getClusterHeadDeviceID(),
	}

	if err := n.cipher.SignDAR(dar); err != nil {
		return nil, err
	}

	return dar, nil
}

func (n *Node) initMaster(ctx context.Context) error {
	ctx, logger := n.logger.StartTrace(ctx, "init master")
	defer logger.FinishTrace()

	firstBlock := n.chain.GetFirstBlock()
	if firstBlock != nil {
		n.authBlockHash = firstBlock.Hash
		return nil
	}

	dar, err := n.createDAR()
	if err != nil {
		logger.Errorf("create dar: %s", err)
		return err
	}

	n.chain.SetGenesisHash([]byte(n.cfg.GenesisHash))

	block, err := n.mineBlock(ctx, dar)
	if err != nil {
		logger.Errorf("mine block: %s", err)
		return err
	}

	n.authBlockHash = block.Hash

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

	firstBlock := n.chain.GetFirstBlock()

	switch firstBlock != nil {
	case true:
		n.authBlockHash = firstBlock.Hash

	case false:
		dar, err := n.createDAR()
		if err != nil {
			logger.Errorf("create dar: %s", err)
			return err
		}

		registerResponse, err := n.clusterHead.Client.RegisterNode(ctx, &types.NodeRegistrationRequest{
			Node: &types.Peer{
				Name:          n.cfg.Name,
				Level:         n.cfg.Level,
				DeviceId:      n.deviceID,
				ClusterHeadId: n.clusterHead.ClusterHeadID,
				GrpcAddress:   n.cfg.GRPC.Address,
			},
			Dar: dar,
		})
		if err != nil {
			logger.Errorf("register node: %s", err)
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

		n.Sync(ctx)

		n.chain.SetGenesisHash(registerResponse.GenesisHash)

		block, err := n.mineBlock(ctx, dar)
		if err != nil {
			logger.Errorf("mine block: %s", err)
			return err
		}

		n.authBlockHash = block.Hash
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

// func (n *Node) verifyDAR(ctx context.Context, dar *types.DeviceAuthenticationRequest) (bool, error) {
// 	ctx, logger := n.logger.StartTrace(ctx, "verify dar")
// 	defer logger.FinishTrace()
//
// 	if err := cipher.VerifyDAR(dar); err != nil {
// 		if errors.Is(err, cipher.ErrDARVerification) {
// 			return false, nil
// 		}
//
// 		return false, err
// 	}
//
// 	return true, nil
// }

// validateBlock validates the block.
func (n *Node) validateBlock(ctx context.Context, block *types.Block) error {
	ctx, logger := n.logger.StartTrace(ctx, "validate block")
	defer logger.FinishTrace()

	if !bytes.Equal(block.Dar.ClusterHeadId, n.deviceID) ||
		bytes.Equal(block.Dar.ClusterHeadId, n.getClusterHeadDeviceID()) {
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

// bucketAuthTableLevel returns the name of the bucket that will store authentication table by level.
func bucketAuthTableLevel(level uint32) string {
	return fmt.Sprintf("%s %d", types.BucketAuthenticationTable, level)
}
