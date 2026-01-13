// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package digitaltwinbuilderflow

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeDigitalTwinBuilderFlow
	ItemFormatTypeDefault     = fabricitem.DefinitionFormatDefault
	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/digital-twin-builder-flow-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Digital Twin Builder Flow",
	Type:           "digital_twin_builder_flow",
	Names:          "Digital Twin Builder Flows",
	Types:          "digital_twin_builder_flows",
	DocsURL:        "https://learn.microsoft.com/fabric/real-time-intelligence/digital-twin-builder/overview",
	IsPreview:      true,
	IsSPNSupported: true,
}
