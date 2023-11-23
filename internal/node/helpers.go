/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"bytes"
	"fmt"

	"github.com/nutsdb/nutsdb"
	"google.golang.org/protobuf/proto"

	"authentication-chains/internal/types"
)

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

// bucketAuthTableLevel returns the name of the bucket that will store authentication table by level.
func bucketAuthTableLevel(level uint) string {
	return fmt.Sprintf("%s %d", types.BucketAuthenticationTable, level)
}
