// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package dashboard

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName           = "Dashboard"
	ItemTFName         = "dashboard"
	ItemsName          = "Dashboards"
	ItemsTFName        = "dashboards"
	ItemType           = fabcore.ItemTypeDashboard
	ItemDocsSPNSupport = common.DocsSPNNotSupported
	ItemDocsURL        = "https://learn.microsoft.com/power-bi/consumer/end-user-dashboards"
)
