// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacegit_test

import (
	"regexp"
	"testing"
	"time"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_WorkspaceGitDataSource_AzDO(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := NewRandomGitConnection(fabcore.GitProviderTypeAzureDevOps)
	entityCredentials := NewRandomGitCredentialsResponse(fabcore.GitCredentialsSourceAutomatic)

	fakes.FakeServer.ServerFactory.Core.GitServer.GetConnection = fakeGitGetConnection(entity)
	fakes.FakeServer.ServerFactory.Core.GitServer.GetMyGitCredentials = fakeGitGetMyGitCredentials(entityCredentials)

	resource.Test(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// error - workspace_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":    workspaceID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "git_connection_state", (*string)(entity.GitConnectionState)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "git_sync_details.head", entity.GitSyncDetails.Head),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "git_sync_details.last_sync_time", entity.GitSyncDetails.LastSyncTime.Format(time.RFC3339)),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "git_provider_details.git_provider_type"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "git_provider_details.git_provider_type", string(fabcore.GitProviderTypeAzureDevOps)),
				resource.TestCheckNoResourceAttr(testDataSourceItemFQN, "git_provider_details.owner_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "git_provider_details.organization_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "git_provider_details.project_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "git_provider_details.repository_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "git_provider_details.branch_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "git_provider_details.directory_name"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "git_credentials.source", string(fabcore.GitCredentialsSourceAutomatic)),
				resource.TestCheckNoResourceAttr(testDataSourceItemFQN, "git_credentials.connection_id"),
			),
		},
	}))
}

func TestUnit_WorkspaceGitDataSource_GitHub(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := NewRandomGitConnection(fabcore.GitProviderTypeGitHub)
	entityCredentials := NewRandomGitCredentialsResponse(fabcore.GitCredentialsSourceConfiguredConnection)

	fakes.FakeServer.ServerFactory.Core.GitServer.GetConnection = fakeGitGetConnection(entity)
	fakes.FakeServer.ServerFactory.Core.GitServer.GetMyGitCredentials = fakeGitGetMyGitCredentials(entityCredentials)

	resource.Test(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// error - workspace_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":    workspaceID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "git_connection_state", (*string)(entity.GitConnectionState)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "git_sync_details.head", entity.GitSyncDetails.Head),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "git_sync_details.last_sync_time", entity.GitSyncDetails.LastSyncTime.Format(time.RFC3339)),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "git_provider_details.git_provider_type"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "git_provider_details.git_provider_type", string(fabcore.GitProviderTypeGitHub)),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "git_provider_details.owner_name"),
				resource.TestCheckNoResourceAttr(testDataSourceItemFQN, "git_provider_details.organization_name"),
				resource.TestCheckNoResourceAttr(testDataSourceItemFQN, "git_provider_details.project_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "git_provider_details.repository_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "git_provider_details.branch_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "git_provider_details.directory_name"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "git_credentials.source", string(fabcore.GitCredentialsSourceConfiguredConnection)),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "git_credentials.connection_id"),
			),
		},
	}))
}
