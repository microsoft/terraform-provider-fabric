// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package folder

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Folder",
	Type:           "folder",
	Names:          "Folders",
	Types:          "folders",
	DocsURL:        "https://learn.microsoft.com/fabric/fundamentals/workspaces-folders",
	IsPreview:      true,
	IsSPNSupported: true,
}
