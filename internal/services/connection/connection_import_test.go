// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package connection_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"testing"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

func connectionImportMetadata(server *fabfake.ServerFactory) {
	server.Core.ConnectionsServer.NewListSupportedConnectionTypesPager = func(_ *fabcore.ConnectionsClientListSupportedConnectionTypesOptions) azfake.PagerResponder[fabcore.ConnectionsClientListSupportedConnectionTypesResponse] {
		var resp azfake.PagerResponder[fabcore.ConnectionsClientListSupportedConnectionTypesResponse]
		parameters := []fabcore.ConnectionCreationParameter{
			{Name: new("server"), DataType: new(fabcore.DataTypeText), Required: new(false)},
			{Name: new("database"), DataType: new(fabcore.DataTypeText), Required: new(false)},
		}
		resp.AddPage(http.StatusOK, fabcore.ConnectionsClientListSupportedConnectionTypesResponse{
			Value: []fabcore.ConnectionCreationMetadata{{
				Type: new("FTP"),
				CreationMethods: []fabcore.ConnectionCreationMethod{
					{Name: new("FTP.Contents"), Parameters: parameters},
					{Name: new("FTP.Alternative"), Parameters: parameters},
				},
				SupportedCredentialTypes: []fabcore.CredentialType{
					fabcore.CredentialTypeAnonymous, fabcore.CredentialTypeBasic,
					fabcore.CredentialTypeKey, fabcore.CredentialTypeKeyPair,
					fabcore.CredentialTypeServicePrincipal, fabcore.CredentialTypeSharedAccessSignature,
				},
				SupportedConnectionEncryptionTypes: []fabcore.ConnectionEncryption{fabcore.ConnectionEncryptionNotEncrypted},
				SupportsSkipTestConnection:         new(true),
			}},
		}, nil)

		return resp
	}
}

func TestUnit_ConnectionResource_ImportWithoutSecrets(t *testing.T) {
	for _, tc := range []struct {
		credentialType fabcore.CredentialType
		block          string
		attributes     map[string]any
	}{
		{fabcore.CredentialTypeBasic, "basic_credentials", map[string]any{"username": "existing-user"}},
		{fabcore.CredentialTypeKey, "key_credentials", map[string]any{}},
		{fabcore.CredentialTypeServicePrincipal, "service_principal_credentials", map[string]any{"tenant_id": testhelp.RandomUUID(), "client_id": testhelp.RandomUUID()}},
		{fabcore.CredentialTypeSharedAccessSignature, "shared_access_signature_credentials", map[string]any{}},
		{fabcore.CredentialTypeKeyPair, "key_pair_credentials", map[string]any{"identifier": "existing-key"}},
	} {
		t.Run(string(tc.credentialType), func(t *testing.T) {
			server := fakes.NewFakeServer()
			connectionImportMetadata(server.ServerFactory)

			entity := fakes.NewRandomShareableCloudConnection()
			entity.GatewayID = new(testhelp.RandomUUID())
			entity.ConnectionDetails = &fabcore.ListConnectionDetails{Type: new("FTP"), Path: new("existing.example.com")}
			entity.CredentialDetails = &fabcore.ListCredentialDetails{
				CredentialType: new(tc.credentialType), ConnectionEncryption: new(fabcore.ConnectionEncryptionNotEncrypted),
				SingleSignOnType: new(fabcore.SingleSignOnTypeNone), SkipTestConnection: new(false),
			}
			server.Upsert(entity)
			update := server.ServerFactory.Core.ConnectionsServer.UpdateConnection
			server.ServerFactory.Core.ConnectionsServer.UpdateConnection = func(ctx context.Context, id string, req fabcore.UpdateConnectionRequestClassification, opts *fabcore.ConnectionsClientUpdateConnectionOptions) (azfake.Responder[fabcore.ConnectionsClientUpdateConnectionResponse], azfake.ErrorResponder) {
				details := req.(*fabcore.UpdateShareableCloudConnectionRequest).CredentialDetails
				if details == nil || details.Credentials != nil {
					t.Error("secretless updates must omit credentials while preserving credential settings")
				}

				return update(ctx, id, req, opts)
			}
			config := at.CompileConfig(testResourceItemHeader, map[string]any{
				"display_name": *entity.DisplayName, "connectivity_type": "ShareableCloud", "privacy_level": string(*entity.PrivacyLevel),
				"connection_details": map[string]any{"type": "FTP", "creation_method": "FTP.Contents", "parameters": []map[string]any{{"name": "server", "value": "existing.example.com"}}},
				"credential_details": map[string]any{"credential_type": string(tc.credentialType), tc.block: tc.attributes},
			})
			importConfig := config + fmt.Sprintf("\nimport {\n to = %s\n id = %q\n}\n", testResourceItemFQN, *entity.ID)
			renamedConfig := regexp.MustCompile(regexp.QuoteMeta(*entity.DisplayName)).ReplaceAllString(config, "renamed-"+*entity.DisplayName)
			changedCredentialConfig := at.CompileConfig(testResourceItemHeader, map[string]any{
				"display_name": "renamed-" + *entity.DisplayName, "connectivity_type": "ShareableCloud", "privacy_level": string(*entity.PrivacyLevel),
				"connection_details": map[string]any{"type": "FTP", "creation_method": "FTP.Contents", "parameters": []map[string]any{{"name": "server", "value": "existing.example.com"}}},
				"credential_details": map[string]any{"credential_type": "Basic", "basic_credentials": map[string]any{"username": "changed-user"}},
			})
			resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, server.ServerFactory, nil, []resource.TestStep{
				{Config: config, PlanOnly: true, ExpectError: regexp.MustCompile("Missing credential secret")},
				{
					Config: importConfig, PlanOnly: true, ExpectNonEmptyPlan: true,
					ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testResourceItemFQN, plancheck.ResourceActionUpdate),
					}},
				},
				{
					Config: importConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(testResourceItemFQN, "id", *entity.ID),
						resource.TestCheckResourceAttr(testResourceItemFQN, "connection_details.creation_method", "FTP.Contents"),
						resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.credential_type", string(tc.credentialType)),
					),
				},
				{Config: config, PlanOnly: true},
				{
					Config: renamedConfig,
					Check:  resource.TestCheckResourceAttr(testResourceItemFQN, "id", *entity.ID),
				},
				{
					Config:   changedCredentialConfig,
					PlanOnly: true, ExpectError: regexp.MustCompile("Missing credential secret"),
				},
				{
					Config:   regexp.MustCompile("FTP.Contents").ReplaceAllString(renamedConfig, "FTP.Alternative"),
					PlanOnly: true, ExpectError: regexp.MustCompile("Missing credential secret"),
				},
			}))
		})
	}
}

