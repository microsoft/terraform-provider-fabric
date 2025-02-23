// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package testhelp

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"software.sslmate.com/src/go-pkcs12"
)

func RandomName(length ...int) string {
	var size int
	if len(length) == 0 || len(length) < 1 {
		size = 20 // default size
	} else {
		size = length[0]
	}

	return acctest.RandStringFromCharSet(size, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func RandomUUID() string {
	result, _ := uuid.GenerateUUID()

	return result
}

func RandomURI() string {
	return fmt.Sprintf("https://%s.com", strings.ToLower(RandomName()))
}

func RandomP12CertB64(password string) string {
	cert := RandomP12Cert(password)

	return base64.StdEncoding.EncodeToString(cert)
}

func RandomP12Cert(password string) []byte {
	certPEM, privateKeyPEM, _ := acctest.RandTLSCert("test")
	p12, _ := createP12Bundle(certPEM, privateKeyPEM, password)

	return p12
}

func RandomInt32Max(max int32) int32 {
	return rand.Int31n(max)
}

func RandomInt32Range(min int32, max int32) int32 {
	return min + rand.Int31n(max-min+1)
}

func RandomElement[T any](slice []T) T {
	return slice[rand.Intn(len(slice))]
}

func createP12Bundle(certPEMStr, privateKeyPEMStr, password string) ([]byte, error) {
	// Decode the private key PEM block
	block, _ := pem.Decode([]byte(privateKeyPEMStr))
	if block == nil {
		return nil, errors.New("failed to parse private key PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %s", err.Error())
	}

	// Decode the certificate PEM block
	block, _ = pem.Decode([]byte(certPEMStr))
	if block == nil {
		return nil, errors.New("failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %s", err.Error())
	}

	p12Bytes, err := pkcs12.Modern2023.Encode(privateKey, cert, nil, password)
	if err != nil {
		return nil, fmt.Errorf("unable to encode p12: %s", err.Error())
	}

	return p12Bytes, nil
}

// Helper function to create a base type name.
func TypeName(providerName, typeName string) string {
	return fmt.Sprintf("%s_%s", providerName, typeName)
}

// Helper function to create a resource FQN.
func ResourceFQN(providerName, typeName, resourceName string) string {
	return fmt.Sprintf("%s.%s", TypeName(providerName, typeName), resourceName)
}

// Helper function to create a data source FQN.
func DataSourceFQN(providerName, typeName, dataSourceName string) string {
	return fmt.Sprintf("data.%s.%s", TypeName(providerName, typeName), dataSourceName)
}

// Helper function to create a function Header.
func FunctionHeader(providerName, functionName string) string {
	return fmt.Sprintf(`provider::%s::%s`, providerName, functionName)
}

// Helper function to create a reference by FQN.
func RefByFQN(objectFQN, path string) string {
	return fmt.Sprintf("${%s.%s}", objectFQN, path)
}

// Helper function to deep copy a map.
func CopyMap(src map[string]any) map[string]any {
	cp := make(map[string]any)

	for k, v := range src {
		switch v := v.(type) {
		case map[string]any:
			cp[k] = CopyMap(v)
		default:
			cp[k] = v
		}
	}

	return cp
}
