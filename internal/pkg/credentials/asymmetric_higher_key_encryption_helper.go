package credentials

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"math/big"
)

// KeyLengths represents the lengths of encryption keys
type KeyLengths byte

const (
	// KeyLength32 represents a key length of 32 bytes
	KeyLength32 KeyLengths = iota
	// KeyLength64 represents a key length of 64 bytes
	KeyLength64
)

const (
	KeyLengthsPrefix = 2
	HmacKeySizeBytes = 64
	AesKeySizeBytes  = 32
)

// AsymmetricHigherKeyEncrypt encrypts plaintext using RSA with a higher key size
func AsymmetricHigherKeyEncrypt(plainTextBytes, modulusBytes, exponentBytes []byte) (string, error) {
	// Generate ephemeral keys for encryption (32 bytes), hmac (64 bytes)
	keyEnc, err := getRandomBytes(AesKeySizeBytes)
	if err != nil {
		return "", err
	}

	keyMac, err := getRandomBytes(HmacKeySizeBytes)
	if err != nil {
		return "", err
	}

	// Encrypt message using ephemeral keys and Authenticated Encryption
	ciphertext, err := AuthenticatedEncrypt(keyEnc, keyMac, plainTextBytes)
	if err != nil {
		return "", err
	}

	// Encrypt ephemeral keys using RSA
	keys := make([]byte, KeyLengthsPrefix+len(keyEnc)+len(keyMac))

	// Prefixing length of Keys. Symmetric Key length followed by HMAC key length
	keys[0] = byte(KeyLength32)
	keys[1] = byte(KeyLength64)

	copy(keys[2:], keyEnc)
	copy(keys[2+len(keyEnc):], keyMac)

	// Create RSA public key from modulus and exponent
	pub := &rsa.PublicKey{
		N: new(big.Int).SetBytes(modulusBytes),
		E: int(new(big.Int).SetBytes(exponentBytes).Int64()),
	}

	// Encrypt keys with RSA-OAEP-SHA256
	encryptedKeys, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, keys, nil)
	if err != nil {
		return "", err
	}

	// Prepare final payload
	return base64.StdEncoding.EncodeToString(encryptedKeys) + base64.StdEncoding.EncodeToString(ciphertext), nil
}

// getRandomBytes generates cryptographically secure random bytes
func getRandomBytes(size int) ([]byte, error) {
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
