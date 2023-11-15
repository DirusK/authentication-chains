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
	"encoding/gob"
)

//go:generate ifacemaker -f cipher.go -s cipher -p cipher -i Cipher -y "Cipher - describe an interface for working with crypto operations."

const (
	// bytes512 is 4096 bits
	bytes512 = 4096
	// bytes256 is 2048 bits
	bytes256 = 2048
)

// cipher is an implementation of the Cipher interface
type cipher struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func New() Cipher {
	// The GenerateKey method takes in a reader that returns random bits, and the number of bits
	privateKey, err := rsa.GenerateKey(rand.Reader, bytes256)
	if err != nil {
		panic(err)
	}

	return cipher{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}
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

// Sign signs the given data using the private key.
func (c cipher) Sign(data []byte) ([]byte, error) {
	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHash := sha256.New()
	if _, err := msgHash.Write(data); err != nil {
		return nil, err
	}

	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message
	signature, err := rsa.SignPSS(rand.Reader, c.PrivateKey, crypto.SHA256, msgHash.Sum(nil), nil)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// Verify verifies the given signature against the given data using the public key.
func (c cipher) Verify(signature []byte, data []byte) error {
	// Before verifying the signature, we need to hash our message
	// The hash is what we actually verify
	msgHash := sha256.New()
	if _, err := msgHash.Write(data); err != nil {
		return err
	}

	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now
	if err := rsa.VerifyPSS(c.PublicKey, crypto.SHA256, msgHash.Sum(nil), signature, nil); err != nil {
		return err
	}

	return nil
}

func (c cipher) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(c.PrivateKey); err != nil {
		return nil, err
	}
	if err := encoder.Encode(c.PublicKey); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (c cipher) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(&c.PrivateKey); err != nil {
		return err
	}
	if err := decoder.Decode(&c.PublicKey); err != nil {
		return err
	}

	return nil
}
