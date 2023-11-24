/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package cipher

import (
	"crypto/rsa"

	"authentication-chains/internal/types"
)

// Cipher - describe an interface for working with crypto operations.
type Cipher interface {
	// GetPublicKey returns the public rsa key.
	GetPublicKey() *rsa.PublicKey
	// SerializePublicKey serializes the public key.
	SerializePublicKey() []byte
	// SerializePrivateKey serializes the private key.
	SerializePrivateKey() []byte
	// GetPrivateKey returns the private rsa key.
	GetPrivateKey() *rsa.PrivateKey
	// Encrypt encrypts the given plain text using the public key.
	Encrypt(plainText []byte) ([]byte, error)
	// Decrypt decrypts the given cipher text using the private key.
	Decrypt(cipherText []byte) ([]byte, error)
	// Sign signs the given data using the private key.
	Sign(data []byte) ([]byte, error)
	// VerifySignature verifies the given signature against the given data using the public key.
	VerifySignature(signature []byte, data []byte) error
	// Serialize serializes the Cipher into bytes.
	Serialize() []byte
	// SignDAR signs the given DeviceAuthenticationRequest.
	SignDAR(dar *types.DeviceAuthenticationRequest) error
	// VerifyDAR verifies the given DeviceAuthenticationRequest.
	VerifyDAR(dar *types.DeviceAuthenticationRequest) error
	// EncryptContent encrypts the given Content.
	EncryptContent(content *types.Content) ([]byte, error)
	// HashBlock without a hash field.
	HashBlock(block *types.Block) ([]byte, error)
	// DecryptContent decrypts the given Content.
	DecryptContent(cipherText []byte) (*types.Content, error)
}
