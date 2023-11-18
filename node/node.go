/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"bytes"
	"crypto/rsa"
	"fmt"

	"github.com/nutsdb/nutsdb"
	"google.golang.org/protobuf/proto"

	"authentication-chains/blockchain"
	"authentication-chains/cipher"
	"authentication-chains/types"
)

type (
	// Node implements node logic.
	Node struct {
		name          string
		cipher        cipher.Cipher
		chain         blockchain.Blockchain
		db            *nutsdb.DB
		clusterHead   *KnownNode
		clusterNodes  []KnownNode
		level         uint
		isClusterHead bool
	}

	// KnownNode is a node that is known to the current node.
	KnownNode struct {
		name          string
		deviceID      rsa.PublicKey
		clusterHeadID rsa.PublicKey
	}
)

// New creates a new node instance.
func New(cfg Config) (Node, error) {
	db, err := nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir(cfg.Storage.Directory),
	)
	if err != nil {
		return Node{}, err
	}

	chain, err := blockchain.New(db)
	if err != nil {
		return Node{}, err
	}

	return Node{
		name:          cfg.Name,
		cipher:        cipher.New(),
		chain:         chain,
		db:            db,
		isClusterHead: cfg.Cluster.IsClusterHead,
		clusterHead:   nil,
		clusterNodes:  nil,
	}, nil
}

// Name returns the name of the node.
func (n *Node) Name() string {
	return n.name
}

// SerializePublicKey serializes the public key.
func (n *Node) SerializePublicKey() []byte {
	return n.cipher.SerializePublicKey()
}

// addAuthenticationEntry registers a device in authentication table.
func (n *Node) addAuthenticationEntry(entry *types.AuthenticationEntry, level uint) error {
	if n.level < level {
		return fmt.Errorf("can't add entry from upper blockchain: node level %d < entry level %d", n.level, level)
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
		for i := n.level; i >= 0; i-- {
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
