package cipher

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"fmt"
)

// Deserialize deserializes the given data into a Cipher.
func Deserialize(data []byte) (Cipher, error) {
	var c cipher

	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	if err := decoder.Decode(&c); err != nil {
		return nil, fmt.Errorf("cipher deserialization error: %w", err)
	}

	return c, nil
}

// Hash hashes the given data using SHA256
func Hash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// VerifySignature verifies the given signature against the given data using the public key.
func VerifySignature(publicKey *rsa.PublicKey, signature, data []byte) error {
	return rsa.VerifyPSS(publicKey, crypto.SHA256, Hash(data), signature, nil)
}

// DeserializePublicKey deserializes a public key from bytes.
func DeserializePublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != typeRSAPublicKey {
		return nil, ErrFailedDecode
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedParsePublicKey, err)
	}

	return publicKey, nil
}

// DeserializePrivateKey deserializes a private key from bytes.
func DeserializePrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != typeRSAPrivateKey {
		return nil, ErrFailedDecode
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedParsePrivateKey, err)
	}

	return privateKey, nil
}
