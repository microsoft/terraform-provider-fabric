// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package transforms_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func TestUnit_SourceFileToPayload(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Create test files

	// Create a regular text file
	regularContent := "Hello World"
	regularFilePath := filepath.Join(tempDir, testhelp.RandomUUID()+".txt")

	err := os.WriteFile(regularFilePath, []byte(regularContent), 0o600)
	require.NoError(t, err, "Failed to write test file")

	// Create a JSON file
	jsonContent := `{"name":"test","value":123}`
	jsonFilePath := filepath.Join(tempDir, testhelp.RandomUUID()+".json")

	err = os.WriteFile(jsonFilePath, []byte(jsonContent), 0o600)
	require.NoError(t, err, "Failed to write test file")

	// Create a template file with regular content
	templateRegularContent := `Hello {{.Name}}!`
	templateRegularFilePath := filepath.Join(tempDir, testhelp.RandomUUID()+".tmpl")

	err = os.WriteFile(templateRegularFilePath, []byte(templateRegularContent), 0o600)
	require.NoError(t, err, "Failed to write test file")

	// Create a template file with json content
	templateJSONContent := `{"name": "{{.Name}}", "value": {{.Value}}}`
	templateJSONFilePath := filepath.Join(tempDir, testhelp.RandomUUID()+".tmpl")

	err = os.WriteFile(templateJSONFilePath, []byte(templateJSONContent), 0o600)
	require.NoError(t, err, "Failed to write test file")

	// Create a template file with invalid regular content
	invalidTemplateRegularContent := `Hello {{.Name!}}`
	invalidTemplateRegularFilePath := filepath.Join(tempDir, testhelp.RandomUUID()+".tmpl")

	err = os.WriteFile(invalidTemplateRegularFilePath, []byte(invalidTemplateRegularContent), 0o600)
	require.NoError(t, err, "Failed to write test file")

	// Create a template file with invalid JSON content
	invalidTemplateJSONContent := `{ "name": "{{.Name}}", "value": }`
	invalidTemplateJSONFilePath := filepath.Join(tempDir, testhelp.RandomUUID()+".tmpl")

	err = os.WriteFile(invalidTemplateJSONFilePath, []byte(invalidTemplateJSONContent), 0o600)
	require.NoError(t, err, "Failed to write test file")

	// Create a binary file (image)
	binaryFilePath := filepath.Join(tempDir, testhelp.RandomUUID()+".png")
	binaryData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG file signature

	err = os.WriteFile(binaryFilePath, binaryData, 0o600)
	require.NoError(t, err, "Failed to write binary file")

	nonExistentPath := filepath.Join(tempDir, testhelp.RandomUUID()+".txt")

	tests := []struct {
		name           string
		filePath       string
		tokens         map[string]string
		expectError    bool
		content        string
		processingMode string
	}{
		{
			name:           "regular_file_none",
			filePath:       regularFilePath,
			tokens:         nil,
			expectError:    false,
			content:        regularContent,
			processingMode: "None",
		},
		{
			name:           "json_file_none",
			filePath:       jsonFilePath,
			tokens:         nil,
			expectError:    false,
			content:        jsonContent,
			processingMode: "None",
		},
		{
			name:           "template_regular_file_none",
			filePath:       templateRegularFilePath,
			tokens:         map[string]string{"Name": "World"},
			expectError:    false,
			content:        templateRegularContent,
			processingMode: "None",
		},
		{
			name:           "template_json_file_none",
			filePath:       templateJSONFilePath,
			tokens:         map[string]string{"Name": "World", "Value": "123"},
			expectError:    false,
			content:        templateJSONContent,
			processingMode: "None",
		},
		{
			name:           "invalid_template_regular_none",
			filePath:       invalidTemplateRegularFilePath,
			tokens:         map[string]string{"Name": "World"},
			expectError:    false,
			content:        invalidTemplateRegularContent,
			processingMode: "None",
		},
		{
			name:           "invalid_template_json_none",
			filePath:       invalidTemplateJSONFilePath,
			tokens:         map[string]string{"Name": "World"},
			expectError:    false,
			content:        invalidTemplateJSONContent,
			processingMode: "None",
		},
		{
			name:           "non_existent_file_none",
			filePath:       nonExistentPath,
			tokens:         nil,
			expectError:    true,
			content:        "",
			processingMode: "None",
		},
		{
			name:           "binary_file_none",
			filePath:       binaryFilePath,
			tokens:         nil,
			expectError:    false,
			content:        "",
			processingMode: "None",
		},
		{
			name:           "regular_file",
			filePath:       regularFilePath,
			tokens:         nil,
			expectError:    false,
			content:        regularContent,
			processingMode: "GoTemplate",
		},
		{
			name:           "json_file",
			filePath:       jsonFilePath,
			tokens:         nil,
			expectError:    false,
			content:        jsonContent,
			processingMode: "GoTemplate",
		},
		{
			name:           "template_regular_file",
			filePath:       templateRegularFilePath,
			tokens:         map[string]string{"Name": "World"},
			expectError:    false,
			content:        "Hello World!",
			processingMode: "GoTemplate",
		},
		{
			name:           "template_json_file",
			filePath:       templateJSONFilePath,
			tokens:         map[string]string{"Name": "World", "Value": "123"},
			expectError:    false,
			content:        `{"name":"World","value":123}`,
			processingMode: "GoTemplate",
		},
		{
			name:           "invalid_template_regular",
			filePath:       invalidTemplateRegularFilePath,
			tokens:         map[string]string{"Name": "World"},
			expectError:    true,
			content:        "",
			processingMode: "GoTemplate",
		},
		{
			name:           "invalid_template_json",
			filePath:       invalidTemplateJSONFilePath,
			tokens:         map[string]string{"Name": "World"},
			expectError:    false,
			content:        `{ "name": "World", "value": }`,
			processingMode: "GoTemplate",
		},
		{
			name:           "non_existent_file",
			filePath:       nonExistentPath,
			tokens:         nil,
			expectError:    true,
			content:        "",
			processingMode: "GoTemplate",
		},
		{
			name:           "binary_file",
			filePath:       binaryFilePath,
			tokens:         nil,
			expectError:    false,
			content:        "",
			processingMode: "GoTemplate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentB64, contentSha256, diags := transforms.SourceFileToPayload(tt.filePath, tt.processingMode, tt.tokens, nil)
			errCount := diags.ErrorsCount()

			if tt.expectError {
				require.Positive(t, errCount, "Expected error diagnostics, got none")
			} else {
				assert.Equal(t, 0, errCount, "Unexpected error diagnostics: %v", diags)
				require.NotEmpty(t, contentB64, "Expected non-empty contentB64")
				require.NotEmpty(t, contentSha256, "Expected non-empty contentSha256")

				if tt.content != "" {
					// Decode the base64 content
					decodedContent, err := transforms.Base64Decode(contentB64)
					require.NoError(t, err)

					// Verify the decoded content matches the expected content
					assert.Equal(t, tt.content, decodedContent, "Decoded content does not match expected content")
				}
			}
		})
	}

	// Additional test for null/unknown tokens
	t.Run("null_tokens", func(t *testing.T) {
		var tokens map[string]string

		contentB64, contentSha256, diags := transforms.SourceFileToPayload(regularFilePath, "GoTemplate", tokens, nil)

		assert.False(t, diags.HasError(), "Unexpected error diagnostics with null tokens: %v", diags)
		require.NotEmpty(t, contentB64, "Expected non-empty contentB64")
		require.NotEmpty(t, contentSha256, "Expected non-empty contentSha256")
	})
}

