// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection_test

import (
	"os"
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

//nolint:maintidx
func TestUnit_ConnectionResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// step 1: error - missing attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// step 2: error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// step 3: error - missing required attributes - display_name
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
		},
		// step 4: error - missing required attributes - connectivity_type
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":  "test",
					"privacy_level": "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "connectivity_type" is required, but no definition was found.`),
		},
		// step 5: error - missing required attributes - connection_details
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "connection_details" is required, but no definition was found.`),
		},
		// step 6: error - missing required attributes - credential_details
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "credential_details" is required, but no definition was found.`),
		},
		// step 7: error - invalid attribute value - connectivity_type
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test",
					"connectivity_type": "InvalidType",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
		// step 8: error - invalid attribute value - connectivity_type (PersonalCloud)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test",
					"connectivity_type": "PersonalCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
		// step 9: error - invalid attribute value - privacy_level
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "InvalidLevel",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
		// step 10: error - invalid uuid - gateway_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test",
					"connectivity_type": "VirtualNetworkGateway",
					"privacy_level":     "Organizational",
					"gateway_id":        "invalid uuid",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// step 11: error - VirtualNetworkGateway connection missing required gateway_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test",
					"connectivity_type": "VirtualNetworkGateway",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute gateway_id`),
		},
		// step 12: error - VirtualNetworkGateway connection with AllowConnectionUsageInGateway set
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":                      "test",
					"connectivity_type":                 "VirtualNetworkGateway",
					"privacy_level":                     "Organizational",
					"gateway_id":                        testhelp.RandomUUID(),
					"allow_connection_usage_in_gateway": true,
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute allow_connection_usage_in_gateway`),
		},
		// step 13: error - ShareableCloud connection with gateway_id set
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"gateway_id":        testhelp.RandomUUID(),
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute gateway_id`),
		},
	}))
}

func TestUnit_ConnectionResource_WriteOnly_Credentials(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Step 1: Test with basic_credentials and another credential type that will fail
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-writeonly-basic",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeBasic),
						"basic_credentials": map[string]any{
							"username":            "testuser",
							"password_wo":         "testpassword",
							"password_wo_version": 1,
						},
						"key_credentials": map[string]any{
							"key_wo":         "invalid-key",
							"key_wo_version": 1,
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
		},
		// Step 2: Test credential type mismatch - Basic type but using Key credentials structure
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-writeonly-mismatch",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeKey),
						"basic_credentials": map[string]any{
							"username":            "testuser",
							"password_wo":         "testpassword",
							"password_wo_version": 1,
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute`),
		},
		// Step 3: Test correct basic credentials configuration
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-writeonly-correct",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeBasic),
						"basic_credentials": map[string]any{
							"username":            "testuser",
							"password_wo":         "testpassword",
							"password_wo_version": 1,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", "test-writeonly-correct"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connectivity_type", "ShareableCloud"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "privacy_level", "Organizational"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.credential_type", string(fabcore.CredentialTypeBasic)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.basic_credentials.username", "testuser"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.basic_credentials.password_wo_version", "1"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.skip_test_connection", "true"),
			),
		},
		// Step 4: Test password update with version = 2
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-writeonly-correct",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeBasic),
						"basic_credentials": map[string]any{
							"username":            "testuser",
							"password_wo":         "updatedpassword",
							"password_wo_version": 2,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", "test-writeonly-correct"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connectivity_type", "ShareableCloud"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "privacy_level", "Organizational"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.credential_type", string(fabcore.CredentialTypeBasic)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.basic_credentials.username", "testuser"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.basic_credentials.password_wo_version", "2"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.skip_test_connection", "true"),
			),
		},
	}))
}

