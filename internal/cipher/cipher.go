/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package cipher

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"fmt"

	"google.golang.org/protobuf/proto"

	"authentication-chains/internal/types"
)

//go:generate ifacemaker -f cipher.go -s cipher -p cipher -i Cipher -y "Cipher - describe an interface for working with crypto operations."

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

// cipher is an implementation of the Cipher interface
type cipher struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func New() Cipher {
	// The GenerateKey method takes in a reader that returns random bits, and the number of bits
	privateKey, err := rsa.GenerateKey(rand.Reader, bytes256)
	if err != nil {
		panic(err)
	}

	return cipher{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}
}

// PublicKey returns the public rsa key.
func (c cipher) PublicKey() *rsa.PublicKey {
	return c.publicKey
}

// SerializePublicKey serializes the public key.
func (c cipher) SerializePublicKey() []byte {
	publicKeyBytes := x509.MarshalPKCS1PublicKey(c.publicKey)

	b := pem.Block{
		Type:  typeRSAPublicKey,
		Bytes: publicKeyBytes,
	}

	return pem.EncodeToMemory(&b)
}

// SerializePrivateKey serializes the private key.
func (c cipher) SerializePrivateKey() []byte {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(c.privateKey)

	b := pem.Block{
		Type:  typeRSAPrivateKey,
		Bytes: privateKeyBytes,
	}

	return pem.EncodeToMemory(&b)
}

// PrivateKey returns the private rsa key.
func (c cipher) PrivateKey() *rsa.PrivateKey {
	return c.privateKey
}

// Encrypt encrypts the given plain text using the public key.
func (c cipher) Encrypt(plainText []byte) ([]byte, error) {
	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, c.publicKey, plainText, nil)
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}

// Decrypt decrypts the given cipher text using the private key.
func (c cipher) Decrypt(cipherText []byte) ([]byte, error) {
	plainText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, c.privateKey, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

// Sign signs the given data using the private key.
func (c cipher) Sign(data []byte) ([]byte, error) {
	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message
	return rsa.SignPSS(rand.Reader, c.privateKey, crypto.SHA256, Hash(data), nil)
}

// VerifySignature verifies the given signature against the given data using the public key.
func (c cipher) VerifySignature(signature []byte, data []byte) error {
	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now
	return rsa.VerifyPSS(c.publicKey, crypto.SHA256, Hash(data), signature, nil)
}

// Serialize serializes the Cipher into bytes.
func (c cipher) Serialize() []byte {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)

	if err := encoder.Encode(c); err != nil {
		panic("cipher serialization error: " + err.Error())
	}

	return buffer.Bytes()
}

// SignDar signs the given DeviceAuthenticationRequest.
func (c cipher) SignDar(dar *types.DeviceAuthenticationRequest) error {
	data, err := proto.Marshal(dar)
	if err != nil {
		return fmt.Errorf("failed to marshal dar: %w", err)
	}

	dar.Signature, err = c.Sign(data)
	if err != nil {
		return fmt.Errorf("failed to sign dar: %w", err)
	}

	return nil
}

// EncryptContent encrypts the given Content.
func (c cipher) EncryptContent(content *types.Content) ([]byte, error) {
	data, err := proto.Marshal(content)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal content: %w", err)
	}

	return c.Encrypt(data)
}
