// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Lakehouse",
	Type:           "lakehouse",
	Names:          "Lakehouses",
	Types:          "lakehouses",
	DocsURL:        "https://learn.microsoft.com/training/modules/get-started-lakehouses",
	IsPreview:      false,
	IsSPNSupported: true,
}

const FabricItemType = fabcore.ItemTypeLakehouse
