// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package datamart

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName           = "Datamart"
	ItemTFName         = "datamart"
	ItemsName          = "Datamarts"
	ItemsTFName        = "datamarts"
	ItemType           = fabcore.ItemTypeDatamart
	ItemDocsSPNSupport = common.DocsSPNNotSupported
	ItemDocsURL        = "https://learn.microsoft.com/power-bi/transform-model/datamarts/datamarts-overview"
)
