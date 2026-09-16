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
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
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
					MarkdownDescription: "The workspace ID.",
					CustomType:          customtypes.UUIDType{},
					Required:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"encryption_details": superschema.SuperSingleNestedAttributeOf[encryptionDetailsModel]{
				Resource: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The workspace encryption details.",
					Required:            true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					MarkdownDescription: "The workspace encryption details.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					"encryption_status": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The workspace encryption status.",
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleWorkspaceEncryptionStatusValues(), true)...),
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
				},
			},
			"previous_encryption_details": superschema.SuperSingleNestedAttributeOf[encryptionDetailsModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The previous workspace encryption details.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					"encryption_status": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The previous workspace encryption status.",
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleWorkspaceEncryptionStatusValues(), true)...),
							},
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