func TestUnit_ConnectionResource_CRUD(t *testing.T) {
	entityExist := fakes.NewRandomShareableCloudConnection()
	entityBefore := fakes.NewRandomShareableCloudConnection()
	entityAfter := fakes.NewRandomShareableCloudConnection()

	fakes.FakeServer.Upsert(fakes.NewRandomShareableCloudConnection())
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(fakes.NewRandomShareableCloudConnection())

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - create - existing entity
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      *entityExist.DisplayName,
					"connectivity_type": string(fabcore.ConnectivityTypeShareableCloud),
					"privacy_level":     string(fabcore.PrivacyLevelOrganizational),
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorCreateHeader),
		},
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      *entityBefore.DisplayName,
					"connectivity_type": string(fabcore.ConnectivityTypeShareableCloud),
					"privacy_level":     string(*entityBefore.PrivacyLevel),
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(*entityBefore.CredentialDetails.ConnectionEncryption),
						"single_sign_on_type":   string(*entityBefore.CredentialDetails.SingleSignOnType),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connectivity_type", string(fabcore.ConnectivityTypeShareableCloud)),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "privacy_level", (*string)(entityBefore.PrivacyLevel)),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "credential_details.connection_encryption", (*string)(entityBefore.CredentialDetails.ConnectionEncryption)),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "credential_details.single_sign_on_type", (*string)(entityBefore.CredentialDetails.SingleSignOnType)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.credential_type", string(fabcore.CredentialTypeAnonymous)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "allow_connection_usage_in_gateway", "false"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "gateway_id"),
			),
		},
		// Update and Read - no replacement
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":                      *entityAfter.DisplayName,
					"connectivity_type":                 string(fabcore.ConnectivityTypeShareableCloud),
					"privacy_level":                     string(*entityAfter.PrivacyLevel),
					"allow_connection_usage_in_gateway": true,
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(*entityBefore.CredentialDetails.ConnectionEncryption),
						"single_sign_on_type":   string(*entityBefore.CredentialDetails.SingleSignOnType),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connectivity_type", string(fabcore.ConnectivityTypeShareableCloud)),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "privacy_level", (*string)(entityAfter.PrivacyLevel)),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "credential_details.connection_encryption", (*string)(entityBefore.CredentialDetails.ConnectionEncryption)),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "credential_details.single_sign_on_type", (*string)(entityBefore.CredentialDetails.SingleSignOnType)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.credential_type", string(fabcore.CredentialTypeAnonymous)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "allow_connection_usage_in_gateway", "true"),
			),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":                      *entityAfter.DisplayName,
					"connectivity_type":                 string(fabcore.ConnectivityTypeShareableCloud),
					"privacy_level":                     string(*entityAfter.PrivacyLevel),
					"allow_connection_usage_in_gateway": false,
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(*entityAfter.CredentialDetails.ConnectionEncryption),
						"single_sign_on_type":   string(*entityAfter.CredentialDetails.SingleSignOnType),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connectivity_type", string(fabcore.ConnectivityTypeShareableCloud)),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "privacy_level", (*string)(entityAfter.PrivacyLevel)),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "credential_details.connection_encryption", (*string)(entityBefore.CredentialDetails.ConnectionEncryption)),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "credential_details.single_sign_on_type", (*string)(entityBefore.CredentialDetails.SingleSignOnType)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.credential_type", string(fabcore.CredentialTypeAnonymous)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "allow_connection_usage_in_gateway", "false"),
			),
		},
		// Update connectivity type - replacement
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      *entityAfter.DisplayName,
					"connectivity_type": string(fabcore.ConnectivityTypeVirtualNetworkGateway),
					"gateway_id":        testhelp.RandomUUID(),
					"privacy_level":     string(*entityAfter.PrivacyLevel),
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(*entityBefore.CredentialDetails.ConnectionEncryption),
						"single_sign_on_type":   string(*entityBefore.CredentialDetails.SingleSignOnType),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connectivity_type", string(fabcore.ConnectivityTypeVirtualNetworkGateway)),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "privacy_level", (*string)(entityAfter.PrivacyLevel)),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "gateway_id"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "allow_connection_usage_in_gateway"),
			),
		},
		// Creation method with empty parameters list
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-empty-params",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "ConnectionWithEmptyParametersList",
						"creation_method": "ConnectionWithEmptyParametersList.Actions",
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", "test-empty-params"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connectivity_type", "ShareableCloud"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "privacy_level", "Organizational"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.credential_type", string(fabcore.CredentialTypeAnonymous)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.skip_test_connection", "true"),
			),
		},
	}))
}

