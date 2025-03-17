// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	superobjectvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/objectvalidator"
	superstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

const (
	CredentialTypeOnPremisesGateway fabcore.CredentialType = "OnPremisesGateway"
	// CredentialTypeOnPremisesGatewayPersonal fabcore.CredentialType = "OnPremisesGatewayPersonal".
)

func getResourceConnectionAttributes(ctx context.Context) map[string]schema.Attribute {
	possibleConnectivityTypeValues := utils.RemoveSlicesByValues(fabcore.PossibleConnectivityTypeValues(), []fabcore.ConnectivityType{fabcore.ConnectivityTypeOnPremisesGateway, fabcore.ConnectivityTypeOnPremisesGatewayPersonal, fabcore.ConnectivityTypePersonalCloud})

	possibleCredentialTypeValues := fabcore.PossibleCredentialTypeValues()
	possibleCredentialTypeValues = append(possibleCredentialTypeValues, CredentialTypeOnPremisesGateway)
	// possibleCredentialTypeValues = append(possibleCredentialTypeValues, CredentialTypeOnPremisesGatewayPersonal)

	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The object ID of the connection.",
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "The display name of the connection.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(123),
			},
		},
		"connectivity_type": schema.StringAttribute{
			MarkdownDescription: "The connectivity type of the connection. Accepted values: " + utils.ConvertStringSlicesToString(possibleConnectivityTypeValues, true, true),
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(possibleConnectivityTypeValues, false)...),
			},
		},
		"privacy_level": schema.StringAttribute{
			MarkdownDescription: "The privacy level of the connection. Accepted values: " + utils.ConvertStringSlicesToString(fabcore.PossiblePrivacyLevelValues(), true, true),
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossiblePrivacyLevelValues(), false)...),
			},
		},
		"gateway_id": schema.StringAttribute{
			MarkdownDescription: "The gateway object ID of the connection.",
			Optional:            true,
			CustomType:          customtypes.UUIDType{},
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
		"connection_details": schema.SingleNestedAttribute{
			MarkdownDescription: "The connection details of the connection.",
			Required:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[rsConnectionDetailsModel](ctx),
			Attributes: map[string]schema.Attribute{
				"path": schema.StringAttribute{
					MarkdownDescription: "The path of the connection.",
					Computed:            true,
				},
				"type": schema.StringAttribute{
					MarkdownDescription: "The type of the connection.",
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"creation_method": schema.StringAttribute{
					MarkdownDescription: "The creation method used to create the connection.",
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"parameters": schema.MapAttribute{
					MarkdownDescription: "A map of key/value pairs of connection parameters.",
					Optional:            true,
					CustomType:          supertypes.NewMapTypeOf[string](ctx),
				},
			},
		},
		"credential_details": schema.SingleNestedAttribute{
			MarkdownDescription: "The credential details of the connection.",
			Required:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[rsCredentialDetailsModel](ctx),
			Attributes: map[string]schema.Attribute{
				"connection_encryption": schema.StringAttribute{
					MarkdownDescription: "The connection encryption type of the connection. Accepted values: " + utils.ConvertStringSlicesToString(fabcore.PossibleConnectionEncryptionValues(), true, true),
					Required:            true,
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleConnectionEncryptionValues(), false)...),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"single_sign_on_type": schema.StringAttribute{
					MarkdownDescription: "The single sign-on type of the connection. Accepted values: " + utils.ConvertStringSlicesToString(fabcore.PossibleSingleSignOnTypeValues(), true, true),
					Required:            true,
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleSingleSignOnTypeValues(), false)...),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"skip_test_connection": schema.BoolAttribute{
					MarkdownDescription: "Whether the connection should skip the test connection during creation and update. `True` - Skip the test connection, `False` - Do not skip the test connection.",
					Required:            true,
				},
				"credential_type": schema.StringAttribute{
					MarkdownDescription: "The credential type of the connection. Possible values: " + utils.ConvertStringSlicesToString(possibleCredentialTypeValues, true, true),
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(possibleCredentialTypeValues, false)...),
						// superstringvalidator.NullIfAttributeIsOneOf(
						// 	path.MatchRoot("connectivity_type"),
						// 	[]attr.Value{types.StringValue(string(fabcore.ConnectivityTypeOnPremisesGateway))},
						// ),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"basic_credentials": schema.SingleNestedAttribute{
					MarkdownDescription: "The basic credentials.",
					Optional:            true,
					CustomType:          supertypes.NewSingleNestedObjectTypeOf[credentialsBasicModel](ctx),
					Attributes: map[string]schema.Attribute{
						"username": schema.StringAttribute{
							MarkdownDescription: "The username.",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"password": schema.StringAttribute{
							MarkdownDescription: "The password.",
							Optional:            true,
							Sensitive:           true,
							Validators: []validator.String{
								// Throws a warning diagnostic encouraging practitioners to use
								// password_wo if password has a known value.
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
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"password_wo": schema.StringAttribute{
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
						"password_wo_version": schema.StringAttribute{
							MarkdownDescription: "The version of the password_wo.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("password"),
								),
								stringvalidator.AlsoRequires(
									path.MatchRelative().AtParent().AtName("password_wo"),
								),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
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
				"key_credentials": schema.SingleNestedAttribute{
					MarkdownDescription: "The key credentials.",
					Optional:            true,
					CustomType:          supertypes.NewSingleNestedObjectTypeOf[credentialsKeyModel](ctx),
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							MarkdownDescription: "The key.",
							Optional:            true,
							Sensitive:           true,
							Validators: []validator.String{
								// Throws a warning diagnostic encouraging practitioners to use
								// key_wo if key has a known value.
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
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"key_wo": schema.StringAttribute{
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
						"key_wo_version": schema.StringAttribute{
							MarkdownDescription: "The version of the key_wo.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("key"),
								),
								stringvalidator.AlsoRequires(
									path.MatchRelative().AtParent().AtName("key_wo"),
								),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
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
				"service_principal_credentials": schema.SingleNestedAttribute{
					MarkdownDescription: "The service principal credentials.",
					Optional:            true,
					CustomType:          supertypes.NewSingleNestedObjectTypeOf[credentialsServicePrincipalModel](ctx),
					Attributes: map[string]schema.Attribute{
						"tenant_id": schema.StringAttribute{
							MarkdownDescription: "The tenant ID.",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"client_id": schema.StringAttribute{
							MarkdownDescription: "The client ID.",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"client_secret": schema.StringAttribute{
							MarkdownDescription: "The client secret.",
							Optional:            true,
							Sensitive:           true,
							Validators: []validator.String{
								// Throws a warning diagnostic encouraging practitioners to use
								// client_secret_wo if client_secret has a known value.
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
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"client_secret_wo": schema.StringAttribute{
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
						"client_secret_wo_version": schema.StringAttribute{
							MarkdownDescription: "The version of the client_secret_wo.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("client_secret"),
								),
								stringvalidator.AlsoRequires(
									path.MatchRelative().AtParent().AtName("client_secret_wo"),
								),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
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
				"shared_access_signature_credentials": schema.SingleNestedAttribute{
					MarkdownDescription: "The shared access signature credentials.",
					Optional:            true,
					CustomType:          supertypes.NewSingleNestedObjectTypeOf[credentialsSharedAccessSignatureModel](ctx),
					Attributes: map[string]schema.Attribute{
						"token": schema.StringAttribute{
							MarkdownDescription: "The token.",
							Optional:            true,
							Sensitive:           true,
							Validators: []validator.String{
								// Throws a warning diagnostic encouraging practitioners to use
								// token_wo if token has a known value.
								stringvalidator.PreferWriteOnlyAttribute(
									path.MatchRelative().AtParent().AtName("token_wo"),
								),
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("token_wo"),
								),
								stringvalidator.ExactlyOneOf(
									path.MatchRelative().AtParent().AtName("token"),
									path.MatchRelative().AtParent().AtName("token_wo"),
								),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"token_wo": schema.StringAttribute{
							MarkdownDescription: "The token (WO).",
							Optional:            true,
							WriteOnly:           true,
							Validators: []validator.String{
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("token"),
								),
								stringvalidator.AlsoRequires(
									path.MatchRelative().AtParent().AtName("token_wo_version"),
								),
								stringvalidator.ExactlyOneOf(
									path.MatchRelative().AtParent().AtName("token"),
									path.MatchRelative().AtParent().AtName("token_wo"),
								),
							},
						},
						"token_wo_version": schema.StringAttribute{
							MarkdownDescription: "The version of the token_wo.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("token"),
								),
								stringvalidator.AlsoRequires(
									path.MatchRelative().AtParent().AtName("token_wo"),
								),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
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
				"windows_credentials": schema.SingleNestedAttribute{
					MarkdownDescription: "The Windows credentials.",
					Optional:            true,
					CustomType:          supertypes.NewSingleNestedObjectTypeOf[credentialsWindowsModel](ctx),
					Attributes: map[string]schema.Attribute{
						"username": schema.StringAttribute{
							MarkdownDescription: "The username.",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"password": schema.StringAttribute{
							MarkdownDescription: "The password.",
							Optional:            true,
							Sensitive:           true,
							Validators: []validator.String{
								// Throws a warning diagnostic encouraging practitioners to use password_wo if password has a known value.
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
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"password_wo": schema.StringAttribute{
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
						"password_wo_version": schema.StringAttribute{
							MarkdownDescription: "The version of the password_wo.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("password"),
								),
								stringvalidator.AlsoRequires(
									path.MatchRelative().AtParent().AtName("password_wo"),
								),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
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

				// "encrypted_credentials": schema.SingleNestedAttribute{
				// 	MarkdownDescription: "The encrypted serialized .json of the list of name value pairs. Name is a credential name and value is a credential value. Encryption is performed using the Rivest-Shamir-Adleman (RSA) encryption algorithm with the on-premises gateway member's public key.",
				// 	Optional:            true,
				// 	CustomType:          supertypes.NewSingleNestedObjectTypeOf[credentialsEncryptedModel](ctx),
				// 	Attributes: map[string]schema.Attribute{
				// 		"value": schema.StringAttribute{
				// 			MarkdownDescription: "The value.",
				// 			Required:            true,
				// 			Sensitive:           true,
				// 		},
				// 	},
				// 	Validators: []validator.Object{
				// 		objectvalidator.ConflictsWith(
				// 			path.MatchRelative().AtParent().AtName("basic_credentials"),
				// 			path.MatchRelative().AtParent().AtName("key_credentials"),
				// 			path.MatchRelative().AtParent().AtName("service_principal_credentials"),
				// 			path.MatchRelative().AtParent().AtName("windows_credentials"),
				// 			path.MatchRelative().AtParent().AtName("shared_access_signature_credentials"),
				// 		),
				// 		superobjectvalidator.RequireIfAttributeIsOneOf(
				// 			path.MatchRoot("connectivity_type"),
				// 			[]attr.Value{
				// 				types.StringValue(string(fabcore.ConnectivityTypeOnPremisesGateway)),
				// 				// types.StringValue(string(fabcore.ConnectivityTypeOnPremisesGatewayPersonal)),
				// 			},
				// 		),
				// 	},
				// },
			},
		},
	}
}
