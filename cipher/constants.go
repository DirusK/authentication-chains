/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package cipher

import (
	"errors"
)

const (
	// bytes512 is 4096 bits
	bytes512 = 4096
	// bytes256 is 2048 bits
	bytes256 = 2048
)

// Types for PEM encoding
const (
	typeRSAPublicKey  = "RSA PUBLIC KEY"
	typeRSAPrivateKey = "RSA PRIVATE KEY"
)

var (
	ErrFailedDecode          = errors.New("failed to decode PEM block containing public key")
	ErrFailedParsePublicKey  = errors.New("failed to parse encoded public key")
	ErrFailedParsePrivateKey = errors.New("failed to parse encoded private key")
)
