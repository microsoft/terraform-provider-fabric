// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// WorkspacePrivateLinksPolicy is a policy that modifies URLs to use workspace-specific private link endpoints
// when workspace private links are enabled.
type WorkspacePrivateLinksPolicy struct {
	Enabled bool
}

// workspaceAPIRegex matches workspace-specific API endpoints that should be transformed for private links.
var workspaceAPIRegex = regexp.MustCompile(`^(/v[\d.]+/workspaces/([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})/.*)`)

func (p WorkspacePrivateLinksPolicy) Do(req *policy.Request) (*http.Response, error) {
	if !p.Enabled {
		return req.Next()
	}

	// Check if this is a workspace-specific API call that should use private links
	if matches := workspaceAPIRegex.FindStringSubmatch(req.Raw().URL.Path); len(matches) > 2 {
		workspaceID := matches[2]
		if err := p.transformURLForWorkspacePrivateLinks(req.Raw().URL, workspaceID); err != nil {
			// If transformation fails, continue with original URL
			// This ensures backwards compatibility
		}
	}

	return req.Next()
}

// transformURLForWorkspacePrivateLinks modifies the URL to use workspace-specific private link endpoints.
func (p WorkspacePrivateLinksPolicy) transformURLForWorkspacePrivateLinks(reqURL *url.URL, workspaceID string) error {
	// Transform the host to include workspace identifier for private links
	// Format: <workspace-id>.api.fabric.microsoft.com
	if strings.Contains(reqURL.Host, "api.fabric.microsoft.com") {
		reqURL.Host = fmt.Sprintf("%s.api.fabric.microsoft.com", workspaceID)

		return nil
	}

	return fmt.Errorf("unsupported host for workspace private links: %s", reqURL.Host)
}

var _ policy.Policy = WorkspacePrivateLinksPolicy{}

// WithWorkspacePrivateLinks returns a policy that enables workspace private links URL transformation.
func WithWorkspacePrivateLinks(enabled bool) policy.Policy {
	return WorkspacePrivateLinksPolicy{Enabled: enabled}
}
