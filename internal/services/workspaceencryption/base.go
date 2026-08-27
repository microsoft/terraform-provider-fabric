// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceencryption

import (
	"time"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Workspace Encryption",
	Type:           "workspace_encryption",
	DocsURL:        "https://learn.microsoft.com/fabric/security/workspace-customer-managed-keys",
	IsPreview:      false,
	IsSPNSupported: true,
}

// Applying or resetting a customer-managed key is asynchronous, so the status is polled at a fixed interval.
const encryptionPollInterval = 30 * time.Second
