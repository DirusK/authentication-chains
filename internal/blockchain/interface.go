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
		// CreateBlock creates a new block from provided device authentication request.
		CreateBlock(dar *types.DeviceAuthenticationRequest) (*types.Block, error)
		// AddBlock adds a block to the chain.
		AddBlock(block *types.Block) error
		// GetBlock returns a block by index.
		GetBlock(index uint64) (*types.Block, error)
		// GetAllBlocks returns all blocks from the chain with pagination.
		GetAllBlocks(from, to uint64) ([]*types.Block, error)
		// GetLastBlock returns the last block of the chain.
		GetLastBlock() *types.Block
		// SetGenesisHash sets the genesis block hash.
		SetGenesisHash(hash []byte)
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
