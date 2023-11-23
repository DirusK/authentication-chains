/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package blockchain

import (
	"encoding/binary"
)

func uint64ToBytes(num uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, num)
	return bytes
}
