// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package semanticmodelconnectionbinding

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabsemanticmodel "github.com/microsoft/fabric-sdk-go/fabric/semanticmodel"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	superstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func itemSchema(_ context.Context, isList bool) superschema.Schema { //revive:disable-line:flag-parameter
	dsTimeout := &superschema.DatasourceTimeoutAttribute{Read: true}
	var (
		idDataSource                *schemaD.StringAttribute
		connectivityDataSource      *schemaD.StringAttribute
		connectionDetailsDataSource *schemaD.SingleNestedAttribute
		connectionPathDataSource    *schemaD.StringAttribute
		connectionTypeDataSource    *schemaD.StringAttribute
		displayNameDataSource       *schemaD.StringAttribute
		gatewayIDDataSource         *schemaD.StringAttribute
	)

	if !isList {
		idDataSource = &schemaD.StringAttribute{Computed: true}
		connectivityDataSource = &schemaD.StringAttribute{
			MarkdownDescription: "The current connectivity type. `None` means the data-source reference is currently unbound.",
			Computed:            true,
		}
		connectionDetailsDataSource = &schemaD.SingleNestedAttribute{
			Required: true,
		}
		connectionPathDataSource = &schemaD.StringAttribute{
			Required: true,
		}
		connectionTypeDataSource = &schemaD.StringAttribute{
			Required: true,
		}
		displayNameDataSource = &schemaD.StringAttribute{Computed: true}
		gatewayIDDataSource = &schemaD.StringAttribute{Computed: true}
	}

	possibleConnectivityTypeValues := utils.RemoveSlicesByValues(
		fabsemanticmodel.PossibleConnectivityTypeValues(),
		[]fabsemanticmodel.ConnectivityType{fabsemanticmodel.ConnectivityTypeNone},
	)
	connectionIDRequiredConnectivityTypeValues := utils.RemoveSlicesByValues(
		possibleConnectivityTypeValues,
		[]fabsemanticmodel.ConnectivityType{fabsemanticmodel.ConnectivityTypeAutomatic},
	)

	requiredConnectionTypes := make([]attr.Value, len(connectionIDRequiredConnectivityTypeValues))
	for i, connectivityType := range connectionIDRequiredConnectivityTypeValues {
		requiredConnectionTypes[i] = types.StringValue(string(connectivityType))
	}

	attributes := superschema.Attributes{
		"workspace_id": superschema.SuperStringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The workspace ID that contains the semantic model.",
				CustomType:          customtypes.UUIDType{},
			},
			Resource: &schemaR.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			DataSource: &schemaD.StringAttribute{Required: true},
		},
		"semantic_model_id": superschema.SuperStringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The semantic model ID.",
				CustomType:          customtypes.UUIDType{},
			},
			Resource: &schemaR.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			DataSource: &schemaD.StringAttribute{Required: true},
		},
		"id": superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "A stable identifier for this binding, derived from its semantic model and connection details.",
			},
			Resource: &schemaR.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			DataSource: idDataSource,
		},
		"connection_id": superschema.SuperStringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The ID of the bound Fabric connection.",
				CustomType:          customtypes.UUIDType{},
			},
			Resource: &schemaR.StringAttribute{
				MarkdownDescription: "Optional when `connectivity_type` is `Automatic`.",
				Optional:            true,
				Validators: []validator.String{
					superstringvalidator.RequireIfAttributeIsOneOf(
						path.MatchRoot("connectivity_type"),
						requiredConnectionTypes,
					),
					superstringvalidator.NullIfAttributeIsOneOf(
						path.MatchRoot("connectivity_type"),
						[]attr.Value{types.StringValue(string(fabsemanticmodel.ConnectivityTypeAutomatic))},
					),
				},
			},
			DataSource: &schemaD.StringAttribute{
				MarkdownDescription: "Null when `connectivity_type` is `Automatic` or `None` (unbound).",
				Computed:            true,
			},
		},
		"connectivity_type": superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The Fabric connectivity type of the binding.",
			},
			Resource: &schemaR.StringAttribute{
				MarkdownDescription: "Select the connectivity type to use for this data-source reference. `None` is reserved for unbinding when the resource is deleted.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(possibleConnectivityTypeValues, true)...),
				},
			},
			DataSource: connectivityDataSource,
		},
		"connection_details": superschema.SuperSingleNestedAttributeOf[connectionDetailsModel]{
			Common: &schemaR.SingleNestedAttribute{
				MarkdownDescription: "The semantic model data-source reference for the binding. Its path and type must exactly match a reference returned by the List Item Connections API.",
			},
			Resource:   &schemaR.SingleNestedAttribute{Required: true},
			DataSource: connectionDetailsDataSource,
			Attributes: superschema.Attributes{
				"path": superschema.StringAttribute{
					Common: &schemaR.StringAttribute{
						MarkdownDescription: "The path of the semantic model data-source reference.",
					},
					Resource: &schemaR.StringAttribute{
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					DataSource: connectionPathDataSource,
				},
				"type": superschema.StringAttribute{
					Common: &schemaR.StringAttribute{
						MarkdownDescription: "The type of the semantic model data-source reference.",
					},
					Resource: &schemaR.StringAttribute{
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					DataSource: connectionTypeDataSource,
				},
			},
		},
		"display_name": superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The display name of the bound Fabric connection.",
			},
			Resource:   &schemaR.StringAttribute{Computed: true},
			DataSource: displayNameDataSource,
		},
		"gateway_id": superschema.SuperStringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The ID of the gateway used by the bound connection.",
				CustomType:          customtypes.UUIDType{},
			},
			Resource:   &schemaR.StringAttribute{Computed: true},
			DataSource: gatewayIDDataSource,
		},
		"timeouts": superschema.TimeoutAttribute{
			Resource: &superschema.ResourceTimeoutAttribute{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
			DataSource: dsTimeout,
		},
	}

	if isList {
		delete(attributes, "id")
		delete(attributes, "connection_id")
		delete(attributes, "connectivity_type")
		delete(attributes, "connection_details")
		delete(attributes, "display_name")
		delete(attributes, "gateway_id")

		attributes["values"] = superschema.SuperSetNestedAttributeOf[connectionBindingPropertiesModel]{
			DataSource: &schemaD.SetNestedAttribute{
				MarkdownDescription: "The set of connection bindings for the semantic model.",
				Computed:            true,
			},
			Attributes: superschema.Attributes{
				"id": superschema.StringAttribute{
					DataSource: &schemaD.StringAttribute{
						MarkdownDescription: "A stable identifier for this binding.",
						Computed:            true,
					},
				},
				"connection_id": superschema.SuperStringAttribute{
					Common: &schemaR.StringAttribute{
						MarkdownDescription: "The ID of the Fabric connection.",
						CustomType:          customtypes.UUIDType{},
					},
					DataSource: &schemaD.StringAttribute{Computed: true},
				},
				"connectivity_type": superschema.StringAttribute{
					DataSource: &schemaD.StringAttribute{
						MarkdownDescription: "The current connectivity type. `None` means the data-source reference is currently unbound.",
						Computed:            true,
					},
				},
				"connection_details": superschema.SuperSingleNestedAttributeOf[connectionDetailsModel]{
					DataSource: &schemaD.SingleNestedAttribute{
						MarkdownDescription: "The semantic model data-source reference for this binding.",
						Computed:            true,
					},
					Attributes: superschema.Attributes{
						"path": superschema.StringAttribute{
							DataSource: &schemaD.StringAttribute{
								MarkdownDescription: "The path of the semantic model data-source reference.",
								Computed:            true,
							},
						},
						"type": superschema.StringAttribute{
							DataSource: &schemaD.StringAttribute{
								MarkdownDescription: "The type of the semantic model data-source reference.",
								Computed:            true,
							},
						},
					},
				},
				"display_name": superschema.StringAttribute{
					DataSource: &schemaD.StringAttribute{
						MarkdownDescription: "The display name of the bound Fabric connection.",
						Computed:            true,
					},
				},
				"gateway_id": superschema.SuperStringAttribute{
					Common: &schemaR.StringAttribute{
						MarkdownDescription: "The ID of the gateway used by the bound connection.",
						CustomType:          customtypes.UUIDType{},
					},
					DataSource: &schemaD.StringAttribute{Computed: true},
				},
			},
		}
	}

	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false),
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, isList),
		},
		Attributes: attributes,
	}
}

func dataSourceSchema(ctx context.Context, isList bool) schemaD.Schema {
	return itemSchema(ctx, isList).GetDataSource(ctx)
}
