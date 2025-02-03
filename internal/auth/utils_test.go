// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package auth_test

import (
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/auth"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func TestUnit_ConvertFileToBase64(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		fileContent string
		expected    string
		expectError bool
	}{
		"valid file": {
			fileContent: "test content",
			expected:    base64.StdEncoding.EncodeToString([]byte("test content")),
			expectError: false,
		},
		"empty file": {
			fileContent: "",
			expected:    base64.StdEncoding.EncodeToString([]byte("")),
			expectError: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Create a temporary file
			tmpfile, err := os.CreateTemp(t.TempDir(), ".testfile-*.tmp")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			// Write content to the file if not expecting a read error
			if !testCase.expectError {
				if _, err := tmpfile.WriteString(testCase.fileContent); err != nil {
					t.Fatal(err)
				}
			}

			// Close the file
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			// Run the function
			result, err := auth.ConvertFileToBase64(tmpfile.Name())

			// Check for expected error
			if testCase.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.expected, result, "they should be equal")
			}
		})
	}
}

func TestUnit_ConvertBase64ToCert(t *testing.T) {
	certPass := testhelp.RandomName()

	t.Parallel()

	testCases := map[string]struct {
		b64           string
		password      string
		expectedCerts []*x509.Certificate
		expectedKey   crypto.PrivateKey
		expectError   bool
	}{
		"valid cert with password": {
			b64:         testhelp.RandomP12CertB64(certPass),
			password:    certPass,
			expectError: false,
		},
		"valid cert without password": {
			b64:         testhelp.RandomP12CertB64(""),
			password:    "",
			expectError: false,
		},
		"invalid cert": {
			b64:         "invalid base64",
			password:    certPass,
			expectError: true,
		},
		"empty cert": {
			b64:         "",
			password:    "",
			expectError: true,
		},
		"invalid password": {
			b64:         testhelp.RandomP12CertB64(certPass),
			password:    testhelp.RandomName(),
			expectError: true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, _, err := auth.ConvertBase64ToCert(testCase.b64, testCase.password)
			if testCase.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