func TestUnit_ConnectionResource_ImportImmutableDetails(t *testing.T) {
	for _, withParameters := range []bool{false, true} {
		t.Run(strconv.FormatBool(withParameters), func(t *testing.T) {
			server := fakes.NewFakeServer()
			connectionImportMetadata(server.ServerFactory)

			entity := fakes.NewRandomShareableCloudConnection()
			entity.ConnectionDetails = &fabcore.ListConnectionDetails{Type: new("FTP"), Path: new("existing.example.com")}
			entity.CredentialDetails = &fabcore.ListCredentialDetails{
				CredentialType: new(fabcore.CredentialTypeAnonymous), ConnectionEncryption: new(fabcore.ConnectionEncryptionNotEncrypted),
				SingleSignOnType: new(fabcore.SingleSignOnTypeNone), SkipTestConnection: new(false),
			}
			server.Upsert(entity)

			details := map[string]any{"type": "FTP", "creation_method": "FTP.Contents"}
			if withParameters {
				details["parameters"] = []map[string]any{{"name": "server", "value": "existing.example.com"}}
			}

			attributes := map[string]any{
				"display_name": *entity.DisplayName, "connectivity_type": "ShareableCloud", "privacy_level": string(*entity.PrivacyLevel),
				"connection_details": details, "credential_details": map[string]any{"credential_type": "Anonymous"},
			}
			config := at.CompileConfig(testResourceItemHeader, attributes)
			details["parameters"] = []map[string]any{{"name": "server", "value": "changed.example.com"}}
			changedParameters := at.CompileConfig(testResourceItemHeader, attributes)
			details["parameters"] = []map[string]any{{"name": "database", "value": "other"}}
			changedName := at.CompileConfig(testResourceItemHeader, attributes)

			delete(details, "parameters")

			removedParameters := at.CompileConfig(testResourceItemHeader, attributes)
			details["creation_method"] = "FTP.Alternative"
			changedMethod := at.CompileConfig(testResourceItemHeader, attributes)

			steps := []resource.TestStep{
				{
					Config:           config + fmt.Sprintf("\nimport {\n to = %s\n id = %q\n}\n", testResourceItemFQN, *entity.ID),
					ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(testResourceItemFQN, plancheck.ResourceActionUpdate)}},
					Check:            resource.TestCheckResourceAttr(testResourceItemFQN, "id", *entity.ID),
				},
				{Config: config, PlanOnly: true},
			}
			for _, changed := range []string{changedParameters, changedName, changedMethod} {
				steps = append(steps, resource.TestStep{
					Config: changed, PlanOnly: true, ExpectNonEmptyPlan: true,
					ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectResourceAction(testResourceItemFQN, plancheck.ResourceActionReplace)}},
				})
			}

			if withParameters {
				steps = append(steps, resource.TestStep{
					Config: removedParameters, PlanOnly: true, ExpectNonEmptyPlan: true,
					ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectResourceAction(testResourceItemFQN, plancheck.ResourceActionReplace)}},
				})
			}

			resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, server.ServerFactory, nil, steps))
		})
	}
}

