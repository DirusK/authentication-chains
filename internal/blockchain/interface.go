/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package blockchain

import (
	"authentication-chains/internal/types"
)

// Blockchain - describe an interface for working with blockchain.
type (
	// Blockchain - describe an interface for working with blockchain.
	Blockchain interface {
		// AddToMemPool adds a device authentication request to the mem-pool.
		AddToMemPool(request *types.DeviceAuthenticationRequest)
		// MineBlock creates a new block from the mem-pool.
		MineBlock() (*types.Block, error)
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
		// DeleteBlocks delete blocks from the chain.
		DeleteBlocks(from, to uint64) error
	}

	// MemPool - describe an interface for working with memory pool.
	MemPool interface {
		// GetFirst returns the first device authentication request from the mem-pool.
		GetFirst() *types.DeviceAuthenticationRequest
		// GetAll returns all device authentication requests from the mem-pool.
		GetAll() []*types.DeviceAuthenticationRequest
		// Add adds a device authentication request to the mem-pool.
		Add(request *types.DeviceAuthenticationRequest)
		// Remove removes a device authentication request from the mem-pool.
		Remove(request *types.DeviceAuthenticationRequest)
	}
)
