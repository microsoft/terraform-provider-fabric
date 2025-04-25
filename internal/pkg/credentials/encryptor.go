package credentials

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

// EncryptCredentials encrypts credentials using RSA algorithm
func EncryptCredentials(credentialData Credentials, publicKey fabcore.PublicKey) (string, error) {
	if publicKey.Exponent == nil {
		return "", errors.New("publicKey.Exponent is required")
	}

	if publicKey.Modulus == nil {
		return "", errors.New("publicKey.Modulus is required")
	}

	plainTextBytes, err := json.Marshal(credentialData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal credentials: %v", err)
	}

	modulusBytes, err := base64.StdEncoding.DecodeString(*publicKey.Modulus)
	if err != nil {
		return "", fmt.Errorf("failed to decode modulus: %v", err)
	}

	exponentBytes, err := base64.StdEncoding.DecodeString(*publicKey.Exponent)
	if err != nil {
		return "", fmt.Errorf("failed to decode exponent: %v", err)
	}

	// Choose the encryption method based on key size
	if len(modulusBytes) == 128 {
		return Asymmetric1024KeyEncrypt(plainTextBytes, modulusBytes, exponentBytes)
	}
	return AsymmetricHigherKeyEncrypt(plainTextBytes, modulusBytes, exponentBytes)
}