func TestUnit_PayloadToGzip(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		isJSON      bool
		expectError bool
	}{
		{
			name:        "valid_json",
			content:     `{"name":"test","value":123}`,
			isJSON:      true,
			expectError: false,
		},
		{
			name:        "plain_text",
			content:     "Hello World",
			isJSON:      false,
			expectError: false,
		},
		{
			name:        "invalid_base64",
			content:     "this is not valid base64!",
			isJSON:      false,
			expectError: true,
		},
		{
			name:        "invalid_json",
			content:     `{"name":"test", "value": 123`,
			isJSON:      false, // We pass true to test the JSON path
			expectError: false,
		},
		{
			name:        "nil_content",
			content:     "",
			isJSON:      false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "nil_content" {
				result, diags := transforms.PayloadToGzip("")
				assert.Empty(t, result)
				assert.False(t, diags.HasError())

				return
			}

			// First encode content as base64 as the function expects
			encodedContent, err := transforms.Base64Encode(tt.content)
			require.NoError(t, err)

			// Add special handling for the "invalid_base64" test case
			if tt.name == "invalid_base64" {
				encodedContent = "this is not valid base64!"
			}

			result, diags := transforms.PayloadToGzip(encodedContent)

			if tt.expectError {
				assert.Empty(t, result)
				require.True(t, diags.HasError(), "Expected error diagnostics, got none")
			} else {
				require.False(t, diags.HasError(), "Unexpected error diagnostics: %v", diags)
				require.NotEmpty(t, result, "Expected non-empty result")

				// Try to decode to verify it's valid
				decodedContent, err := transforms.Base64GzipDecode(result)
				require.NoError(t, err)

				assert.Equal(t, tt.content, decodedContent)

				// For JSON content, verify it's valid JSON after decoding
				if tt.isJSON {
					assert.True(t, transforms.IsJSON(decodedContent))
				}
			}
		})
	}
}
