/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package blockchain

import (
	"bytes"
	"fmt"

	"github.com/nutsdb/nutsdb"

	"authentication-chains/types"
)

//go:generate ifacemaker -f blockchain.go -s blockchain -p blockchain -i Blockchain -y "Blockchain - describe an interface for working with blockchain."

// blockchain implements chain logic.
type blockchain struct {
	lastBlock *types.Block
	db        *nutsdb.DB
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
		db:        db,
	}, nil
}

// AddBlock adds a block to the chain.
func (b *blockchain) AddBlock(block *types.Block) error {
	if err := b.db.Update(func(tx *nutsdb.Tx) error {
		if b.lastBlock != nil {
			data, _ := tx.Get(types.BucketBlocks, block.Hash)
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

		if err := tx.Put(types.BucketBlocks, block.Hash, block.Serialize(), types.InfinityTTL); err != nil {
			return err
		}

		b.lastBlock = block

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// GetBlock returns a block by hash.
func (b *blockchain) GetBlock(hash []byte) (*types.Block, error) {
	var block *types.Block

	if err := b.db.View(func(tx *nutsdb.Tx) error {
		entry, err := tx.Get(types.BucketBlocks, hash)
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
func (b *blockchain) GetAllBlocks(offset, limit int, reverse bool) ([]*types.Block, error) {
	var (
		blocks []*types.Block
		skip   int
	)

	if err := b.db.View(func(tx *nutsdb.Tx) error {
		iterator := nutsdb.NewIterator(tx, types.BucketBlocks, nutsdb.IteratorOptions{Reverse: reverse})

		for {
			if skip >= offset {
				data, err := iterator.Value()
				if err != nil {
					return err
				}

				blocks = append(blocks, types.DeserializeBlock(data))
			} else {
				skip++
			}

			if len(blocks) >= limit || !iterator.Next() {
				break
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return blocks, nil
}

// GetLastBlock returns the last block of the chain.
func (b *blockchain) GetLastBlock() *types.Block {
	return b.lastBlock
}

// DeleteLastBlock deletes the last block from the chain.
func (b *blockchain) DeleteLastBlock() error {
	if err := b.db.Update(func(tx *nutsdb.Tx) error {
		if b.lastBlock != nil {
			if err := tx.Delete(types.BucketBlocks, b.lastBlock.Hash); err != nil {
				return err
			}

			b.lastBlock = nil

			iterator := nutsdb.NewIterator(tx, types.BucketBlocks, nutsdb.IteratorOptions{Reverse: true})

			data, err := iterator.Value()
			if err != nil {
				return nil // no blocks in the chain
			}

			b.lastBlock = types.DeserializeBlock(data)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
