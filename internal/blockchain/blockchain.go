/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package blockchain

import (
	"bytes"
	"fmt"

	"github.com/nutsdb/nutsdb"

	"authentication-chains/internal/types"
)

//go:generate ifacemaker -f blockchain.go -s blockchain -p blockchain -i Blockchain -y "Blockchain - describe an interface for working with blockchain."

// blockchain implements chain logic.
type blockchain struct {
	lastBlock   *types.Block
	mempool     MemPool
	genesisHash []byte
	db          *nutsdb.DB
}

// New creates a new blockchain instance.
func New(db *nutsdb.DB) (Blockchain, error) {
	var lastBlock *types.Block

	if err := db.View(func(tx *nutsdb.Tx) error {
		iterator := nutsdb.NewIterator(tx, types.BucketBlocks, nutsdb.IteratorOptions{Reverse: true})

		data, err := iterator.Value()
		if err != nil {
			return nil // no blocks in the chain
		}

		lastBlock = types.DeserializeBlock(data)

		return nil
	}); err != nil {
		return nil, err
	}

	return &blockchain{
		lastBlock: lastBlock,
		mempool:   newMemPool(),
		db:        db,
	}, nil
}

// AddToMemPool adds a device authentication request to the mem-pool.
func (b *blockchain) AddToMemPool(request *types.DeviceAuthenticationRequest) {
	b.mempool.Add(request)
}

// SetGenesisHash sets the genesis block hash.
func (b *blockchain) SetGenesisHash(hash []byte) {
	b.genesisHash = hash
}

// CreateBlock creates a new block from the mem-pool.
func (b *blockchain) CreateBlock() (*types.Block, error) {
	dar := b.mempool.GetFirst()
	if dar == nil {
		return nil, ErrEmptyMemPool
	}

	if b.lastBlock != nil {
		return types.NewBlock(b.lastBlock.Hash, b.lastBlock.Index, dar), nil
	}

	return types.NewBlock(b.genesisHash, 0, dar), nil
}

// AddBlock adds a block to the chain.
func (b *blockchain) AddBlock(block *types.Block) error {
	return b.db.Update(func(tx *nutsdb.Tx) error {
		if b.lastBlock != nil {
			if block.Index != b.lastBlock.Index+1 {
				return fmt.Errorf("%w: block %x: index is not valid", ErrBlockValidation, block.Hash)
			}

			data, _ := tx.Get(types.BucketBlocks, uint64ToBytes(block.Index))
			if data != nil {
				return fmt.Errorf("%w: block %x: already exists", ErrBlockValidation, block.Hash)
			}

			if !bytes.Equal(block.PrevHash, b.lastBlock.Hash) {
				return fmt.Errorf("%w: block %x: prev hash is not valid", ErrBlockValidation, block.Hash)
			}
		}

		if err := tx.Put(types.BucketBlocks, block.Hash, block.Serialize(), types.InfinityTTL); err != nil {
			return err
		}

		b.lastBlock = block

		b.mempool.Remove(block.Dar)

		return nil
	})
}

// GetBlock returns a block by index.
func (b *blockchain) GetBlock(index uint64) (*types.Block, error) {
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
	if b.lastBlock == nil {
		return nil
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

// DeleteLastBlock deletes the last block from the chain.
func (b *blockchain) DeleteLastBlock() error {
	return b.db.Update(func(tx *nutsdb.Tx) error {
		lastIndex := b.lastBlock.Index

		if b.lastBlock != nil {
			if err := tx.Delete(types.BucketBlocks, uint64ToBytes(lastIndex)); err != nil {
				return err
			}

			b.lastBlock = nil

			entry, err := tx.Get(types.BucketBlocks, uint64ToBytes(lastIndex-1))
			if err != nil {
				return err
			}

			b.lastBlock = types.DeserializeBlock(entry.Value)

			b.mempool.Add(b.lastBlock.Dar)
		}

		return nil
	})
}

// DeleteBlocks delete blocks from the chain.
func (b *blockchain) DeleteBlocks(from, to uint64) error {
	return b.db.Update(func(tx *nutsdb.Tx) error {
		entries, err := tx.RangeScan(types.BucketBlocks, uint64ToBytes(from), uint64ToBytes(to))
		if err != nil {
			return err
		}

		for _, entry := range entries {
			block := types.DeserializeBlock(entry.Value)

			if err = tx.Delete(types.BucketBlocks, entry.Key); err != nil {
				return err
			}

			b.mempool.Add(block.Dar)
		}

		return nil
	})
}
