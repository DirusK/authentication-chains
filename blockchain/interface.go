/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package blockchain

import (
	"authentication-chains/types"
)

// Blockchain - describe an interface for working with blockchain.
type Blockchain interface {
	// AddBlock adds a block to the chain.
	AddBlock(block *types.Block) error
	// GetBlock returns a block by index.
	GetBlock(index uint64) (*types.Block, error)
	// GetAllBlocks returns all blocks from the chain with pagination.
	GetAllBlocks(from, to uint64) ([]*types.Block, error)
	// GetLastBlock returns the last block of the chain.
	GetLastBlock() *types.Block
	// DeleteLastBlock deletes the last block from the chain.
	DeleteLastBlock() error
}
