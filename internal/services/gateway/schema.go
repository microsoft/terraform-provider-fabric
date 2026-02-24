// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	superint32validator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/int32validator"
	superobjectvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/objectvalidator"
	superstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

//nolint:maintidx
func itemSchema(isList bool) superschema.Schema { //revive:disable-line:flag-parameter
	var dsTimeout *superschema.DatasourceTimeoutAttribute

	if !isList {
		dsTimeout = &superschema.DatasourceTimeoutAttribute{
			Read: true,
		}
	}

	possibleGatewayTypeValues := utils.RemoveSlicesByValues(fabcore.PossibleGatewayTypeValues(), []fabcore.GatewayType{fabcore.GatewayTypeOnPremises, fabcore.GatewayTypeOnPremisesPersonal})

	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false),
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, isList),
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Optional: !isList,
					Computed: true,
				},
			},
			"display_name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " display name.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.LengthAtMost(200),
						superstringvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.GatewayTypeOnPremises)),
								types.StringValue(string(fabcore.GatewayTypeVirtualNetwork)),
							}),
						superstringvalidator.NullIfAttributeIsOneOf(path.MatchRoot("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.GatewayTypeOnPremisesPersonal)),
							}),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Optional: !isList,
					Computed: true,
				},
			},
			"type": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " type.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(possibleGatewayTypeValues, true)...),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleGatewayTypeValues(), true)...),
					},
				},
			},
			"capacity_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The capacity ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						superstringvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.GatewayTypeVirtualNetwork)),
							}),
						superstringvalidator.NullIfAttributeIsOneOf(path.MatchRoot("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.GatewayTypeOnPremises)),
								types.StringValue(string(fabcore.GatewayTypeOnPremisesPersonal)),
							}),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"inactivity_minutes_before_sleep": superschema.Int32Attribute{
				Common: &schemaR.Int32Attribute{
					MarkdownDescription: "The inactivity minutes before sleep.",
					Validators: []validator.Int32{
						int32validator.OneOf(PossibleInactivityMinutesBeforeSleepValues...),
					},
				},
				Resource: &schemaR.Int32Attribute{
					Optional: true,
					Computed: true,
					Validators: []validator.Int32{
						superint32validator.RequireIfAttributeIsOneOf(path.MatchRoot("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.GatewayTypeVirtualNetwork)),
							}),
						superint32validator.NullIfAttributeIsOneOf(path.MatchRoot("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.GatewayTypeOnPremises)),
								types.StringValue(string(fabcore.GatewayTypeOnPremisesPersonal)),
							}),
					},
				},
				DataSource: &schemaD.Int32Attribute{
					Computed: true,
				},
			},
			"number_of_member_gateways": superschema.Int32Attribute{
				Common: &schemaR.Int32Attribute{
					MarkdownDescription: "The number of member gateways.",
					Validators: []validator.Int32{
						int32validator.Between(MinNumberOfMemberGatewaysValues, MaxNumberOfMemberGatewaysValues),
					},
				},
				Resource: &schemaR.Int32Attribute{
					Optional: true,
					Computed: true,
					Validators: []validator.Int32{
						superint32validator.RequireIfAttributeIsOneOf(path.MatchRoot("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.GatewayTypeVirtualNetwork)),
							}),
						superint32validator.NullIfAttributeIsOneOf(path.MatchRoot("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.GatewayTypeOnPremises)),
								types.StringValue(string(fabcore.GatewayTypeOnPremisesPersonal)),
							}),
					},
				},
				DataSource: &schemaD.Int32Attribute{
					Computed: true,
				},
			},
			"virtual_network_azure_resource": superschema.SuperSingleNestedAttributeOf[virtualNetworkAzureResourceModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The Azure virtual network resource.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.RequiresReplace(),
					},
					Validators: []validator.Object{
						superobjectvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.GatewayTypeVirtualNetwork)),
							}),
						superobjectvalidator.NullIfAttributeIsOneOf(path.MatchRoot("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.GatewayTypeOnPremises)),
								types.StringValue(string(fabcore.GatewayTypeOnPremisesPersonal)),
							}),
					},
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"virtual_network_name": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The virtual network name.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"subnet_name": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The subnet name.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"resource_group_name": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The resource group name.",
							CustomType:          customtypes.CaseInsensitiveStringType{},
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"subscription_id": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The subscription ID.",
							CustomType:          customtypes.UUIDType{},
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"allow_cloud_connection_refresh": superschema.BoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Allow cloud connection refresh.",
				},
				Resource: &schemaR.BoolAttribute{
					Computed: true,
				},
				DataSource: &schemaD.BoolAttribute{
					Computed: true,
				},
			},
			"allow_custom_connectors": superschema.BoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Allow custom connectors.",
				},
				Resource: &schemaR.BoolAttribute{
					Computed: true,
				},
				DataSource: &schemaD.BoolAttribute{
					Computed: true,
				},
			},
			"load_balancing_setting": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The load balancing setting",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleLoadBalancingSettingValues(), true)...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"public_key": superschema.SuperSingleNestedAttributeOf[publicKeyModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The public key of the primary gateway member. Used to encrypt the credentials for creating and updating connections.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Computed: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"exponent": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The exponent.",
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"modulus": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The modulus.",
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"version": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " version.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
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
		},
	}
}
