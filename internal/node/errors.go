/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"errors"
)

var (
	ErrBlockValidation  = errors.New("block validation failed")
	ErrVerification     = errors.New("verification failed")
	ErrInvalidDAR       = errors.New("invalid device authentication request")
	ErrWaitBlockTimeout = errors.New("wait block timeout ")
	ErrBlockHasNotMined = errors.New("block has not mined")
	ErrBlockHasNoDAR    = errors.New("block has no device authentication request")
)