func TestUnit_ConnectionResource_ImportIgnoredDetails(t *testing.T) {
	for _, credentialResponse := range []string{"complete", "type-only", "absent"} {
		t.Run(credentialResponse, func(t *testing.T) {
			server := fakes.NewFakeServer()
			entity := fakes.NewRandomShareableCloudConnection()
			entity.ConnectionDetails.Type = new("FTP")

			entity.CredentialDetails = &fabcore.ListCredentialDetails{
				CredentialType: new(fabcore.CredentialTypeBasic), ConnectionEncryption: new(fabcore.ConnectionEncryptionNotEncrypted),
				SingleSignOnType: new(fabcore.SingleSignOnTypeNone), SkipTestConnection: new(false),
			}
			if credentialResponse == "type-only" {
				entity.CredentialDetails = &fabcore.ListCredentialDetails{CredentialType: new(fabcore.CredentialTypeBasic)}
			}

			if credentialResponse == "absent" {
				entity.CredentialDetails = nil
			}

			action := plancheck.ResourceActionNoop
			if credentialResponse == "type-only" {
				action = plancheck.ResourceActionUpdate
			}

			server.Upsert(entity)
			config := fmt.Sprintf(`
%s {
 display_name = %q
 connectivity_type = "ShareableCloud"
 privacy_level = %q
 connection_details = { type = "FTP", creation_method = "FTP.Contents" }
 credential_details = { credential_type = "Basic", basic_credentials = { username = "existing" } }
 lifecycle { ignore_changes = [connection_details, credential_details] }
}
import {
 to = %s
 id = %q
}
`, testResourceItemHeader, *entity.DisplayName, *entity.PrivacyLevel, testResourceItemFQN, *entity.ID)
			resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, server.ServerFactory, nil, []resource.TestStep{
				{
					Config:           config,
					ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(testResourceItemFQN, action)}},
					Check:            resource.TestCheckResourceAttr(testResourceItemFQN, "id", *entity.ID),
				},
				{
					Config: regexp.MustCompile(regexp.QuoteMeta(*entity.DisplayName)).ReplaceAllString(config, "renamed-"+*entity.DisplayName),
					Check:  resource.TestCheckResourceAttr(testResourceItemFQN, "id", *entity.ID),
				},
			}))
		})
	}
}

func TestUnit_ConnectionResource_ImportVirtualNetworkGateway(t *testing.T) {
	server := fakes.NewFakeServer()
	connectionImportMetadata(server.ServerFactory)

	entity := fakes.NewRandomVirtualNetworkGatewayConnection()
	entity.ConnectionDetails = &fabcore.ListConnectionDetails{Type: new("FTP"), Path: new("existing.example.com")}
	entity.CredentialDetails = &fabcore.ListCredentialDetails{
		CredentialType: new(fabcore.CredentialTypeBasic), ConnectionEncryption: new(fabcore.ConnectionEncryptionNotEncrypted),
		SingleSignOnType: new(fabcore.SingleSignOnTypeNone), SkipTestConnection: new(false),
	}
	server.Upsert(entity)
	config := at.CompileConfig(testResourceItemHeader, map[string]any{
		"display_name": *entity.DisplayName, "connectivity_type": "VirtualNetworkGateway", "gateway_id": *entity.GatewayID,
		"privacy_level":      string(*entity.PrivacyLevel),
		"connection_details": map[string]any{"type": "FTP", "creation_method": "FTP.Contents", "parameters": []map[string]any{{"name": "server", "value": "existing.example.com"}}},
		"credential_details": map[string]any{"credential_type": "Basic", "basic_credentials": map[string]any{"username": "existing-user"}},
	})
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, server.ServerFactory, nil, []resource.TestStep{
		{
			Config:           config + fmt.Sprintf("\nimport {\n to = %s\n id = %q\n}\n", testResourceItemFQN, *entity.ID),
			ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(testResourceItemFQN, plancheck.ResourceActionUpdate)}},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "id", *entity.ID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "gateway_id", *entity.GatewayID),
			),
		},
		{Config: config, PlanOnly: true},
	}))
}

