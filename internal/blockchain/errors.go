/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package blockchain

import (
	"errors"
)

var (
	ErrBlockValidation = errors.New("block validation failed")
	ErrEmptyMemPool    = errors.New("mempool is empty")
)
