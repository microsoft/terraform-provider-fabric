// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	superboolvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/boolvalidator"
	superobjectvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/objectvalidator"
	superstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

//nolint:maintidx
func itemSchema(ctx context.Context, isList bool) superschema.Schema { //revive:disable-line:flag-parameter
	var dsTimeout *superschema.DatasourceTimeoutAttribute

	if !isList {
		dsTimeout = &superschema.DatasourceTimeoutAttribute{
			Read: true,
		}
	}

	possibleSupportedConnectivityTypes := []string{
		string(fabcore.ConnectivityTypeShareableCloud),
		string(fabcore.ConnectivityTypeVirtualNetworkGateway),
	}

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
					Required: !isList,
					Computed: isList,
				},
			},
			"display_name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " display name.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.LengthAtMost(123),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"connectivity_type": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " connectivity type.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.OneOf(possibleSupportedConnectivityTypes...),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"privacy_level": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " privacy level.",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossiblePrivacyLevelValues(), true)...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Computed: true,
					Default:  stringdefault.StaticString(string(fabcore.PrivacyLevelOrganizational)),
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"gateway_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " gateway object ID.",
					CustomType:          customtypes.UUIDType{},
					Validators: []validator.String{
						superstringvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("connectivity_type"),
							[]attr.Value{
								types.StringValue(string(fabcore.ConnectivityTypeVirtualNetworkGateway)),
							}),
						superstringvalidator.NullIfAttributeIsOneOf(path.MatchRoot("connectivity_type"),
							[]attr.Value{
								types.StringValue(string(fabcore.ConnectivityTypeShareableCloud)),
							}),
					},
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
					Optional: true,
				},
			},
			"allow_connection_usage_in_gateway": superschema.BoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Allow this connection to be utilized with either on-premises data gateways or VNet data gateways.",
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.Bool{
						superboolvalidator.NullIfAttributeIsOneOf(path.MatchRoot("connectivity_type"),
							[]attr.Value{
								types.StringValue(string(fabcore.ConnectivityTypeVirtualNetworkGateway)),
							}),
					},
				},
				DataSource: &schemaD.BoolAttribute{
					Computed: true,
					Optional: true,
				},
			},
			"connection_details": superschema.SuperSingleNestedAttribute{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " connection details.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Required:   true,
					CustomType: supertypes.NewSingleNestedObjectTypeOf[rsConnectionDetailsModel](ctx),
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed:   true,
					CustomType: supertypes.NewSingleNestedObjectTypeOf[dsConnectionDetailsModel](ctx),
				},
				Attributes: map[string]superschema.Attribute{
					"path": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The path of the connection.",
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"type": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The type of the connection.",
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
					"creation_method": superschema.StringAttribute{
						Resource: &schemaR.StringAttribute{
							MarkdownDescription: "The creation method used to create the connection.",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
					"parameters": superschema.SuperSetNestedAttributeOf[connectionParametersModel]{
						Resource: &schemaR.SetNestedAttribute{
							MarkdownDescription: "A set of connection parameters.",
							Optional:            true,
							PlanModifiers: []planmodifier.Set{
								setplanmodifier.RequiresReplace(),
							},
						},
						Attributes: superschema.Attributes{
							"name": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The name of the parameter.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(`\S`),
											"Name must contain at least one non-whitespace character.",
										),
									},
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.RequiresReplace(),
									},
								},
							},
							"data_type": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The data type of the parameter.",
									Computed:            true,
								},
							},
							"value": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The value of the parameter.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(`\S`),
											"Value must contain at least one non-whitespace character.",
										),
									},
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.RequiresReplace(),
									},
								},
							},
						},
					},
				},
			},
			"credential_details": superschema.SuperSingleNestedAttribute{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " credential details.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Required:   true,
					CustomType: supertypes.NewSingleNestedObjectTypeOf[rsCredentialDetailsModel](ctx),
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed:   true,
					CustomType: supertypes.NewSingleNestedObjectTypeOf[dsCredentialDetailsModel](ctx),
				},
				Attributes: map[string]superschema.Attribute{
					"connection_encryption": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The connection encryption type.",
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleConnectionEncryptionValues(), true)...),
							},
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  stringdefault.StaticString(string(fabcore.ConnectionEncryptionNotEncrypted)),
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"single_sign_on_type": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The single sign-on type.",
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleSingleSignOnTypeValues(), true)...),
							},
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  stringdefault.StaticString(string(fabcore.SingleSignOnTypeNone)),
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"skip_test_connection": superschema.BoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Whether the connection should skip the test connection during creation and update. `True` - Skip the test connection, `False` - Do not skip the test connection.",
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
						DataSource: &schemaD.BoolAttribute{
							Computed: true,
						},
					},
					"credential_type": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The credential type.",
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(
									utils.RemoveSliceByValue(fabcore.PossibleCredentialTypeValues(), fabcore.CredentialTypeOAuth2),
									true)...),
							},
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"basic_credentials": superschema.SuperSingleNestedAttributeOf[credentialsBasicModel]{
						Resource: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "The basic credentials.",
							Optional:            true,
							Validators: []validator.Object{
								superobjectvalidator.RequireIfAttributeIsOneOf(
									path.MatchRelative().AtParent().AtName("credential_type"),
									[]attr.Value{
										types.StringValue(string(fabcore.CredentialTypeBasic)),
									},
								),
							},
						},
						Attributes: superschema.Attributes{
							"username": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The username.",
									Required:            true,
								},
							},
							"password_wo": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The password (WO).",
									Required:            true,
									WriteOnly:           true,
								},
							},
							"password_wo_version": superschema.Int32Attribute{
								Resource: &schemaR.Int32Attribute{
									MarkdownDescription: "The version of the `password_wo`.",
									Required:            true,
								},
							},
						},
					},
					"key_credentials": superschema.SuperSingleNestedAttributeOf[credentialsKeyModel]{
						Resource: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "The key credentials.",
							Optional:            true,
							Validators: []validator.Object{
								superobjectvalidator.RequireIfAttributeIsOneOf(
									path.MatchRelative().AtParent().AtName("credential_type"),
									[]attr.Value{
										types.StringValue(string(fabcore.CredentialTypeKey)),
									},
								),
							},
						},
						Attributes: superschema.Attributes{
							"key_wo": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The key (WO).",
									Required:            true,
									WriteOnly:           true,
								},
							},
							"key_wo_version": superschema.Int32Attribute{
								Resource: &schemaR.Int32Attribute{
									MarkdownDescription: "The version of the `key_wo`.",
									Required:            true,
								},
							},
						},
					},
					"service_principal_credentials": superschema.SuperSingleNestedAttributeOf[credentialsServicePrincipalModel]{
						Resource: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "The service principal credentials.",
							Optional:            true,
							Validators: []validator.Object{
								superobjectvalidator.RequireIfAttributeIsOneOf(
									path.MatchRelative().AtParent().AtName("credential_type"),
									[]attr.Value{
										types.StringValue(string(fabcore.CredentialTypeServicePrincipal)),
									},
								),
							},
						},
						Attributes: superschema.Attributes{
							"tenant_id": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The tenant ID.",
									Required:            true,
								},
							},
							"client_id": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The client ID.",
									Required:            true,
								},
							},
							"client_secret_wo": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The client secret (WO).",
									Required:            true,
									WriteOnly:           true,
								},
							},
							"client_secret_wo_version": superschema.Int32Attribute{
								Resource: &schemaR.Int32Attribute{
									MarkdownDescription: "The version of the `client_secret_wo`.",
									Required:            true,
								},
							},
						},
					},
					"shared_access_signature_credentials": superschema.SuperSingleNestedAttributeOf[credentialsSharedAccessSignatureModel]{
						Resource: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "The shared access signature credentials.",
							Optional:            true,
							Validators: []validator.Object{
								superobjectvalidator.RequireIfAttributeIsOneOf(
									path.MatchRelative().AtParent().AtName("credential_type"),
									[]attr.Value{
										types.StringValue(string(fabcore.CredentialTypeSharedAccessSignature)),
									},
								),
							},
						},
						Attributes: superschema.Attributes{
							"token_wo": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The token (WO).",
									Required:            true,
									WriteOnly:           true,
								},
							},
							"token_wo_version": superschema.Int32Attribute{
								Resource: &schemaR.Int32Attribute{
									MarkdownDescription: "The version of the `token_wo`.",
									Required:            true,
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
				DataSource: dsTimeout,
			},
		},
	}
}