func TestUnit_ConnectionResource_UnknownCredentialReference(t *testing.T) {
	server := fakes.NewFakeServer()
	connectionImportMetadata(server.ServerFactory)

	vaultID := testhelp.RandomUUID()
	config := fmt.Sprintf(`
resource "terraform_data" "vault" {
 input = {
  connection_id = %q
  secret_name = "existing-secret"
 }
}
%s {
 display_name = %q
 connectivity_type = "ShareableCloud"
 connection_details = { type = "FTP", creation_method = "FTP.Contents" }
 credential_details = {
  credential_type = "Key"
  key_credentials = { key_reference = terraform_data.vault.output }
 }
}
`, vaultID, testResourceItemHeader, testhelp.RandomUUID())
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, server.ServerFactory, nil, []resource.TestStep{{
		Config: config,
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
			resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details.key_credentials.key_reference.connection_id", vaultID),
		),
	}}))
}

func TestUnit_ConnectionResource_CredentialSecretsCRUD(t *testing.T) {
	for _, tc := range []struct {
		credentialType string
		block          string
		secret         string
		version        string
		attributes     map[string]any
	}{
		{"Basic", "basic_credentials", "password_wo", "password_wo_version", map[string]any{"username": "user"}},
		{"Key", "key_credentials", "key_wo", "key_wo_version", map[string]any{}},
		{"KeyPair", "key_pair_credentials", "private_key_wo", "private_key_wo_version", map[string]any{"identifier": "key"}},
		{"ServicePrincipal", "service_principal_credentials", "client_secret_wo", "client_secret_wo_version", map[string]any{"tenant_id": testhelp.RandomUUID(), "client_id": testhelp.RandomUUID()}},
		{"SharedAccessSignature", "shared_access_signature_credentials", "token_wo", "token_wo_version", map[string]any{}},
	} {
		t.Run(tc.credentialType, func(t *testing.T) {
			server := fakes.NewFakeServer()
			connectionImportMetadata(server.ServerFactory)
			create := server.ServerFactory.Core.ConnectionsServer.CreateConnection
			server.ServerFactory.Core.ConnectionsServer.CreateConnection = func(ctx context.Context, req fabcore.CreateConnectionRequestClassification, opts *fabcore.ConnectionsClientCreateConnectionOptions) (azfake.Responder[fabcore.ConnectionsClientCreateConnectionResponse], azfake.ErrorResponder) {
				credentials := req.(*fabcore.CreateCloudConnectionRequest).CredentialDetails.Credentials
				if credentials == nil || string(*credentials.GetCredentials().CredentialType) != tc.credentialType {
					t.Error("create must send the configured credentials")
				}

				return create(ctx, req, opts)
			}
			update := server.ServerFactory.Core.ConnectionsServer.UpdateConnection
			server.ServerFactory.Core.ConnectionsServer.UpdateConnection = func(ctx context.Context, id string, req fabcore.UpdateConnectionRequestClassification, opts *fabcore.ConnectionsClientUpdateConnectionOptions) (azfake.Responder[fabcore.ConnectionsClientUpdateConnectionResponse], azfake.ErrorResponder) {
				credentials := req.(*fabcore.UpdateShareableCloudConnectionRequest).CredentialDetails.Credentials
				if credentials == nil || string(*credentials.GetCredentials().CredentialType) != tc.credentialType {
					t.Error("credential rotation must send the configured credentials")
				}

				return update(ctx, id, req, opts)
			}
			tc.attributes[tc.secret] = "test-secret"
			tc.attributes[tc.version] = 1
			attributes := map[string]any{
				"display_name": testhelp.RandomUUID(), "connectivity_type": "ShareableCloud",
				"connection_details": map[string]any{"type": "FTP", "creation_method": "FTP.Contents"},
				"credential_details": map[string]any{"credential_type": tc.credentialType, tc.block: tc.attributes},
			}
			createConfig := at.CompileConfig(testResourceItemHeader, attributes)
			tc.attributes[tc.secret] = "rotated-test-secret"
			tc.attributes[tc.version] = 2
			updateConfig := at.CompileConfig(testResourceItemHeader, attributes)
			resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, server.ServerFactory, nil, []resource.TestStep{
				{Config: createConfig, Check: resource.TestCheckResourceAttrSet(testResourceItemFQN, "id")},
				{
					Config:           updateConfig,
					ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(testResourceItemFQN, plancheck.ResourceActionUpdate)}},
					Check:            resource.TestCheckResourceAttr(testResourceItemFQN, "credential_details."+tc.block+"."+tc.version, "2"),
				},
			}))
		})
	}
}
