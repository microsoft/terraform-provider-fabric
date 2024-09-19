// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace_test

import (
	"regexp"
	"testing"
	"time"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testDataSourceWorkspaceGit       = testhelp.DataSourceFQN("fabric", workspaceGitTFName, "test")
	testDataSourceWorkspaceGitHeader = at.DataSourceHeader(testhelp.TypeName("fabric", workspaceGitTFName), "test")
)

func TestUnit_WorkspaceGitDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := NewRandomGitConnection()

	fakes.FakeServer.ServerFactory.Core.GitServer.GetConnection = fakeGitGetConnection(entity)

	resource.Test(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceWorkspaceGitHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// error - workspace_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceWorkspaceGitHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceWorkspaceGitHeader,
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
				testDataSourceWorkspaceGitHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceWorkspaceGit, "workspace_id", workspaceID),
				resource.TestCheckResourceAttrPtr(testDataSourceWorkspaceGit, "git_connection_state", (*string)(entity.GitConnectionState)),
				resource.TestCheckResourceAttrPtr(testDataSourceWorkspaceGit, "git_sync_details.head", entity.GitSyncDetails.Head),
				resource.TestCheckResourceAttr(testDataSourceWorkspaceGit, "git_sync_details.last_sync_time", entity.GitSyncDetails.LastSyncTime.Format(time.RFC3339)),
				resource.TestCheckResourceAttrSet(testDataSourceWorkspaceGit, "git_provider_details.git_provider_type"),
				resource.TestCheckResourceAttrSet(testDataSourceWorkspaceGit, "git_provider_details.organization_name"),
				resource.TestCheckResourceAttrSet(testDataSourceWorkspaceGit, "git_provider_details.project_name"),
				resource.TestCheckResourceAttrSet(testDataSourceWorkspaceGit, "git_provider_details.repository_name"),
				resource.TestCheckResourceAttrSet(testDataSourceWorkspaceGit, "git_provider_details.branch_name"),
				resource.TestCheckResourceAttrSet(testDataSourceWorkspaceGit, "git_provider_details.directory_name"),
			),
		},
	}))
}
