/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package blockchain

import (
	"github.com/nutsdb/nutsdb"
)

// Blockchain implements chain logic.
type Blockchain struct {
	lastBlockHash []byte
	db            *nutsdb.DB
}

func New(db *nutsdb.DB) *Blockchain {
	err = db.Update(func(tx *nutsdb.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = tx.Put([]byte("l"), genesis.Hash)

			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := Blockchain{tip, db}

	return &Blockchain{db}
}
