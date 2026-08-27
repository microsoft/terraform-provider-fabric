// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceencryption

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator" //revive:disable-line:import-alias-naming
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"  //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"    //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

// Fabric requires a versionless key identifier. The host is intentionally unconstrained, because the vault DNS
// suffix varies by Azure environment and the API validates the vault itself.
var keyIdentifierRegex = regexp.MustCompile(`^https://[^/]+/keys/[^/]+/?$`)

func itemSchema() superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false),
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, false),
		},
		Attributes: map[string]superschema.Attribute{
			"workspace_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Workspace ID.",
					CustomType:          customtypes.UUIDType{},
					Required:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"key_identifier": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Azure Key Vault key identifier. Changing this value rotates the customer-managed key.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(keyIdentifierRegex, "must be a versionless key identifier, for example: https://myvault.vault.azure.net/keys/mykey/"),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"encryption_status": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Workspace encryption status.",
					Computed:            true,
				},
			},
			"previous_encryption_detail": superschema.SuperSingleNestedAttributeOf[encryptionDetailModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The previous workspace encryption detail.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					"encryption_status": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The previous workspace encryption status.",
							Computed:            true,
						},
					},
					"key_identifier": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The previous key identifier.",
							Computed:            true,
						},
					},
				},
			},
			"workspace_encryption_items_details": superschema.SuperSetNestedAttributeOf[workspaceEncryptionItemsDetailModel]{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "The encryption status of items in the workspace.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					"encryption_status": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The encryption status for the items.",
							Computed:            true,
						},
					},
					"items": superschema.SuperSetNestedAttributeOf[workspaceEncryptionItemModel]{
						Common: &schemaR.SetNestedAttribute{
							MarkdownDescription: "The array of workspace item details.",
							Computed:            true,
						},
						Attributes: superschema.Attributes{
							"id": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The item ID.",
									CustomType:          customtypes.UUIDType{},
									Computed:            true,
								},
							},
							"display_name": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The item display name.",
									Computed:            true,
								},
							},
							"type": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The item type.",
									Computed:            true,
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
