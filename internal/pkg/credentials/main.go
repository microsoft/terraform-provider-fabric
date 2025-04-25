package credentials

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

func example() {
	// Create gateway public key (you would get this from your actual system)
	publicKey := fabcore.PublicKey{
		Modulus:  to.Ptr("your-modulus-base64"),
		Exponent: to.Ptr("your-exponent-base64"),
	}

	// Create credentials
	credentials, err := NewBasicCredentials("username", "password")
	if err != nil {
		log.Fatalf("Failed to create credentials: %v", err)
	}

	// Encrypt credentials
	encryptedCredentials, err := EncryptCredentials(*credentials, publicKey)
	if err != nil {
		log.Fatalf("Failed to encrypt credentials: %v", err)
	}

	fmt.Printf("Encrypted credentials: %s\n", encryptedCredentials)
}
