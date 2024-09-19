// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sqlendpoint

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName           = "SQL Endpoint"
	ItemTFName         = "sql_endpoint"
	ItemsName          = "SQL Endpoints"
	ItemsTFName        = "sql_endpoints"
	ItemType           = fabcore.ItemTypeSQLEndpoint
	ItemDocsSPNSupport = common.DocsSPNNotSupported
	ItemDocsURL        = "https://learn.microsoft.com/power-bi/transform-model/sqlendpoints/sqlendpoints-overview"
)
