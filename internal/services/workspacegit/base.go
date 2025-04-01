// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspacegit

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Workspace Git",
	Type:           "workspace_git",
	DocsURL:        "https://learn.microsoft.com/fabric/cicd/git-integration/intro-to-git-integration",
	IsPreview:      true,
	IsSPNSupported: false,
}
