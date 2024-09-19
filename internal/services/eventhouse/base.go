// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventhouse

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName           = "Eventhouse"
	ItemTFName         = "eventhouse"
	ItemsName          = "Eventhouses"
	ItemsTFName        = "eventhouses"
	ItemType           = fabcore.ItemTypeEventhouse
	ItemDocsSPNSupport = common.DocsSPNSupported
	ItemDocsURL        = "https://learn.microsoft.com/fabric/real-time-intelligence/eventhouse"
)
