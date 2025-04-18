// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
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

func itemSchema(ctx context.Context, isList bool) superschema.Schema { //revive:disable-line:flag-parameter
	var dsTimeout *superschema.DatasourceTimeoutAttribute

	if !isList {
		dsTimeout = &superschema.DatasourceTimeoutAttribute{
			Read: true,
		}
	}

	possibleConnectivityTypeValues := utils.RemoveSlicesByValues(
		fabcore.PossibleConnectivityTypeValues(),
		[]fabcore.ConnectivityType{fabcore.ConnectivityTypeOnPremisesGateway, fabcore.ConnectivityTypeOnPremisesGatewayPersonal, fabcore.ConnectivityTypePersonalCloud},
	)

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
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(possibleConnectivityTypeValues, false)...),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleConnectivityTypeValues(), false)...),
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
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						superstringvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("connectivity_type"),
							[]attr.Value{
								types.StringValue(string(fabcore.ConnectivityTypeVirtualNetworkGateway)),
								types.StringValue(string(fabcore.ConnectivityTypeOnPremisesGateway)),
								types.StringValue(string(fabcore.ConnectivityTypeOnPremisesGatewayPersonal)),
							}),
						superstringvalidator.NullIfAttributeIsOneOf(path.MatchRoot("connectivity_type"),
							[]attr.Value{
								types.StringValue(string(fabcore.ConnectivityTypeAutomatic)),
								types.StringValue(string(fabcore.ConnectivityTypeNone)),
								types.StringValue(string(fabcore.ConnectivityTypePersonalCloud)),
								types.StringValue(string(fabcore.ConnectivityTypeShareableCloud)),
							}),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
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
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleConnectionEncryptionValues(), true)...),
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
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleCredentialTypeValues(), true)...),
							},
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
					"basic_credentials": superschema.SuperSingleNestedAttributeOf[credentialsBasicModel]{
						Resource: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "The basic credentials.",
							Optional:            true,
							Validators: []validator.Object{
								objectvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("key_credentials"),
									path.MatchRelative().AtParent().AtName("service_principal_credentials"),
									path.MatchRelative().AtParent().AtName("shared_access_signature_credentials"),
									path.MatchRelative().AtParent().AtName("windows_credentials"),
									// path.MatchRelative().AtParent().AtName("encrypted_credentials"),
								),
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
							"password": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The password.",
									DeprecationMessage:  "This attribute is deprecated. Use `password_wo` instead.",
									Optional:            true,
									Sensitive:           true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.RequiresReplace(),
									},
									Validators: []validator.String{
										stringvalidator.PreferWriteOnlyAttribute(
											path.MatchRelative().AtParent().AtName("password_wo"),
										),
										stringvalidator.ConflictsWith(
											path.MatchRelative().AtParent().AtName("password_wo"),
										),
										stringvalidator.ExactlyOneOf(
											path.MatchRelative().AtParent().AtName("password"),
											path.MatchRelative().AtParent().AtName("password_wo"),
										),
									},
								},
							},
							"password_wo": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The password (WO).",
									Optional:            true,
									WriteOnly:           true,
									Validators: []validator.String{
										stringvalidator.ConflictsWith(
											path.MatchRelative().AtParent().AtName("password"),
										),
										stringvalidator.AlsoRequires(
											path.MatchRelative().AtParent().AtName("password_wo_version"),
										),
										stringvalidator.ExactlyOneOf(
											path.MatchRelative().AtParent().AtName("password"),
											path.MatchRelative().AtParent().AtName("password_wo"),
										),
									},
								},
							},
							"password_wo_version": superschema.Int32Attribute{
								Resource: &schemaR.Int32Attribute{
									MarkdownDescription: "The version of the password_wo.",
									Optional:            true,
									PlanModifiers: []planmodifier.Int32{
										int32planmodifier.RequiresReplace(),
									},
									Validators: []validator.Int32{
										int32validator.ConflictsWith(
											path.MatchRelative().AtParent().AtName("password"),
										),
										int32validator.AlsoRequires(
											path.MatchRelative().AtParent().AtName("password_wo"),
										),
									},
								},
							},
						},
					},
					"key_credentials": superschema.SuperSingleNestedAttributeOf[credentialsKeyModel]{
						Resource: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "The key credentials.",
							Optional:            true,
							Validators: []validator.Object{
								objectvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("basic_credentials"),
									path.MatchRelative().AtParent().AtName("service_principal_credentials"),
									path.MatchRelative().AtParent().AtName("shared_access_signature_credentials"),
									path.MatchRelative().AtParent().AtName("windows_credentials"),
									// path.MatchRelative().AtParent().AtName("encrypted_credentials"),
								),
								superobjectvalidator.RequireIfAttributeIsOneOf(
									path.MatchRelative().AtParent().AtName("credential_type"),
									[]attr.Value{
										types.StringValue(string(fabcore.CredentialTypeKey)),
									},
								),
							},
						},
						Attributes: superschema.Attributes{
							"key": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The key.",
									DeprecationMessage:  "This attribute is deprecated. Use `key_wo` instead.",
									Optional:            true,
									Sensitive:           true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.RequiresReplace(),
									},
									Validators: []validator.String{
										stringvalidator.PreferWriteOnlyAttribute(
											path.MatchRelative().AtParent().AtName("key_wo"),
										),
										stringvalidator.ConflictsWith(
											path.MatchRelative().AtParent().AtName("key_wo"),
										),
										stringvalidator.ExactlyOneOf(
											path.MatchRelative().AtParent().AtName("key"),
											path.MatchRelative().AtParent().AtName("key_wo"),
										),
									},
								},
							},
							"key_wo": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The key (WO).",
									Optional:            true,
									WriteOnly:           true,
									Validators: []validator.String{
										stringvalidator.ConflictsWith(
											path.MatchRelative().AtParent().AtName("key"),
										),
										stringvalidator.AlsoRequires(
											path.MatchRelative().AtParent().AtName("key_wo_version"),
										),
										stringvalidator.ExactlyOneOf(
											path.MatchRelative().AtParent().AtName("key"),
											path.MatchRelative().AtParent().AtName("key_wo"),
										),
									},
								},
							},
							"key_wo_version": superschema.Int32Attribute{
								Resource: &schemaR.Int32Attribute{
									MarkdownDescription: "The version of the key_wo.",
									Optional:            true,
									PlanModifiers: []planmodifier.Int32{
										int32planmodifier.RequiresReplace(),
									},
									Validators: []validator.Int32{
										int32validator.ConflictsWith(
											path.MatchRelative().AtParent().AtName("key"),
										),
										int32validator.AlsoRequires(
											path.MatchRelative().AtParent().AtName("key_wo"),
										),
									},
								},
							},
						},
					},
					"service_principal_credentials": superschema.SuperSingleNestedAttributeOf[credentialsServicePrincipalModel]{
						Resource: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "The service principal credentials.",
							Optional:            true,
							Validators: []validator.Object{
								objectvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("basic_credentials"),
									path.MatchRelative().AtParent().AtName("key_credentials"),
									path.MatchRelative().AtParent().AtName("shared_access_signature_credentials"),
									path.MatchRelative().AtParent().AtName("windows_credentials"),
									// path.MatchRelative().AtParent().AtName("encrypted_credentials"),
								),
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

							"client_secret": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The client secret.",
									DeprecationMessage:  "This attribute is deprecated. Use `client_secret_wo` instead.",
									Optional:            true,
									Sensitive:           true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.RequiresReplace(),
									},
									Validators: []validator.String{
										stringvalidator.PreferWriteOnlyAttribute(
											path.MatchRelative().AtParent().AtName("client_secret_wo"),
										),
										stringvalidator.ConflictsWith(
											path.MatchRelative().AtParent().AtName("client_secret_wo"),
										),
										stringvalidator.ExactlyOneOf(
											path.MatchRelative().AtParent().AtName("client_secret"),
											path.MatchRelative().AtParent().AtName("client_secret_wo"),
										),
									},
								},
							},
							"client_secret_wo": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The client secret (WO).",
									Optional:            true,
									WriteOnly:           true,
									Validators: []validator.String{
										stringvalidator.ConflictsWith(
											path.MatchRelative().AtParent().AtName("client_secret"),
										),
										stringvalidator.AlsoRequires(
											path.MatchRelative().AtParent().AtName("client_secret_version"),
										),
										stringvalidator.ExactlyOneOf(
											path.MatchRelative().AtParent().AtName("client_secret"),
											path.MatchRelative().AtParent().AtName("client_secret_wo"),
										),
									},
								},
							},
							"client_secret_wo_version": superschema.Int32Attribute{
								Resource: &schemaR.Int32Attribute{
									MarkdownDescription: "The version of the client_secret_wo.",
									Optional:            true,
									PlanModifiers: []planmodifier.Int32{
										int32planmodifier.RequiresReplace(),
									},
									Validators: []validator.Int32{
										int32validator.ConflictsWith(
											path.MatchRelative().AtParent().AtName("client_secret"),
										),
										int32validator.AlsoRequires(
											path.MatchRelative().AtParent().AtName("client_secret_wo"),
										),
									},
								},
							},
						},
					},
					"shared_access_signature_credentials": superschema.SuperSingleNestedAttributeOf[credentialsSharedAccessSignatureModel]{
						Resource: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "The hared access signature credentials.",
							Optional:            true,
							Validators: []validator.Object{
								objectvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("basic_credentials"),
									path.MatchRelative().AtParent().AtName("key_credentials"),
									path.MatchRelative().AtParent().AtName("service_principal_credentials"),
									path.MatchRelative().AtParent().AtName("windows_credentials"),
									// path.MatchRelative().AtParent().AtName("encrypted_credentials"),
								),
								superobjectvalidator.RequireIfAttributeIsOneOf(
									path.MatchRelative().AtParent().AtName("credential_type"),
									[]attr.Value{
										types.StringValue(string(fabcore.CredentialTypeSharedAccessSignature)),
									},
								),
							},
						},
						Attributes: superschema.Attributes{
							"sas": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The sas.",
									DeprecationMessage:  "This attribute is deprecated. Use `sas_wo` instead.",
									Optional:            true,
									Sensitive:           true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.RequiresReplace(),
									},
									Validators: []validator.String{
										stringvalidator.PreferWriteOnlyAttribute(
											path.MatchRelative().AtParent().AtName("sas_wo"),
										),
										stringvalidator.ConflictsWith(
											path.MatchRelative().AtParent().AtName("sas_wo"),
										),
										stringvalidator.ExactlyOneOf(
											path.MatchRelative().AtParent().AtName("sas"),
											path.MatchRelative().AtParent().AtName("sas_wo"),
										),
									},
								},
							},
							"sas_wo": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The sas (WO).",
									Optional:            true,
									WriteOnly:           true,
									Validators: []validator.String{
										stringvalidator.ConflictsWith(
											path.MatchRelative().AtParent().AtName("sas"),
										),
										stringvalidator.AlsoRequires(
											path.MatchRelative().AtParent().AtName("sas_wo_version"),
										),
										stringvalidator.ExactlyOneOf(
											path.MatchRelative().AtParent().AtName("sas"),
											path.MatchRelative().AtParent().AtName("sas_wo"),
										),
									},
								},
							},
							"sas_wo_version": superschema.Int32Attribute{
								Resource: &schemaR.Int32Attribute{
									MarkdownDescription: "The version of the sas_wo.",
									Optional:            true,
									PlanModifiers: []planmodifier.Int32{
										int32planmodifier.RequiresReplace(),
									},
									Validators: []validator.Int32{
										int32validator.ConflictsWith(
											path.MatchRelative().AtParent().AtName("sas"),
										),
										int32validator.AlsoRequires(
											path.MatchRelative().AtParent().AtName("sas_wo"),
										),
									},
								},
							},
						},
					},
					"windows_credentials": superschema.SuperSingleNestedAttributeOf[credentialsWindowsModel]{
						Resource: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "The Windows credentials.",
							Optional:            true,
							Validators: []validator.Object{
								objectvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("basic_credentials"),
									path.MatchRelative().AtParent().AtName("key_credentials"),
									path.MatchRelative().AtParent().AtName("service_principal_credentials"),
									path.MatchRelative().AtParent().AtName("shared_access_signature_credentials"),
									// path.MatchRelative().AtParent().AtName("encrypted_credentials"),
								),
								superobjectvalidator.RequireIfAttributeIsOneOf(
									path.MatchRelative().AtParent().AtName("credential_type"),
									[]attr.Value{
										types.StringValue(string(fabcore.CredentialTypeWindows)),
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
							"password": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The password.",
									DeprecationMessage:  "This attribute is deprecated. Use `password_wo` instead.",
									Optional:            true,
									Sensitive:           true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.RequiresReplace(),
									},
									Validators: []validator.String{
										stringvalidator.PreferWriteOnlyAttribute(
											path.MatchRelative().AtParent().AtName("password_wo"),
										),
										stringvalidator.ConflictsWith(
											path.MatchRelative().AtParent().AtName("password_wo"),
										),
										stringvalidator.ExactlyOneOf(
											path.MatchRelative().AtParent().AtName("password"),
											path.MatchRelative().AtParent().AtName("password_wo"),
										),
									},
								},
							},
							"password_wo": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The password (WO).",
									Optional:            true,
									WriteOnly:           true,
									Validators: []validator.String{
										stringvalidator.ConflictsWith(
											path.MatchRelative().AtParent().AtName("password"),
										),
										stringvalidator.AlsoRequires(
											path.MatchRelative().AtParent().AtName("password_wo_version"),
										),
										stringvalidator.ExactlyOneOf(
											path.MatchRelative().AtParent().AtName("password"),
											path.MatchRelative().AtParent().AtName("password_wo"),
										),
									},
								},
							},
							"password_wo_version": superschema.Int32Attribute{
								Resource: &schemaR.Int32Attribute{
									MarkdownDescription: "The version of the password_wo.",
									Optional:            true,
									PlanModifiers: []planmodifier.Int32{
										int32planmodifier.RequiresReplace(),
									},
									Validators: []validator.Int32{
										int32validator.ConflictsWith(
											path.MatchRelative().AtParent().AtName("password"),
										),
										int32validator.AlsoRequires(
											path.MatchRelative().AtParent().AtName("password_wo"),
										),
									},
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
