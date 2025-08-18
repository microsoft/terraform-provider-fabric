// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package testhelp

import (
	"crypto/sha1" //nolint:gosec
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
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

	return acctest.RandStringFromCharSet(size, acctest.CharSetAlpha+strings.ToUpper(acctest.CharSetAlpha))
}

func RandomTime(min, max time.Time) time.Time {
	if !min.Before(max) {
		panic("min must be before max")
	}

	// Duration between min and max
	delta := max.Sub(min)

	// Random duration offset
	offset := time.Duration(rand.Int64N(int64(delta)))

	return min.Add(offset)
}

// RandomTimeDefault returns a random time between Unix epoch and now.
func RandomTimeDefault() time.Time {
	return RandomTime(time.Unix(0, 0), time.Now())
}

// RandomIntRange returns a random integer between minInt (inclusive) and maxInt (exclusive).
func RandomIntRange[T ~int | ~int8 | ~int16 | ~int32 | ~int64](minInt, maxInt T) T {
	if minInt >= maxInt {
		panic(fmt.Sprintf("minInt %d must be less than maxInt %d", minInt, maxInt)) // lintignore:R009
	}

	// Generate a random integer in the range [minInt, maxInt)
	return rand.N(maxInt-minInt) + minInt // #nosec G404
}

func RandomBool() bool {
	return RandomIntRange(0, 2) == 1
}

func RandomUUID() string {
	result, err := uuid.GenerateUUID()
	if err != nil {
		panic("failed to generate UUID: " + err.Error()) // lintignore:R009
	}

	return result
}

func RandomElement[T any](elements []T) T {
	return elements[RandomIntRange(0, len(elements))]
}

func RandomURI() string {
	return fmt.Sprintf("https://%s.com", strings.ToLower(RandomName()))
}

func RandomSHA1() string {
	hash := sha1.Sum([]byte(RandomUUID())) //nolint:gosec

	return hex.EncodeToString(hash[:])
}

func RandomP12CertB64(password string) string {
	cert := RandomP12Cert(password)

	return base64.StdEncoding.EncodeToString(cert)
}

func RandomP12Cert(password string) []byte {
	certPEM, privateKeyPEM, err := acctest.RandTLSCert("test")
	if err != nil {
		panic("failed to generate random TLS cert: " + err.Error()) // lintignore:R009
	}

	p12, err := createP12Bundle(certPEM, privateKeyPEM, password)
	if err != nil {
		panic("failed to create p12 bundle: " + err.Error()) // lintignore:R009
	}

	return p12
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

func TFDataSource(providerName, typeName, dataSourceName string) (fqn, header string) { //nolint:nonamedreturns
	fqn = DataSourceFQN(providerName, typeName, dataSourceName)
	header = at.DataSourceHeader(TypeName(providerName, typeName), dataSourceName)

	return fqn, header
}

func TFResource(providerName, typeName, resourceName string) (fqn, header string) { //nolint:nonamedreturns
	fqn = ResourceFQN(providerName, typeName, resourceName)
	header = at.ResourceHeader(TypeName(providerName, typeName), resourceName)

	return fqn, header
}

func TFEphemeral(providerName, typeName, ephemeralResourceName string) (fqn, header string) { //nolint:nonamedreturns
	fqn = EphemeralResourceFQN(providerName, typeName, ephemeralResourceName)
	header = EphemeralResourceHeader(TypeName(providerName, typeName), ephemeralResourceName)

	return fqn, header
}

func TFEphemeralEcho(ephemeralResourceFQN string) (config, fqn string) { //nolint:nonamedreturns
	fqn = "echo.test"

	// lintignore:AT004
	config = fmt.Sprintf(`
					provider "echo" {
						data = %[1]s
					}

					resource "echo" "test" {}
				`, ephemeralResourceFQN)

	return fqn, config
}

func EphemeralResourceHeader(ephemeralResourceType, ephemeralResourceName string) string {
	const f = `ephemeral %q %q`

	return fmt.Sprintf(f, ephemeralResourceType, ephemeralResourceName)
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

// Helper function to create a ephemeral resource FQN.
func EphemeralResourceFQN(providerName, typeName, ephemeralResource string) string {
	return fmt.Sprintf("ephemeral.%s.%s", TypeName(providerName, typeName), ephemeralResource)
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
