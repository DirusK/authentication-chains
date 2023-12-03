/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package blockchain

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/nutsdb/nutsdb"
	"google.golang.org/protobuf/proto"

	"authentication-chains/internal/cipher"
	"authentication-chains/internal/types"
)

//go:generate ifacemaker -f blockchain.go -s blockchain -p blockchain -i Blockchain -y "Blockchain - describe an interface for working with blockchain."

// blockchain implements chain logic.
type blockchain struct {
	db          *nutsdb.DB
	lastBlock   *types.Block
	genesisHash []byte
	mutex       sync.RWMutex
}

// New creates a new blockchain instance.
func New(db *nutsdb.DB) (Blockchain, error) {
	var lastBlock *types.Block

	db.View(func(tx *nutsdb.Tx) error {
		lastBlockIndex, err := tx.Get(types.BucketIndexes, types.KeyLastBlock)
		if err != nil {
			return err
		}

		block, err := tx.Get(types.BucketBlocks, lastBlockIndex.Value)
		if err != nil {
			return err
		}

		lastBlock = types.DeserializeBlock(block.Value)

		return nil
	})

	return &blockchain{
		lastBlock: lastBlock,
		mutex:     sync.RWMutex{},
		db:        db,
	}, nil
}

// // AddToMemPool adds a device authentication request to the mem-pool.
// func (b *blockchain) AddToMemPool(request *types.DeviceAuthenticationRequest) {
// 	b.mempool.Add(request)
// }

// SetGenesisHash sets the genesis block hash.
func (b *blockchain) SetGenesisHash(hash []byte) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.genesisHash = hash
}

// CreateBlock creates a new block from the mem-pool.
func (b *blockchain) CreateBlock(dar *types.DeviceAuthenticationRequest) (*types.Block, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	var block *types.Block
	if b.lastBlock != nil {
		block = types.NewBlock(b.lastBlock.Hash, b.lastBlock.Index+1, dar)
	} else {
		block = types.NewBlock(b.genesisHash, 1, dar)
	}

	data, err := proto.Marshal(block)
	if err != nil {
		return nil, err
	}

	block.Hash = cipher.Hash(data)

	return block, nil
}

// AddBlock adds a block to the chain.
func (b *blockchain) AddBlock(block *types.Block) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.db.Update(func(tx *nutsdb.Tx) error {
		if b.lastBlock != nil {
			data, _ := tx.Get(types.BucketBlocks, uint64ToBytes(block.Index))
			if data != nil {
				return fmt.Errorf("%w: block %x: already exists", ErrBlockValidation, block.Hash)
			}

			if block.Index != b.lastBlock.Index+1 {
				return fmt.Errorf("%w: block %x: index is not valid", ErrBlockValidation, block.Hash)
			}

			if !bytes.Equal(block.PrevHash, b.lastBlock.Hash) {
				return fmt.Errorf("%w: block %x: prev hash is not valid", ErrBlockValidation, block.Hash)
			}
		}

		if err := tx.Put(types.BucketBlocks, uint64ToBytes(block.Index), block.Serialize(), types.InfinityTTL); err != nil {
			return err
		}

		if err := tx.Put(types.BucketIndexes, types.KeyLastBlock, uint64ToBytes(block.Index), types.InfinityTTL); err != nil {
			return err
		}

		b.lastBlock = block

		return nil
	})
}

// GetBlock returns a block by index.
func (b *blockchain) GetBlock(index uint64) (*types.Block, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	var block *types.Block

	if err := b.db.View(func(tx *nutsdb.Tx) error {
		entry, err := tx.Get(types.BucketBlocks, uint64ToBytes(index))
		if err != nil {
			return err
		}

		block = types.DeserializeBlock(entry.Value)

		return nil
	}); err != nil {
		return nil, err
	}

	return block, nil
}

// GetAllBlocks returns all blocks from the chain with pagination.
func (b *blockchain) GetAllBlocks(from, to uint64) ([]*types.Block, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	var blocks []*types.Block

	if err := b.db.View(func(tx *nutsdb.Tx) error {
		entries, err := tx.RangeScan(types.BucketBlocks, uint64ToBytes(from), uint64ToBytes(to))
		if err != nil {
			return err
		}

		for _, entry := range entries {
			blocks = append(blocks, types.DeserializeBlock(entry.Value))
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return blocks, nil
}

// GetLastBlock returns the last block of the chain.
func (b *blockchain) GetLastBlock() *types.Block {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.lastBlock == nil {
		return &types.Block{}
	}

	lastBlock := &types.Block{
		Hash:      b.lastBlock.Hash,
		PrevHash:  b.lastBlock.PrevHash,
		Index:     b.lastBlock.Index,
		Dar:       b.lastBlock.Dar,
		Timestamp: b.lastBlock.Timestamp,
	}

	return lastBlock
}
