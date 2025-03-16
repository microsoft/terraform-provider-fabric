// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package transforms_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
)

func TestUnit_IsJSON(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "valid_json_object",
			content:  `{"name": "test", "value": 123}`,
			expected: true,
		},
		{
			name:     "valid_json_array",
			content:  `[1, 2, 3, 4]`,
			expected: true,
		},
		{
			name:     "invalid_json",
			content:  `{"name": "test", "value": 123`,
			expected: false,
		},
		{
			name:     "plain_text",
			content:  "This is not JSON",
			expected: false,
		},
		{
			name:     "empty_string",
			content:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transforms.IsJSON(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnit_JSONNormalize(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expected    string
		expectError bool
	}{
		{
			name:        "valid_json",
			content:     `{"name":"test", "value": 123}`,
			expected:    `{"name":"test","value":123}`,
			expectError: false,
		},
		{
			name:        "invalid_json",
			content:     `{"name":"test", "value": 123`,
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty_string",
			content:     "",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := transforms.JSONNormalize(tt.content)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}

	t.Run("nil_content", func(t *testing.T) {
		var nilString string
		result, err := transforms.JSONNormalize(nilString)
		assert.Equal(t, "", result)
		require.ErrorIs(t, err, transforms.ErrEmptyJSON)
	})
}

func TestUnit_JSONNormalizePretty(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name:        "valid_json",
			content:     `{"name":"test","value":123}`,
			expectError: false,
		},
		{
			name:        "invalid_json",
			content:     `{"name":"test", "value": 123`,
			expectError: true,
		},
		{
			name:        "empty_string",
			content:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := transforms.JSONNormalizePretty(tt.content)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify it's properly formatted JSON
				var obj any
				jsonErr := json.Unmarshal([]byte(result), &obj)
				require.NoError(t, jsonErr, "Result should be valid JSON")

				// Check that it contains newlines and spaces (pretty print)
				assert.Contains(t, result, "\n", "Pretty JSON should contain newlines")
				assert.Contains(t, result, "  ", "Pretty JSON should contain indentation")
			}
		})
	}

	t.Run("nil_content", func(t *testing.T) {
		var nilString string
		result, err := transforms.JSONNormalizePretty(nilString)
		assert.Equal(t, "", result)
		require.ErrorIs(t, err, transforms.ErrEmptyJSON)
	})
}

func TestUnit_JSONBase64GzipEncode(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name:        "valid_json",
			content:     `{"name":"test","value":123}`,
			expectError: false,
		},
		{
			name:        "invalid_json",
			content:     `{"name":"test", "value": 123`,
			expectError: true,
		},
		{
			name:        "empty_string",
			content:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := transforms.JSONBase64GzipEncode(tt.content)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify it's a base64 string
				_, decodeErr := base64.StdEncoding.DecodeString(result)
				require.NoError(t, decodeErr, "Result should be valid base64")

				// The result should not be equal to the original content
				assert.NotEqual(t, tt.content, result, "Content should be encoded")
			}
		})
	}
}

func TestUnit_JSONBase64GzipDecode(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name:        "valid_encoded_json",
			content:     "", // Will be populated in setup
			expectError: false,
		},
		{
			name:        "invalid_base64",
			content:     "this is not valid base64!",
			expectError: true,
		},
		{
			name:        "valid_base64_invalid_gzip",
			content:     "", // Will be populated in setup
			expectError: true,
		},
		{
			name:        "empty_string",
			content:     "",
			expectError: true,
		},
	}

	// Setup test data
	jsonContent := `{"name":"test","value":123}`
	encodedJSON, err := transforms.JSONBase64GzipEncode(jsonContent)
	require.NoError(t, err, "Failed to encode test JSON content")

	tests[0].content = encodedJSON

	invalidGzip := base64.StdEncoding.EncodeToString([]byte("this is not gzipped"))
	tests[2].content = invalidGzip

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := transforms.JSONBase64GzipDecode(tt.content)

			if tt.expectError {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result, "Decoded result should not be nil")

			// Try to convert result to map to verify it's a JSON object
			jsonMap, ok := result.(map[string]any)
			assert.True(t, ok, "Result should be a JSON object")
			assert.Contains(t, jsonMap, "name")
			assert.Contains(t, jsonMap, "value")
		})
	}
}

func TestUnit_Base64Decode(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expected    string
		expectError bool
	}{
		{
			name:        "valid_base64",
			content:     "SGVsbG8gV29ybGQ=", // "Hello World"
			expected:    "Hello World",
			expectError: false,
		},
		{
			name:        "invalid_base64",
			content:     "this is not valid base64!",
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty_string",
			content:     "",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := transforms.Base64Decode(tt.content)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestUnit_Base64Encode(t *testing.T) {
	tests := []struct {
		name     string
		content  any
		expected string
	}{
		{
			name:     "string_content",
			content:  "Hello World",
			expected: "SGVsbG8gV29ybGQ=",
		},
		{
			name:     "byte_content",
			content:  []byte("Hello World"),
			expected: "SGVsbG8gV29ybGQ=",
		},
		{
			name:     "empty_string",
			content:  "",
			expected: "",
		},
		{
			name:     "empty_bytes",
			content:  []byte{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			var err error

			switch v := tt.content.(type) {
			case string:
				result, err = transforms.Base64Encode(v)
			case []byte:
				result, err = transforms.Base64Encode(v)
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnit_Base64GzipEncode(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name:        "valid_content",
			content:     "Hello World",
			expectError: false,
		},
		{
			name:        "empty_string",
			content:     "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := transforms.Base64GzipEncode(tt.content)

			require.NoError(t, err)

			if tt.content != "" {
				assert.NotEmpty(t, result, "Encoded result should not be empty")
				assert.NotEqual(t, tt.content, result, "Content should be encoded")

				// Verify we can decode it back
				decoded, decodeErr := transforms.Base64GzipDecode(result)
				require.NoError(t, decodeErr)
				assert.Equal(t, tt.content, decoded)
			} else {
				assert.Empty(t, result, "Empty content should result in empty encoded string")
			}
		})
	}
}

func TestUnit_Base64GzipDecode(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expected    string
		expectError bool
	}{
		{
			name:        "valid_encoded_content",
			content:     "", // Will be populated in setup
			expected:    "Hello World",
			expectError: false,
		},
		{
			name:        "invalid_base64",
			content:     "this is not valid base64!",
			expected:    "",
			expectError: true,
		},
		{
			name:        "valid_base64_invalid_gzip",
			content:     "", // Will be populated in setup
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty_string",
			content:     "",
			expected:    "",
			expectError: false,
		},
	}

	// Setup test data
	encodedContent, err := transforms.Base64GzipEncode("Hello World")
	require.NoError(t, err, "Failed to encode test content")

	tests[0].content = encodedContent

	invalidGzip := base64.StdEncoding.EncodeToString([]byte("this is not gzipped"))
	tests[2].content = invalidGzip

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := transforms.Base64GzipDecode(tt.content)

			if tt.expectError {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
