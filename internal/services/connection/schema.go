// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

//nolint:maintidx
func itemSchema(isList bool) superschema.Schema { //revive:disable-line:flag-parameter
	// Define possible values for enums
	connectivityTypeValues := []string{"ShareableCloud", "OnPremisesGateway", "VirtualNetworkGateway"}
	privacyLevelValues := []string{"Organizational", "Private", "Public", "None"}
	singleSignOnTypeValues := []string{"None", "Kerberos", "KerberosDirectQueryAndRefresh", "MicrosoftEntraID", "SecurityAssertionMarkupLanguage"}
	connectionEncryptionValues := []string{"Any", "Encrypted", "NotEncrypted"}
	credentialTypeValues := []string{"Anonymous", "Basic", "Key", "OAuth2", "ServicePrincipal", "SharedAccessSignature", "Windows", "WindowsWithoutImpersonation", "WorkspaceIdentity"}
	parameterDataTypeValues := []string{"Text", "Number", "Boolean", "Date", "DateTime", "Time", "Duration", "DateTimeZone"}

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
						stringvalidator.LengthAtMost(200),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"connectivity_type": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Connectivity type. Possible values: " + utils.ConvertStringSlicesToString(connectivityTypeValues, true, true) + ".",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf(connectivityTypeValues...),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"gateway_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Gateway ID. Required for OnPremisesGateway and VirtualNetworkGateway connectivity types.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"privacy_level": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Privacy level. Possible values: " + utils.ConvertStringSlicesToString(privacyLevelValues, true, true) + ".",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Computed: true,
					Default:  stringdefault.StaticString("Organizational"),
					Validators: []validator.String{
						stringvalidator.OneOf(privacyLevelValues...),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"connection_details": superschema.SuperSingleNestedAttributeOf[connectionDetailsModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Connection details.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Required: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"type": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Connection type.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"creation_method": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Creation method.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"parameters": superschema.SuperListNestedAttributeOf[connectionParameterModel]{
						Common: &schemaR.ListNestedAttribute{
							MarkdownDescription: "Connection parameters.",
						},
						Resource: &schemaR.ListNestedAttribute{
							Required: true,
						},
						DataSource: &schemaD.ListNestedAttribute{
							Computed: true,
						},
						Attributes: map[string]superschema.Attribute{
							"name": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Parameter name.",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"data_type": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Data type. Possible values: " + utils.ConvertStringSlicesToString(parameterDataTypeValues, true, true) + ".",
									Validators: []validator.String{
										stringvalidator.OneOf(parameterDataTypeValues...),
									},
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"value": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Parameter value.",
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
			"credential_details": superschema.SuperSingleNestedAttribute{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Credential details.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Required: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"single_sign_on_type": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Single sign-on type. Possible values: " + utils.ConvertStringSlicesToString(singleSignOnTypeValues, true, true) + ".",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  stringdefault.StaticString("None"),
							Validators: []validator.String{
								stringvalidator.OneOf(singleSignOnTypeValues...),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"connection_encryption": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Connection encryption. Possible values: " + utils.ConvertStringSlicesToString(connectionEncryptionValues, true, true) + ".",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  stringdefault.StaticString("NotEncrypted"),
							Validators: []validator.String{
								stringvalidator.OneOf(connectionEncryptionValues...),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"skip_test_connection": superschema.BoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Skip test connection.",
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
					"credentials": superschema.SuperSingleNestedAttribute{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "Credentials. Required fields depend on the credential_type value.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Required: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: map[string]superschema.Attribute{
							"credential_type": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Credential type. Possible values: " + utils.ConvertStringSlicesToString(credentialTypeValues, true, true) + ".",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf(credentialTypeValues...),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							// For BasicCredentials
							"username": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Username for Basic authentication. Required when credential_type is 'Basic' or 'Windows'.",
									Sensitive:           true,
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"password": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Password for Basic authentication. Required when credential_type is 'Basic' or 'Windows'.",
									Sensitive:           true,
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true, // Changed from false to true - not returned for security reasons
								},
							},
							// For KeyCredentials
							"key": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Key for Key authentication. Required when credential_type is 'Key'.",
									Sensitive:           true,
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							// For ServicePrincipalCredentials
							"application_id": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Application ID for ServicePrincipal authentication. Required when credential_type is 'ServicePrincipal'.",
									Sensitive:           true,
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"application_secret": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Application Secret for ServicePrincipal authentication. Required when credential_type is 'ServicePrincipal'.",
									Sensitive:           true,
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true, // Changed from false to true
								},
							},
							"tenant_id": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Tenant ID for ServicePrincipal authentication. Required when credential_type is 'ServicePrincipal'.",
									Sensitive:           true,
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							// For SharedAccessSignatureCredentials
							"sas_token": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "SAS token for SharedAccessSignature authentication. Required when credential_type is 'SharedAccessSignature'.",
									Sensitive:           true,
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true, // Changed from false to true - not returned for security reasons
								},
							},
							// For WindowsCredentials and WindowsWithoutImpersonationCredentials
							"domain": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Domain for Windows authentication. Required when credential_type is 'Windows' or 'WindowsWithoutImpersonation'.",
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}
