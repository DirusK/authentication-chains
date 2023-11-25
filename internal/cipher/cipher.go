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

	"github.com/nutsdb/nutsdb"
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
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func New(db *nutsdb.DB) (Cipher, error) {
	var c Cipher

	if db != nil {
		if err := db.View(func(tx *nutsdb.Tx) error {
			entry, err := tx.Get(types.BucketCipher, types.KeyCipher)
			if err != nil {
				return err
			}

			c, err = Deserialize(entry.Value)
			if err != nil {
				return err
			}

			return nil
		}); err == nil {
			return c, nil
		}
	}

	// The GenerateKey method takes in a reader that returns random bits, and the number of bits
	privateKey, err := rsa.GenerateKey(rand.Reader, bytes256)
	if err != nil {
		return nil, err
	}

	c = cipher{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}

	if db != nil {
		if err = db.Update(func(tx *nutsdb.Tx) error {
			return tx.Put(types.BucketCipher, types.KeyCipher, c.Serialize(), types.InfinityTTL)
		}); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// GetPublicKey returns the public rsa key.
func (c cipher) GetPublicKey() *rsa.PublicKey {
	return c.PublicKey
}

// SerializePublicKey serializes the public key.
func (c cipher) SerializePublicKey() []byte {
	publicKeyBytes := x509.MarshalPKCS1PublicKey(c.PublicKey)

	b := pem.Block{
		Type:  typeRSAPublicKey,
		Bytes: publicKeyBytes,
	}

	return pem.EncodeToMemory(&b)
}

func (c cipher) ToHexPublicKey() string {
	return fmt.Sprintf("%x", c.SerializePublicKey())
}

// SerializePrivateKey serializes the private key.
func (c cipher) SerializePrivateKey() []byte {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(c.PrivateKey)

	b := pem.Block{
		Type:  typeRSAPrivateKey,
		Bytes: privateKeyBytes,
	}

	return pem.EncodeToMemory(&b)
}

func (c cipher) ToHexPrivateKey() string {
	return fmt.Sprintf("%x", c.SerializePrivateKey())
}

// GetPrivateKey returns the private rsa key.
func (c cipher) GetPrivateKey() *rsa.PrivateKey {
	return c.PrivateKey
}

// Encrypt encrypts the given plain text using the public key.
func (c cipher) Encrypt(plainText []byte) ([]byte, error) {
	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, c.PublicKey, plainText, nil)
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}

// Decrypt decrypts the given cipher text using the private key.
func (c cipher) Decrypt(cipherText []byte) ([]byte, error) {
	plainText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, c.PrivateKey, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

// DecryptContent decrypts the given Content.
func (c cipher) DecryptContent(cipherText []byte) (*types.Content, error) {
	plainText, err := c.Decrypt(cipherText)
	if err != nil {
		return nil, err
	}

	content := new(types.Content)
	if err = proto.Unmarshal(plainText, content); err != nil {
		return nil, err
	}

	return content, nil
}

// Sign signs the given data using the private key.
func (c cipher) Sign(data []byte) ([]byte, error) {
	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message
	return rsa.SignPSS(rand.Reader, c.PrivateKey, crypto.SHA256, Hash(data), nil)
}

// VerifySignature verifies the given signature against the given data using the public key.
func (c cipher) VerifySignature(signature []byte, data []byte) error {
	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now
	return rsa.VerifyPSS(c.PublicKey, crypto.SHA256, Hash(data), signature, nil)
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

// SignDAR signs the given DeviceAuthenticationRequest.
func (c cipher) SignDAR(dar *types.DeviceAuthenticationRequest) error {
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

// VerifyDAR verifies the given DeviceAuthenticationRequest.
func (c cipher) VerifyDAR(dar *types.DeviceAuthenticationRequest) error {
	signature := dar.Signature
	dar.Signature = nil

	data, err := proto.Marshal(dar)
	if err != nil {
		return fmt.Errorf("failed to marshal dar: %w", err)
	}

	if err = c.VerifySignature(signature, data); err != nil {
		return fmt.Errorf("failed to verify dar signature: %w", err)
	}

	return nil
}

// HashBlock without a hash field.
func (c cipher) HashBlock(block *types.Block) ([]byte, error) {
	bc := &types.Block{
		Hash:      nil,
		PrevHash:  block.PrevHash,
		Index:     block.Index,
		Dar:       block.Dar,
		Timestamp: block.Timestamp,
	}

	data, err := proto.Marshal(bc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal block: %w", err)
	}

	result := Hash(data)

	return result, nil
}

// EncryptContent encrypts the given Content.
func (c cipher) EncryptContent(content *types.Content) ([]byte, error) {
	data, err := proto.Marshal(content)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal content: %w", err)
	}

	return c.Encrypt(data)
}
