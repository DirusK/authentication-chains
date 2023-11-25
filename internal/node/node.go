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
func New(ctx context.Context, cfg config.Node, db *nutsdb.DB, workerPool *pond.WorkerPool, logger log.Logger) (*Node, error) {
	chain, err := blockchain.New(db)
	if err != nil {
		return nil, err
	}

	cipher, err := cipher.New(db)
	if err != nil {
		return nil, err
	}

	clusterHead, _ := initPeer(ctx, db, types.BucketClusterHead, types.KeyClusterHead)
	clusterNodes, _ := initPeers(ctx, db, types.BucketClusterNodes)
	childrenNodes, _ := initPeers(ctx, db, types.BucketChildrenNodes)

	return &Node{
		cfg:           cfg,
		cipher:        cipher,
		chain:         chain,
		db:            db,
		logger:        logger.With("node"),
		workerPool:    workerPool,
		deviceID:      cipher.SerializePublicKey(),
		clusterHead:   clusterHead,
		clusterNodes:  clusterNodes,
		childrenNodes: childrenNodes,
	}, nil
}

// Init initializes the node.
func (n *Node) Init(ctx context.Context) error {
	ctx, logger := n.logger.StartTrace(ctx, "init")
	defer logger.FinishTrace()

	var err error

	switch n.cfg.ClusterHeadGRPCAddress == "" {
	case true:
		err = n.initMaster(ctx)
	default:
		err = n.initPeers(ctx)
	}

	if err != nil {
		return err
	}

	logger.Infof("node initialized")

	return nil
}

// Sync syncs blocks from the cluster.
func (n *Node) Sync(ctx context.Context) {
	ctx, logger := n.logger.StartTrace(ctx, "sync")
	defer logger.FinishTrace()

	if n.clusterNodes == nil {
		return
	}

	lastBlock := n.chain.GetLastBlock()

	for _, peer := range n.clusterNodes.GetAll() {
		status, err := peer.Client.GetStatus(ctx, &types.StatusRequest{})
		if err != nil {
			logger.Errorf("get status from peer %s: %s", peer.Name, err)
			continue
		}

		if status.LastBlockIndex > lastBlock.Index {
			logger.Infof("node sync with peer %s from %d block to %d block", peer.Name, lastBlock.Index+1, status.LastBlockIndex)

			if err = n.syncBlocks(ctx, peer, lastBlock.Index+1, status.LastBlockIndex); err != nil {
				continue
			}
		}
	}
}
