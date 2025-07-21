// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package transforms_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSourceFileToPayload_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	// Create a temporary template file using the same format as existing tests
	tmpDir := t.TempDir()
	templateContent := `{"name": "{{.NotebookName}}", "content": "This is notebook {{.NotebookName}}"}`
	templateFile := filepath.Join(tmpDir, "test.json.tmpl")
	err := os.WriteFile(templateFile, []byte(templateContent), 0644)
	require.NoError(t, err)

	// Test concurrent access to the same template file with different tokens
	const numGoroutines = 10
	var wg sync.WaitGroup
	results := make([]string, numGoroutines)
	hasErrors := make([]bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			
			tokens := map[string]string{
				"NotebookName": "notebook_" + string(rune('A'+index)),
			}
			
			_, sha256Value, diags := transforms.SourceFileToPayload(templateFile, tokens)
			if diags.HasError() {
				hasErrors[index] = true
				return
			}
			
			results[index] = sha256Value
		}(i)
	}

	wg.Wait()

	// Verify that no errors occurred
	for i, hasError := range hasErrors {
		assert.False(t, hasError, "goroutine %d should not have error", i)
	}

	// Verify that all results are valid (non-empty) and different
	// (because each uses different tokens)
	seenHashes := make(map[string]bool)
	for i, result := range results {
		if hasErrors[i] {
			continue // Skip results that had errors
		}
		assert.NotEmpty(t, result, "result %d should not be empty", i)
		assert.False(t, seenHashes[result], "result %d should be unique, got duplicate hash: %s", i, result)
		seenHashes[result] = true
	}

	// Verify that we got the expected number of unique results
	expectedResults := numGoroutines
	for _, hasError := range hasErrors {
		if hasError {
			expectedResults--
		}
	}
	assert.Equal(t, expectedResults, len(seenHashes), "should have %d unique hash results", expectedResults)
}

func TestSourceFileToPayload_SameTokensConcurrent(t *testing.T) {
	t.Parallel()

	// Create a temporary template file using the same format as existing tests
	tmpDir := t.TempDir()
	templateContent := `{"name": "{{.NotebookName}}", "id": "{{.NotebookID}}"}`
	templateFile := filepath.Join(tmpDir, "test.json.tmpl")
	err := os.WriteFile(templateFile, []byte(templateContent), 0644)
	require.NoError(t, err)

	// Test concurrent access with the same tokens - should produce identical results
	const numGoroutines = 5
	var wg sync.WaitGroup
	results := make([]string, numGoroutines)
	hasErrors := make([]bool, numGoroutines)
	
	tokens := map[string]string{
		"NotebookName": "shared_notebook",
		"NotebookID":   "42",
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			
			_, sha256Value, diags := transforms.SourceFileToPayload(templateFile, tokens)
			if diags.HasError() {
				hasErrors[index] = true
				return
			}
			
			results[index] = sha256Value
		}(i)
	}

	wg.Wait()

	// Verify that no errors occurred
	for i, hasError := range hasErrors {
		assert.False(t, hasError, "goroutine %d should not have error", i)
	}

	// Find the first successful result to use as expected value
	var expectedHash string
	for i, result := range results {
		if !hasErrors[i] && result != "" {
			expectedHash = result
			break
		}
	}
	assert.NotEmpty(t, expectedHash, "should have at least one successful result")

	// Verify that all successful results are identical (same tokens should produce same hash)
	for i, result := range results {
		if !hasErrors[i] {
			assert.Equal(t, expectedHash, result, "result %d should match the expected hash", i)
		}
	}
}