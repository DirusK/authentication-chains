/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"authentication-chains/blockchain"
	"authentication-chains/cipher"
)

// Node implements node logic.
type Node struct {
	cipher cipher.Cipher
	chain  blockchain.Blockchain
}
