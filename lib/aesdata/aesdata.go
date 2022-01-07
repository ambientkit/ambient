// Package aesdata provide encryption/decryption of JSON data using AES.
package aesdata

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// Resource: https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes/

// EncryptedStorage represents a private key.
type EncryptedStorage struct {
	privatekey string
}

// NewEncryptedStorage returns an encrypted storage object.
func NewEncryptedStorage(privatekey string) *EncryptedStorage {
	return &EncryptedStorage{
		privatekey: privatekey,
	}
}

// Encrypt data using private key.
func (en *EncryptedStorage) Encrypt(data []byte) ([]byte, error) {
	// Convert key to byte array.
	key, err := hex.DecodeString(en.privatekey)
	if err != nil {
		return nil, err
	}

	// Create a new cipher block from the key.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Wrap the cipher block.
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Create a nonce.
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt the data.
	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

// Decrypt data using private key or return empty JSON array.
func (en *EncryptedStorage) Decrypt(enc []byte) ([]byte, error) {
	// Don't decrypt if there is no content.
	if len(enc) == 0 {
		return []byte("{}"), nil
	}

	// Convert key to byte array.
	key, err := hex.DecodeString(en.privatekey)
	if err != nil {
		return nil, err
	}

	// Create a new cipher block from the key.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Wrap the cipher block.
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Create a nonce.
	nonceSize := aesGCM.NonceSize()

	// Extract the nonce from the encrypted data.
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	// Decrypt the data.
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
