// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakediagnosticsetting

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	superobjectvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/objectvalidator"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func itemSchema() superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false) +
				"\n\n**Important:** Destroying this resource disables OneLake diagnostics and clears the destination. Changing `workspace_id` also disables and clears diagnostics on the old workspace because the API does not expose a delete operation.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, false),
		},
		Attributes: map[string]superschema.Attribute{
			"workspace_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the source Workspace. The caller must have an Admin role on this Workspace.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
			},
			"status": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The status of OneLake diagnostic logging. Possible values: `Enabled`, `Disabled`.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf(diagnosticStatusEnabled, diagnosticStatusDisabled),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"destination": superschema.SuperSingleNestedAttributeOf[oneLakeDiagnosticDestinationModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The destination where OneLake diagnostic logs are stored. Required when diagnostics are enabled. When diagnostics are disabled, include the destination to retain it or omit it to clear it. The destination Lakehouse must be in the same capacity as the source Workspace, and the caller must have a Contributor or higher role on the destination Workspace.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
					Validators: []validator.Object{
						superobjectvalidator.RequireIfAttributeIsOneOf(
							path.MatchRoot("status"),
							[]attr.Value{types.StringValue(diagnosticStatusEnabled)},
						),
					},
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"lakehouse": superschema.SuperSingleNestedAttributeOf[oneLakeDiagnosticLakehouseModel]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "The Lakehouse receiving OneLake diagnostic logs. The API currently supports references by ID only.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Required: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: superschema.Attributes{
							"workspace_id": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The ID of the Workspace containing the destination Lakehouse.",
									CustomType:          customtypes.UUIDType{},
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"item_id": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The ID of the destination Lakehouse.",
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
				},
			},
			"timeouts": superschema.TimeoutAttribute{
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Read:   true,
					Update: true,
					Delete: true,
				},
				DataSource: &superschema.DatasourceTimeoutAttribute{
					Read: true,
				},
			},
		},
	}
}
