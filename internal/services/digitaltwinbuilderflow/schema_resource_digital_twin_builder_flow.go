// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package digitaltwinbuilderflow

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

func getResourceDigitalTwinBuilderFlowPropertiesAttributes() map[string]schema.Attribute {
	result := map[string]schema.Attribute{
		"item_id": schema.StringAttribute{
			MarkdownDescription: "The DigitalTwinBuilderFlow item ID.",
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"reference_type": schema.StringAttribute{
			MarkdownDescription: "The DigitalTwinBuilderFlow reference type.",
			Computed:            true,
		},
		"workspace_id": schema.StringAttribute{
			MarkdownDescription: "The workspace ID the DigitalTwinBuilderFlow belongs to.",
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
		},
	}

	return result
}

func getResourceDigitalTwinBuilderFlowConfigurationAttributes() map[string]schema.Attribute {
	result := map[string]schema.Attribute{
		"item_id": schema.StringAttribute{
			MarkdownDescription: "The DigitalTwinBuilderFlow item ID.",
			CustomType:          customtypes.UUIDType{},
			Computed:            true,
		},
		"reference_type": schema.StringAttribute{
			MarkdownDescription: "The DigitalTwinBuilderFlow reference type.",
			Computed:            true,
		},
		"workspace_id": schema.StringAttribute{
			MarkdownDescription: "The workspace ID the DigitalTwinBuilderFlow belongs to.",
			CustomType:          customtypes.UUIDType{},
			Computed:            true,
		},
	}

	return result
}
