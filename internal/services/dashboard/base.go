// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package dashboard

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Dashboard",
	Type:           "dashboard",
	Names:          "Dashboards",
	Types:          "dashboards",
	DocsURL:        "https://learn.microsoft.com/power-bi/consumer/end-user-dashboards",
	IsPreview:      true,
	IsSPNSupported: false,
}

const FabricItemType = fabcore.ItemTypeDashboard
