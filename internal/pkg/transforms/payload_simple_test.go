// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package transforms_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSourceFileToPayload_TemplateNameFix(t *testing.T) {
	// Create a temporary template file with simple content
	tmpDir := t.TempDir()
	templateContent := `{
	"name": "{{ .notebook_name }}",
	"content": "Test notebook"
}`
	templateFile := filepath.Join(tmpDir, "test.json.tmpl")
	err := os.WriteFile(templateFile, []byte(templateContent), 0644)
	require.NoError(t, err)

	// Test single-threaded access
	tokens := map[string]string{
		"notebook_name": "my_notebook",
	}
	
	_, sha256Value, diags := transforms.SourceFileToPayload(templateFile, tokens)
	
	// Check for any errors
	if diags.HasError() {
		for _, diag := range diags.Errors() {
			t.Logf("Error: %s", diag.Summary())
			t.Logf("Detail: %s", diag.Detail())
		}
	}
	
	require.False(t, diags.HasError(), "Should not have errors in single-threaded access")
	assert.NotEmpty(t, sha256Value, "SHA256 value should not be empty")
}