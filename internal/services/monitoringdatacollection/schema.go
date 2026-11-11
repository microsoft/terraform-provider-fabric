// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package monitoringdatacollection

import (
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func itemSchema() superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false) +
				"\n\nManages collection for an existing Monitoring item without managing its definition. " +
				"The caller needs the workspace Admin role and write permission on the item. " +
				"Configure routing in the item's definition before enabling collection. " +
				"Destroying this resource disables collection but does not delete the Monitoring item or change its definition. " +
				"Manage each Monitoring item's collection with only one instance of this resource.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Monitoring item ID.",
					CustomType:          customtypes.UUIDType{},
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"workspace_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Workspace ID.",
					CustomType:          customtypes.UUIDType{},
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"monitoring_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Monitoring item whose data collection is managed.",
					CustomType:          customtypes.UUIDType{},
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"enabled": superschema.BoolAttribute{
				Resource: &schemaR.BoolAttribute{
					MarkdownDescription: "Whether data collection is enabled. Disabling collection preserves the item's routing definition.",
					Required:            true,
				},
			},
			"timeouts": superschema.TimeoutAttribute{
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Read:   true,
					Update: true,
					Delete: true,
				},
			},
		},
	}
}
