// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

func ConvertFileToBase64(path string) (string, error) {
	pfx, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("could not read PFX file: %w", err)
	}

	b64 := base64.StdEncoding.EncodeToString(pfx)

	return b64, nil
}

func ConvertBase64ToCert(b64, password string) ([]*x509.Certificate, crypto.PrivateKey, error) {
	pfx, err := convertBase64ToByte(b64)
	if err != nil {
		return nil, nil, err
	}

	certs, key, err := convertByteToCert(pfx, password)
	if err != nil {
		return nil, nil, err
	}

	return certs, key, nil
}

func convertBase64ToByte(b64 string) ([]byte, error) {
	if b64 == "" {
		return nil, errors.New("got empty base64 certificate data")
	}

	pfx, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return pfx, fmt.Errorf("could not decode base64 certificate data: %w", err)
	}

	return pfx, nil
}

func convertByteToCert(certData []byte, password string) ([]*x509.Certificate, crypto.PrivateKey, error) {
	var key crypto.PrivateKey

	key, cert, _, err := pkcs12.DecodeChain(certData, password)
	if err != nil {
		return nil, nil, err
	}

	if cert == nil {
		return nil, nil, errors.New("found no certificate")
	}

	certs := []*x509.Certificate{cert}

	return certs, key, nil
}
