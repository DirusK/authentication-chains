/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"context"

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

		deviceID []byte

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
		deviceID:   cipher.SerializePublicKey(),
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

// Sync syncs blocks from the cluster.
func (n *Node) Sync(ctx context.Context) {
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
}
