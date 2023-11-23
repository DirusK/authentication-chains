/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"github.com/DirusK/utils/log"
	"github.com/nutsdb/nutsdb"

	"authentication-chains/internal/blockchain"
	"authentication-chains/internal/cipher"
	"authentication-chains/internal/config"
)

type (
	// Node implements node logic.
	Node struct {
		cfg           config.Node
		cipher        cipher.Cipher
		chain         blockchain.Blockchain
		db            *nutsdb.DB
		logger        log.Logger
		clusterHead   *Peer
		clusterNodes  Peers
		childrenNodes Peers
		level         uint
	}
)

// New creates a new node instance.
func New(cfg config.Node, db *nutsdb.DB, logger log.Logger) (*Node, error) {
	chain, err := blockchain.New(db)
	if err != nil {
		return nil, err
	}

	cipher, err := cipher.New(db)
	if err != nil {
		return nil, err
	}

	return &Node{
		cfg:    cfg,
		cipher: cipher,
		chain:  chain,
		db:     db,
		logger: logger,
		level:  cfg.Level,
	}, nil
}
