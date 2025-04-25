package credentials

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"
)

// AeCipher represents the available authenticated encryption cipher options
type AeCipher byte

const (
	// Aes256CbcPkcs7 is AES 256-bit cipher in CBC mode with PKCS#7 padding
	Aes256CbcPkcs7 AeCipher = iota
)

// AeMac represents the available message authentication code (MAC) algorithms
type AeMac byte

const (
	// HMACSHA256 uses SHA-256 hashing algorithm
	HMACSHA256 AeMac = iota
	// HMACSHA384 uses SHA-384 hashing algorithm
	HMACSHA384
	// HMACSHA512 uses SHA-512 hashing algorithm
	HMACSHA512
)

var (
	// Default algorithms
	aeCipher         = Aes256CbcPkcs7
	aeMac            = HMACSHA256
	algorithmChoices = []byte{byte(aeCipher), byte(aeMac)}
)

// Encrypt encrypts a message using the specified keys
func AuthenticatedEncrypt(keyEnc, keyMac, message []byte) ([]byte, error) {
	if keyEnc == nil {
		return nil, errors.New("keyEnc is required")
	}
	if keyMac == nil {
		return nil, errors.New("keyMac is required")
	}
	if len(keyEnc) < 32 {
		return nil, errors.New("encryption key must be at least 256 bits (32 bytes)")
	}
	if len(keyMac) < 32 {
		return nil, errors.New("mac key must be at least 256 bits (32 bytes)")
	}
	if message == nil {
		return nil, errors.New("message is required")
	}

	// Create cipher
	block, err := aes.NewCipher(keyEnc[:32])
	if err != nil {
		return nil, err
	}

	// Generate random IV
	iv := make([]byte, block.BlockSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	// Encrypt the message
	cipherText, err := encrypt(block, iv, message)
	if err != nil {
		return nil, err
	}

	// Calculate MAC
	tagGenerator, err := getMac(aeMac, keyMac)
	if err != nil {
		return nil, err
	}

	// The IV and ciphertext both need to be included in the MAC to prevent tampering
	tagData := bytes.Buffer{}
	tagData.Write(algorithmChoices)
	tagData.Write(iv)
	tagData.Write(cipherText)

	tagGenerator.Write(tagData.Bytes())
	tag := tagGenerator.Sum(nil)

	// Build the final result
	result := bytes.Buffer{}
	result.Write(algorithmChoices)
	result.Write(tag)
	result.Write(iv)
	result.Write(cipherText)

	return result.Bytes(), nil
}

// Helper function to encrypt with AES-CBC
func encrypt(block cipher.Block, iv, plaintext []byte) ([]byte, error) {
	// Apply PKCS#7 padding
	padding := block.BlockSize() - (len(plaintext) % block.BlockSize())
	padtext := make([]byte, len(plaintext)+padding)
	copy(padtext, plaintext)
	for i := len(plaintext); i < len(padtext); i++ {
		padtext[i] = byte(padding)
	}

	// Encrypt
	ciphertext := make([]byte, len(padtext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padtext)

	return ciphertext, nil
}

// Helper function to get MAC algorithm
func getMac(macType AeMac, key []byte) (hash.Hash, error) {
	switch macType {
	case HMACSHA256:
		return hmac.New(sha256.New, key), nil
	case HMACSHA384:
		return hmac.New(sha512.New384, key), nil
	case HMACSHA512:
		return hmac.New(sha512.New, key), nil
	default:
		return nil, errors.New("invalid MAC algorithm")
	}
}