func TestUnit_ConnectionResource_ModifyPlan_Validations(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Test 1: Unsupported connection type validation
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-connection",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "UnsupportedConnectionType",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Unsupported connection type`),
		},
		// Test 2: Unsupported creation method validation
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-connection",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "UnsupportedCreationMethod",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Unsupported creation method`),
		},
		// Test 3: Unsupported connection encryption validation
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-connection",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Unsupported connection encryption`),
		},
		// Test 4: Unsupported credential type validation
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-connection",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeWindows), // Unsupported because of the modify plan validation, there is also a static check in the schema
					},
				},
			),
			ExpectError: regexp.MustCompile(`Unsupported credential type`),
		},
		// Test 5: Unsupported skip test connection validation (false when should be true)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-connection",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Unsupported skip test connection`),
		},
		// Test 6: Unsupported connection parameter key validation
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-connection",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "unsupported_parameter",
								"value": "some_value",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Unsupported connection parameter key`),
		},
		// Test 7: Missing required connection parameter validation
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-connection",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "database", // Missing required 'server' parameter
								"value": "testdb",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Missing connection parameter key`),
		},
		// Test 8: Missing connection parameter value validation
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-connection",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "", // Empty value for required parameter
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Missing connection parameter value`),
		},
	}))
}

func TestUnit_ConnectionResource_ModifyPlan_DataTypeValidations(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Invalid boolean parameter data type validation
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-connection",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
							{
								"name":  "enable_ssl", // Assume this is a Boolean parameter in fake data
								"value": "maybe",      // Invalid boolean value
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid connection parameter value`),
		},
		// Invalid date parameter data type validation
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-connection",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
							{
								"name":  "start_date", // Assume this is a Date parameter in fake data
								"value": "invalid-date",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid connection parameter value`),
		},
		// Invalid number parameter data type validation
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-connection",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
							{
								"name":  "port", // Assume this is a Number parameter in fake data
								"value": "not-a-number",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid connection parameter value`),
		},
		// Valid configuration that should pass all validations
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      "test-connection-valid",
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
							{
								"name":  "database",
								"value": "testdb",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  true,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			// Should not expect any errors for valid configuration
		},
	}))
}

func TestAcc_ConnectionResource_ShareableCloud(t *testing.T) {
	displayName := testhelp.RandomName()
	displayNameUpdated := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      displayName,
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", displayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connectivity_type", "ShareableCloud"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "privacy_level", "Organizational"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "connection_details.path"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connection_details.type", "FTP"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.connection_encryption", string(fabcore.ConnectionEncryptionNotEncrypted)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.credential_type", string(fabcore.CredentialTypeAnonymous)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.single_sign_on_type", string(fabcore.SingleSignOnTypeNone)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.skip_test_connection", "false"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "allow_connection_usage_in_gateway", "false"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "gateway_id"),
			),
		},
		// Update display name and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      displayNameUpdated,
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", displayNameUpdated),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connectivity_type", "ShareableCloud"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "privacy_level", "Organizational"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "connection_details.path"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connection_details.type", "FTP"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.connection_encryption", string(fabcore.ConnectionEncryptionNotEncrypted)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.credential_type", string(fabcore.CredentialTypeAnonymous)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.single_sign_on_type", string(fabcore.SingleSignOnTypeNone)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.skip_test_connection", "false"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "allow_connection_usage_in_gateway", "false"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "gateway_id"),
			),
		},
	},
	))
}

