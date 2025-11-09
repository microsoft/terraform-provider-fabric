// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0
package client

import (
	"net/http"
	"net/url"
	"testing"
)

func TestWorkspacePrivateLinksPolicy_Do(t *testing.T) {
	tests := []struct {
		name               string
		enabled            bool
		originalURL        string
		expectedURL        string
		shouldTransformURL bool
	}{
		{
			name:               "workspace API with private links enabled",
			enabled:            true,
			originalURL:        "https://api.fabric.microsoft.com/v1/workspaces/12345678-1234-1234-1234-123456789012/items",
			expectedURL:        "https://12345678-1234-1234-1234-123456789012.api.fabric.microsoft.com/v1/workspaces/12345678-1234-1234-1234-123456789012/items",
			shouldTransformURL: true,
		},
		{
			name:               "workspace API with private links disabled",
			enabled:            false,
			originalURL:        "https://api.fabric.microsoft.com/v1/workspaces/12345678-1234-1234-1234-123456789012/items",
			expectedURL:        "https://api.fabric.microsoft.com/v1/workspaces/12345678-1234-1234-1234-123456789012/items",
			shouldTransformURL: false,
		},
		{
			name:               "non-workspace API with private links enabled",
			enabled:            true,
			originalURL:        "https://api.fabric.microsoft.com/v1/capacities",
			expectedURL:        "https://api.fabric.microsoft.com/v1/capacities",
			shouldTransformURL: false,
		},
		{
			name:               "different host with private links enabled",
			enabled:            true,
			originalURL:        "https://other.microsoft.com/v1/workspaces/12345678-1234-1234-1234-123456789012/items",
			expectedURL:        "https://other.microsoft.com/v1/workspaces/12345678-1234-1234-1234-123456789012/items",
			shouldTransformURL: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := WorkspacePrivateLinksPolicy{Enabled: tt.enabled}

			// Parse original URL
			originalURL, err := url.Parse(tt.originalURL)
			if err != nil {
				t.Fatalf("Failed to parse original URL: %v", err)
			}

			// Create a mock request
			req := &http.Request{
				URL: originalURL,
			}

			// Create a policy request
			policyReq := &policy.Request{}
			policyReq.SetRaw(req)

			// Mock the Next() call to return the transformed request
			policyReq.SetNext(func(*policy.Request) (*http.Response, error) {
				// Just return a dummy response, we're only interested in URL transformation
				return &http.Response{}, nil
			})

			// Execute the policy
			_, err = policy.Do(policyReq)
			if err != nil {
				t.Fatalf("Policy execution failed: %v", err)
			}

			// Check if URL was transformed correctly
			actualURL := req.URL.String()
			if actualURL != tt.expectedURL {
				t.Errorf("Expected URL %s, got %s", tt.expectedURL, actualURL)
			}
		})
	}
}

func TestWithWorkspacePrivateLinks(t *testing.T) {
	policy := WithWorkspacePrivateLinks(true)
	if policy == nil {
		t.Error("Expected non-nil policy")
	}

	// Verify it's the correct type
	if _, ok := policy.(WorkspacePrivateLinksPolicy); !ok {
		t.Error("Expected WorkspacePrivateLinksPolicy type")
	}
}

func TestWorkspaceAPIRegex(t *testing.T) {
	tests := []struct {
		path        string
		shouldMatch bool
		workspaceID string
	}{
		{
			path:        "/v1/workspaces/12345678-1234-1234-1234-123456789012/items",
			shouldMatch: true,
			workspaceID: "12345678-1234-1234-1234-123456789012",
		},
		{
			path:        "/v2.0/workspaces/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/folders",
			shouldMatch: true,
			workspaceID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		},
		{
			path:        "/v1/capacities",
			shouldMatch: false,
		},
		{
			path:        "/v1/workspaces",
			shouldMatch: false,
		},
		{
			path:        "/v1/workspaces/invalid-uuid/items",
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			matches := workspaceAPIRegex.FindStringSubmatch(tt.path)

			if tt.shouldMatch {
				if len(matches) < 3 {
					t.Errorf("Expected path %s to match regex, but it didn't", tt.path)
					return
				}

				actualWorkspaceID := matches[2]
				if actualWorkspaceID != tt.workspaceID {
					t.Errorf("Expected workspace ID %s, got %s", tt.workspaceID, actualWorkspaceID)
				}
			} else {
				if len(matches) > 0 {
					t.Errorf("Expected path %s to not match regex, but it did", tt.path)
				}
			}
		})
	}
}
