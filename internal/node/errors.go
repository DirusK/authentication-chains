/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"errors"
)

var (
	ErrBlockValidation        = errors.New("block validation failed")
	ErrVerification           = errors.New("verification failed")
	ErrInvalidDAR             = errors.New("invalid device authentication request")
	ErrBlockHasNoDAR          = errors.New("block has no device authentication request")
	ErrInvalidMessageReceiver = errors.New("invalid message receiver")
	ErrNotFoundBlock          = errors.New("block not found")
)