func TestAcc_ConnectionResource_ShareableCloud_SQLServer_WriteOnly(t *testing.T) {
	var (
		sqlUsername, sqlUsernameExist = os.LookupEnv("FABRIC_CONNECTION_SQL_SERVER_USERNAME")
		sqlPassword, sqlPasswordExist = os.LookupEnv("FABRIC_CONNECTION_SQL_SERVER_PASSWORD")
		sqlURL, sqlURLEexist          = os.LookupEnv("FABRIC_CONNECTION_SQL_SERVER_URL")
	)

	var missingEnvVars []string

	if !sqlUsernameExist {
		missingEnvVars = append(missingEnvVars, "FABRIC_CONNECTION_SQL_SERVER_USERNAME")
	}

	if !sqlPasswordExist {
		missingEnvVars = append(missingEnvVars, "FABRIC_CONNECTION_SQL_SERVER_PASSWORD")
	}

	if !sqlURLEexist {
		missingEnvVars = append(missingEnvVars, "FABRIC_CONNECTION_SQL_SERVER_URL")
	}

	t.Fatalf("SQL Username is: %s", sqlUsername)
	t.Fatalf("SQL URL is: %s", sqlURL)

	if len(missingEnvVars) > 0 {
		t.Fatalf("Required environment variables are not set: %v", missingEnvVars)
	}

	displayName := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create SQL connection with basic credentials using environment variables
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      displayName,
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "SQL",
						"creation_method": "Sql",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": sqlURL,
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       "Basic",
						"basic_credentials": map[string]any{
							"username":            sqlUsername,
							"password_wo":         sqlPassword,
							"password_wo_version": 1,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", displayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connectivity_type", "ShareableCloud"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "privacy_level", "Organizational"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "connection_details.path"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connection_details.type", "SQL"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.connection_encryption", string(fabcore.ConnectionEncryptionNotEncrypted)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.credential_type", "Basic"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.single_sign_on_type", string(fabcore.SingleSignOnTypeNone)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.skip_test_connection", "false"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.basic_credentials.username", sqlUsername),
				resource.TestCheckResourceAttr(testResourceItemFQN, "allow_connection_usage_in_gateway", "false"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "gateway_id"),
			),
		},
	},
	))
}

func TestAcc_ConnectionResource_VirtualNetworkGateway(t *testing.T) {
	entityVirtualNetwork := testhelp.WellKnown()["GatewayVirtualNetwork"].(map[string]any)
	entityVirtualNetworkID := entityVirtualNetwork["id"].(string)
	displayName := testhelp.RandomName()
	displayNameUpdated := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      displayName,
					"connectivity_type": "VirtualNetworkGateway",
					"privacy_level":     "Organizational",
					"gateway_id":        entityVirtualNetworkID,
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", displayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connectivity_type", "VirtualNetworkGateway"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "privacy_level", "Organizational"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "connection_details.path"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connection_details.type", "FTP"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.connection_encryption", string(fabcore.ConnectionEncryptionNotEncrypted)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.credential_type", string(fabcore.CredentialTypeAnonymous)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.single_sign_on_type", string(fabcore.SingleSignOnTypeNone)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.skip_test_connection", "false"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "gateway_id", entityVirtualNetworkID),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "allow_connection_usage_in_gateway"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      displayNameUpdated,
					"connectivity_type": "VirtualNetworkGateway",
					"privacy_level":     "Organizational",
					"gateway_id":        entityVirtualNetworkID,
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
						"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
						"skip_test_connection":  false,
						"credential_type":       string(fabcore.CredentialTypeAnonymous),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", displayNameUpdated),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connectivity_type", "VirtualNetworkGateway"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "privacy_level", "Organizational"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "connection_details.path"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connection_details.type", "FTP"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.connection_encryption", string(fabcore.ConnectionEncryptionNotEncrypted)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.credential_type", string(fabcore.CredentialTypeAnonymous)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.single_sign_on_type", string(fabcore.SingleSignOnTypeNone)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.skip_test_connection", "false"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "gateway_id", entityVirtualNetworkID),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "allow_connection_usage_in_gateway"),
			),
		},
	},
	))
}
