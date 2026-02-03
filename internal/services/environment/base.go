// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package environment

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const FabricItemType = fabcore.ItemTypeEnvironment

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Environment",
	Type:           "environment",
	Names:          "Environments",
	Types:          "environments",
	DocsURL:        "https://learn.microsoft.com/fabric/data-engineering/create-and-use-environment",
	IsPreview:      true,
	IsSPNSupported: true,
}
