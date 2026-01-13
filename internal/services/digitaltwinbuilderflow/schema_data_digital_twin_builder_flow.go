// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package digitaltwinbuilderflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

func getDataSourceDigitalTwinBuilderFlowPropertiesAttributes(ctx context.Context) map[string]schema.Attribute {
	result := map[string]schema.Attribute{
		"digital_twin_builder_item_reference": schema.SingleNestedAttribute{
			MarkdownDescription: "An object containing the properties of the SQL endpoint.",
			Computed:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[digitalTwinBuilderItemReferenceModel](ctx),
			Attributes: map[string]schema.Attribute{
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
			},
		},
	}

	return result
}
