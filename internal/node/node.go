/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"context"
	"errors"

	"github.com/DirusK/utils/log"
	"github.com/alitto/pond"
	"github.com/nutsdb/nutsdb"

	"authentication-chains/internal/blockchain"
	"authentication-chains/internal/cipher"
	"authentication-chains/internal/config"
	"authentication-chains/internal/types"
)

type (
	// Node implements node logic.
	Node struct {
		types.UnimplementedNodeServer

		cfg        config.Node
		cipher     cipher.Cipher
		chain      blockchain.Blockchain
		db         *nutsdb.DB
		logger     log.Logger
		workerPool *pond.WorkerPool

		genesisBlockHash []byte
		authBlockHash    []byte

		clusterHead   *Peer
		clusterNodes  *Peers
		childrenNodes *Peers
	}
)

// New creates a new node instance.
func New(cfg config.Node, db *nutsdb.DB, workerPool *pond.WorkerPool, logger log.Logger) (*Node, error) {
	chain, err := blockchain.New(db)
	if err != nil {
		return nil, err
	}

	cipher, err := cipher.New(db)
	if err != nil {
		return nil, err
	}

	return &Node{
		cfg:        cfg,
		cipher:     cipher,
		chain:      chain,
		db:         db,
		logger:     logger.With("node"),
		workerPool: workerPool,
	}, nil
}

// Init initializes the node.
func (n *Node) Init(ctx context.Context) error {
	ctx, logger := n.logger.StartTrace(ctx, "init")
	defer logger.FinishTrace()

	switch n.cfg.ClusterHeadGRPCAddress == "" {
	case true:
		return n.initMaster(ctx)
	default:
		return n.initPeers(ctx)
	}
}

// MineBlock mines a new block.
func (n *Node) MineBlock(ctx context.Context) error {
	ctx, logger := n.logger.StartTrace(ctx, "mine block")
	defer logger.FinishTrace()

	block, err := n.chain.CreateBlock()
	if err != nil {
		if errors.Is(err, blockchain.ErrEmptyMemPool) {
			logger.Debug("empty mem-pool")
			return nil
		}

		return err
	}

	if n.clusterNodes != nil {
		peers := n.clusterNodes.GetAll()

		for _, node := range peers {
			response, err := node.Client.SendBlock(ctx, &types.BlockValidationRequest{Block: block})
			if err != nil {
				logger.Errorf("send block to node %s: %s", node.Name, err)
				return err
			}

			if !response.IsValid {
				logger.Errorf("validation by node %s: block %x is not valid", node.Name, block.Hash)
				return blockchain.ErrBlockValidation
			}
		}
	}

	if n.clusterHead != nil {
		response, err := n.clusterHead.Client.SendBlock(ctx, &types.BlockValidationRequest{Block: block})
		if err != nil {
			logger.Errorf("send block to cluster head: %s", err)
			return err
		}

		if !response.IsValid {
			logger.Errorf("validation by cluster head %s: block %x is not valid", n.clusterHead.Name, block.Hash)
			return blockchain.ErrBlockValidation
		}
	}

	if err = n.chain.AddBlock(block); err != nil {
		return err
	}

	if err = n.addAuthenticationEntry(block, n.cfg.Level); err != nil {
		return err
	}

	return nil
}

// Sync syncs blocks from the cluster.
func (n *Node) Sync(ctx context.Context) error {
	ctx, logger := n.logger.StartTrace(ctx, "sync")
	defer logger.FinishTrace()

	lastBlock := n.chain.GetLastBlock()

	for _, peer := range n.clusterNodes.GetAll() {
		status, err := peer.Client.GetStatus(ctx, &types.StatusRequest{})
		if err != nil {
			logger.Errorf("get status from peer %s: %s", peer.Name, err)
			continue
		}

		if status.LastBlockIndex > lastBlock.Index {
			if err = n.syncBlocks(ctx, peer, lastBlock.Index+1, status.LastBlockIndex); err != nil {
				continue
			}
		}
	}

	return nil
}
