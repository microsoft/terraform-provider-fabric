// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package digitaltwinbuilderflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	fabdigitaltwinbuilderflow "github.com/microsoft/fabric-sdk-go/fabric/digitaltwinbuilderflow"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getResourceDigitalTwinBuilderFlowProperties(ctx context.Context) map[string]schema.Attribute {
	result := map[string]schema.Attribute{
		"digital_twin_builder_item_reference": schema.SingleNestedAttribute{
			MarkdownDescription: "An object containing the properties of the Digital Twin Builder item reference.",
			Optional:            true,
			Computed:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[digitalTwinBuilderItemReferenceModel](ctx),
			Attributes: map[string]schema.Attribute{
				"item_id": schema.StringAttribute{
					MarkdownDescription: "The DigitalTwinBuilderFlow item ID.",
					Required:            true,
					CustomType:          customtypes.UUIDType{},
				},
				"reference_type": schema.StringAttribute{
					MarkdownDescription: "The DigitalTwinBuilderFlow reference type. Must be 'ById'.",
					Required:            true,
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabdigitaltwinbuilderflow.PossibleItemReferenceTypeValues(), true)...),
					},
				},
				"workspace_id": schema.StringAttribute{
					MarkdownDescription: "The workspace ID the DigitalTwinBuilderFlow belongs to.",
					Required:            true,
					CustomType:          customtypes.UUIDType{},
				},
			},
		},
	}

	return result
}

func getResourceDigitalTwinBuilderFlowConfigurationAttributes(ctx context.Context) map[string]schema.Attribute {
	result := map[string]schema.Attribute{
		"digital_twin_builder_item_reference": schema.SingleNestedAttribute{
			MarkdownDescription: "An object containing the properties of the Digital Twin Builder item reference.",
			Optional:            true,
			Computed:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[digitalTwinBuilderItemReferenceModel](ctx),
			Attributes: map[string]schema.Attribute{
				"item_id": schema.StringAttribute{
					MarkdownDescription: "The DigitalTwinBuilderFlow item ID.",
					Required:            true,
					CustomType:          customtypes.UUIDType{},
				},
				"reference_type": schema.StringAttribute{
					MarkdownDescription: "The DigitalTwinBuilderFlow reference type. Must be 'ById'.",
					Required:            true,
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabdigitaltwinbuilderflow.PossibleItemReferenceTypeValues(), true)...),
					},
				},
				"workspace_id": schema.StringAttribute{
					MarkdownDescription: "The workspace ID the DigitalTwinBuilderFlow belongs to.",
					Required:            true,
					CustomType:          customtypes.UUIDType{},
				},
			},
		},
	}

	return result
}
