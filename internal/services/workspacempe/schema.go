// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspacempe

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func itemSchema(dsList bool) superschema.Schema { //revive:disable-line:flag-parameter
	var dsTimeout *superschema.DatasourceTimeoutAttribute

	if !dsList {
		dsTimeout = &superschema.DatasourceTimeoutAttribute{
			Read: true,
		}
	}

	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false),
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, false),
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
					Optional: !dsList,
					Computed: true,
				},
			},
			"workspace_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Workspace ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: !dsList,
					Computed: dsList,
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " name.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.LengthAtMost(64),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Optional: !dsList,
					Computed: true,
				},
			},
			"provisioning_state": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Provisioning state of the endpoint.",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossiblePrivateEndpointProvisioningStateValues(), true)...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"target_private_link_resource_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Resource ID of data source for which private endpoint is created.",
					CustomType:          customtypes.CaseInsensitiveStringType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^/subscriptions/[a-f0-9-]+/resourceGroups/[a-zA-Z0-9-_]+/providers/[a-zA-Z0-9.]+/[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+$`),
							"Resource ID must be in the format `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}`.",
						),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"target_subresource_type": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Sub-resource pointing to [Private-link resource](https://learn.microsoft.com/azure/private-link/private-endpoint-overview#private-link-resource).",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"request_message": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "Request message.",
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.LengthAtMost(140),
					},
				},
			},
			"connection_state": superschema.SuperSingleNestedAttributeOf[connectionStateModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Endpoint connection state of provisioned endpoints.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Computed: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"actions_required": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Actions required to establish connection.",
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"status": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Connection status.",
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleConnectionStatusValues(), true)...),
							},
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"description": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Description message provided on approving or rejecting the end point.",
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
