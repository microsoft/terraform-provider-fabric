// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package report

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName                  = "Report"
	ItemTFName                = "report"
	ItemsName                 = "Reports"
	ItemsTFName               = "reports"
	ItemType                  = fabcore.ItemTypeReport
	ItemDocsSPNSupport        = common.DocsSPNSupported
	ItemDocsURL               = "https://learn.microsoft.com/power-bi/developer/projects/projects-report"
	ItemFormatTypeDefault     = "PBIR-Legacy"
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/report-definition"
)

var (
	ItemFormatTypes               = []string{"PBIR-Legacy"}                                                                                                  //nolint:gochecknoglobals
	ItemDefinitionPathsPBIRLegacy = []string{"report.json", "definition.pbir", "StaticResources/RegisteredResources/*", "StaticResources/SharedResources/*"} //nolint:gochecknoglobals
)
