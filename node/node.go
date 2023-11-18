/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"crypto/rsa"

	"authentication-chains/blockchain"
	"authentication-chains/cipher"
)

type (
	// Node implements node logic.
	Node struct {
		name         string
		cipher       cipher.Cipher
		chain        blockchain.Blockchain
		clusterHead  *KnownNode
		clusterNodes []KnownNode
	}

	// KnownNode is a node that is known to the current node.
	KnownNode struct {
		deviceID            rsa.PublicKey
		clusterHeadDeviceID rsa.PublicKey
	}
)

func New(cfg Config) (Node, error) {

}
