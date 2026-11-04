// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package deploymentplan

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType = fabcore.ItemTypeDeploymentPlan

	ItemDefinitionEmpty = `$schema: https://developer.microsoft.com/json-schemas/fabric/item/deploymentPlan/definition/plan/1.0.0/schema.json
version: 1.0.0
groups: []
`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/deployment-plan-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Deployment Plan",
	Type:           "deployment_plan",
	Names:          "Deployment Plans",
	Types:          "deployment_plans",
	DocsURL:        "https://learn.microsoft.com/fabric/cicd/deployment-plan/deployment-plan-actions",
	IsPreview:      true,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"plan.yml"},
	},
}
