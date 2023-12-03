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

func FromStringPrivateKey(privateKeyString string) (Cipher, error) {
	privateKey, err := DeserializePrivateKey([]byte(privateKeyString))
	if err != nil {
		return nil, err
	}

	return &cipher{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

// Hash hashes the given data using SHA256
func Hash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// HashBlock without a hash field.
func HashBlock(block *types.Block) ([]byte, error) {
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

// VerifyDAR verifies the given DeviceAuthenticationRequest.
func VerifyDAR(dar *types.DeviceAuthenticationRequest) error {
	copyDar := &types.DeviceAuthenticationRequest{
		DeviceId:      dar.DeviceId,
		ClusterHeadId: dar.ClusterHeadId,
	}

	pubKey, err := DeserializePublicKey(copyDar.DeviceId)
	if err != nil {
		return fmt.Errorf("failed to deserialize public key: %w", err)
	}

	data, err := proto.Marshal(copyDar)
	if err != nil {
		return fmt.Errorf("failed to marshal dar: %w", err)
	}

	if err = VerifySignature(pubKey, dar.Signature, data); err != nil {
		return fmt.Errorf("failed to verify dar signature: %w", ErrDARVerification)
	}

	return nil
}

// VerifySignature verifies the given signature against the given data using the public key.
func VerifySignature(publicKey *rsa.PublicKey, signature, data []byte) error {
	return rsa.VerifyPSS(publicKey, crypto.SHA256, Hash(data), signature, nil)
}

func SerializePublicKey(publicKey *rsa.PublicKey) []byte {
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)

	b := pem.Block{
		Type:  typeRSAPublicKey,
		Bytes: publicKeyBytes,
	}

	return pem.EncodeToMemory(&b)
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

func SerializePrivateKey(privateKey *rsa.PrivateKey) []byte {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	b := pem.Block{
		Type:  typeRSAPrivateKey,
		Bytes: privateKeyBytes,
	}

	return pem.EncodeToMemory(&b)
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

func EncryptContent(pubKey *rsa.PublicKey, content *types.Content) ([]byte, error) {
	data, err := proto.Marshal(content)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal content: %w", err)
	}

	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, data, nil)
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}
