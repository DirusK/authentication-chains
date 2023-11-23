/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package cipher

import (
	"errors"
)

var (
	ErrFailedDecode          = errors.New("failed to decode PEM block containing public key")
	ErrFailedParsePublicKey  = errors.New("failed to parse encoded public key")
	ErrFailedParsePrivateKey = errors.New("failed to parse encoded private key")
	ErrDARVerification       = errors.New("failed to verify dar signature")
)
