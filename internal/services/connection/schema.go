// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
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

	possibleConnectionEncryptionValues := []string{
		string(fabcore.ConnectionEncryptionEncrypted),
		string(fabcore.ConnectionEncryptionNotEncrypted),
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
					Optional: !isList,
					Computed: true,
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
					Optional: !isList,
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
					Validators: []validator.String{
						stringvalidator.OneOf(possibleSupportedConnectivityTypes...),
					},
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
					// why optional, why this default value?
					Default: stringdefault.StaticString(string(fabcore.PrivacyLevelOrganizational)),
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"gateway_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " gateway object ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
						stringplanmodifier.UseStateForUnknown(),
					},
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
				DataSource: &schemaD.StringAttribute{
					Computed: true,
					Optional: true,
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
									MarkdownDescription: "Name.",
									Required:            true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.RequiresReplace(),
									},
								},
							},
							"data_type": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "Data Type.",
									Computed:            true,
								},
							},
							"value": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "Value.",
									Required:            true,
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
								stringvalidator.OneOf(possibleConnectionEncryptionValues...),
							},
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
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
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
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
							Required: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.RequiresReplace(),
							},
						},
						DataSource: &schemaD.BoolAttribute{
							Computed: true,
						},
					},
					"credential_type": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The credential type.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								// note to self: which one will need authentication outside of the provider?
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleCredentialTypeValues(), true)...),
							},
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
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.RequiresReplace(),
									},
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
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.RequiresReplace(),
									},
								},
							},
							"client_id": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The client ID.",
									Required:            true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.RequiresReplace(),
									},
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
							MarkdownDescription: "The hared access signature credentials.",
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
							"sas_wo": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The sas (WO).",
									Required:            true,
									WriteOnly:           true,
								},
							},
							"sas_wo_version": superschema.Int32Attribute{
								Resource: &schemaR.Int32Attribute{
									MarkdownDescription: "The version of the `sas_wo`.",
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
