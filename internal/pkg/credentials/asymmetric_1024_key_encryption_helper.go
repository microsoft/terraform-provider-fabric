package credentials

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"math/big"
	"time"
)

const (
	SegmentLength   = 85
	EncryptedLength = 128
	MaxAttempts     = 3
)

// Asymmetric1024KeyEncrypt encrypts plaintext using RSA-1024
func Asymmetric1024KeyEncrypt(plainTextBytes, modulusBytes, exponentBytes []byte) (string, error) {
	// Split the message into different segments, each segment's length is 85
	hasIncompleteSegment := len(plainTextBytes)%SegmentLength != 0

	segmentNumber := len(plainTextBytes) / SegmentLength
	if hasIncompleteSegment {
		segmentNumber++
	}

	encryptedBytes := make([]byte, segmentNumber*EncryptedLength)

	for i := range segmentNumber {
		lengthToCopy := SegmentLength
		if i == segmentNumber-1 && hasIncompleteSegment {
			lengthToCopy = len(plainTextBytes) % SegmentLength
		}

		segment := make([]byte, lengthToCopy)
		copy(segment, plainTextBytes[i*SegmentLength:i*SegmentLength+lengthToCopy])

		segmentEncryptedResult, err := encryptSegment(modulusBytes, exponentBytes, segment)
		if err != nil {
			return "", err
		}

		copy(encryptedBytes[i*EncryptedLength:], segmentEncryptedResult)
	}

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

// encryptSegment encrypts a segment with RSA
func encryptSegment(modulus, exponent, data []byte) ([]byte, error) {
	if data == nil {
		return nil, errors.New("data is required")
	}

	if len(data) == 0 {
		return data, nil
	}

	// Create RSA public key from modulus and exponent
	pub := &rsa.PublicKey{
		N: new(big.Int).SetBytes(modulus),
		E: int(new(big.Int).SetBytes(exponent).Int64()),
	}

	// Try encryption with retries
	for attempt := range MaxAttempts {
		encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
		if err == nil {
			return encryptedBytes, nil
		}

		// Sleep and retry on error
		if attempt < MaxAttempts-1 {
			time.Sleep(50 * time.Millisecond)
		}
	}

	return nil, errors.New("encryption failed after maximum attempts")
}
