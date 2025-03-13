// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testResourceWorkspaceGitFQN    = testhelp.ResourceFQN("fabric", workspaceGitTFName, "test")
	testResourceWorkspaceGitHeader = at.ResourceHeader(testhelp.TypeName("fabric", workspaceGitTFName), "test")
)

func TestUnit_WorkspaceGitResource_AzDO(t *testing.T) {
	gitConnection := NewRandomGitConnection(fabcore.GitProviderTypeAzureDevOps)
	gitCredentials := NewRandomGitCredentialsResponse(fabcore.GitCredentialsSourceAutomatic)
	gitProviderDetails := gitConnection.GitProviderDetails.GetGitProviderDetails()
	gitInit := NewRandomGitInitializeGitConnection()

	fakes.FakeServer.ServerFactory.Core.GitServer.GetConnection = fakeGitGetConnection(gitConnection)
	fakes.FakeServer.ServerFactory.Core.GitServer.GetMyGitCredentials = fakeGitGetMyGitCredentials(gitCredentials)
	fakes.FakeServer.ServerFactory.Core.GitServer.Connect = fakeGitConnect()
	fakes.FakeServer.ServerFactory.Core.GitServer.BeginInitializeConnection = fakeGitInitializeGitConnection(gitInit)
	fakes.FakeServer.ServerFactory.Core.GitServer.BeginCommitToGit = fakeGitCommitToGit()
	fakes.FakeServer.ServerFactory.Core.GitServer.BeginUpdateFromGit = fakeGitUpdateFromGit()
	fakes.FakeServer.ServerFactory.Core.GitServer.Disconnect = fakeGitDisconnect()

	testHelperGitProviderDetails := map[string]any{
		"git_provider_type": string(*gitProviderDetails.GitProviderType),
		"organization_name": "TestOrganization",
		"project_name":      "TestProject",
		"repository_name":   *gitProviderDetails.RepositoryName,
		"branch_name":       *gitProviderDetails.BranchName,
		"directory_name":    *gitProviderDetails.DirectoryName,
	}

	testCaseInvalidGitProviderType := testhelp.CopyMap(testHelperGitProviderDetails)
	testCaseInvalidGitProviderType["git_provider_type"] = "test1"

	testCaseInvalidDirectoryName := testhelp.CopyMap(testHelperGitProviderDetails)
	testCaseInvalidDirectoryName["directory_name"] = "test2"

	testCaseInvalidOwnerName := testhelp.CopyMap(testHelperGitProviderDetails)
	testCaseInvalidOwnerName["owner_name"] = "test3"

	testCaseMissingBranchName := testhelp.CopyMap(testHelperGitProviderDetails)
	delete(testCaseMissingBranchName, "branch_name")

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceWorkspaceGitFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testResourceWorkspaceGitFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceGitHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - no required git_provider_details
		{
			ResourceName: testResourceWorkspaceGitFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceGitHeader,
				map[string]any{
					"workspace_id":            "00000000-0000-0000-0000-000000000000",
					"initialization_strategy": "PreferWorkspace",
				},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - no required initialization_strategy
		{
			ResourceName: testResourceWorkspaceGitFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceGitHeader,
				map[string]any{
					"workspace_id":         "00000000-0000-0000-0000-000000000000",
					"git_provider_details": testHelperGitProviderDetails,
				},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - invalid initialization_strategy
		{
			ResourceName: testResourceWorkspaceGitFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceGitHeader,
				map[string]any{
					"workspace_id":            "00000000-0000-0000-0000-000000000000",
					"initialization_strategy": "test",
					"git_provider_details":    testHelperGitProviderDetails,
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
		},
		// error - invalid git_provider_type
		{
			ResourceName: testResourceWorkspaceGitFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceGitHeader,
				map[string]any{
					"workspace_id":            "00000000-0000-0000-0000-000000000000",
					"initialization_strategy": "PreferWorkspace",
					"git_provider_details":    testCaseInvalidGitProviderType,
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
		},
		// error - invalid directory_name
		{
			ResourceName: testResourceWorkspaceGitFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceGitHeader,
				map[string]any{
					"workspace_id":            "00000000-0000-0000-0000-000000000000",
					"initialization_strategy": "PreferWorkspace",
					"git_provider_details":    testCaseInvalidDirectoryName,
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
		},
		// error - invalid owner_name
		{
			ResourceName: testResourceWorkspaceGitFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceGitHeader,
				map[string]any{
					"workspace_id":            "00000000-0000-0000-0000-000000000000",
					"initialization_strategy": "PreferWorkspace",
					"git_provider_details":    testCaseInvalidOwnerName,
				},
			),
			ExpectError: regexp.MustCompile("Invalid configuration for attribute git_provider_details.owner_name"),
		},
		// error - missing branch_name
		{
			ResourceName: testResourceWorkspaceGitFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceGitHeader,
				map[string]any{
					"workspace_id":            "00000000-0000-0000-0000-000000000000",
					"initialization_strategy": "PreferWorkspace",
					"git_provider_details":    testCaseMissingBranchName,
				},
			),
			ExpectError: regexp.MustCompile(`Incorrect attribute value type`),
		},
		// error - invalid git_credentials
		{
			ResourceName: testResourceWorkspaceGitFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceGitHeader,
				map[string]any{
					"workspace_id":            "00000000-0000-0000-0000-000000000000",
					"initialization_strategy": "PreferWorkspace",
					"git_provider_details":    testHelperGitProviderDetails,
					"git_credentials": map[string]any{
						"connection_id": "00000000-0000-0000-0000-000000000000",
					},
				},
			),
			ExpectError: regexp.MustCompile("Invalid configuration for attribute git_credentials"),
		},
		// ok - PreferWorkspace
		{
			ResourceName: testResourceWorkspaceGitFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceGitHeader,
				map[string]any{
					"workspace_id":            "00000000-0000-0000-0000-000000000000",
					"initialization_strategy": "PreferWorkspace",
					"git_provider_details":    testHelperGitProviderDetails,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceWorkspaceGitFQN, "git_connection_state", (*string)(gitConnection.GitConnectionState)),
			),
		},
		// ok - PreferRemote
		{
			ResourceName: testResourceWorkspaceGitFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceGitHeader,
				map[string]any{
					"workspace_id":            "00000000-0000-0000-0000-000000000000",
					"initialization_strategy": "PreferRemote",
					"git_provider_details":    testHelperGitProviderDetails,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceWorkspaceGitFQN, "git_connection_state", (*string)(gitConnection.GitConnectionState)),
			),
		},
	}))
}

func TestAcc_WorkspaceGitResource_AzDO(t *testing.T) {
	if testhelp.ShouldSkipTest(t) {
		t.Skip("No SPN support")
	}

	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	doPlatform := testhelp.WellKnown()["AzDO"].(map[string]any)
	azdoOrganization := doPlatform["organizationName"].(string)
	azdoProject := doPlatform["projectName"].(string)
	azdoRepository := doPlatform["repositoryName"].(string)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceWorkspaceGitFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceWorkspaceGitFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceWorkspaceGitHeader,
					map[string]any{
						"workspace_id":            testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"initialization_strategy": "PreferWorkspace",
						"git_provider_details": map[string]any{
							"git_provider_type": "AzureDevOps",
							"organization_name": azdoOrganization,
							"project_name":      azdoProject,
							"repository_name":   azdoRepository,
							"branch_name":       "main",
							"directory_name":    "/",
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testResourceWorkspaceGitFQN, "git_sync_details.head"),
				resource.TestCheckResourceAttr(testResourceWorkspaceGitFQN, "git_connection_state", string(fabcore.GitConnectionStateConnectedAndInitialized)),
			),
		},
	},
	))
}

// func TestUnit_WorkspaceGitResource_GitHub(t *testing.T) {
// 	gitConnection := NewRandomGitConnection(fabcore.GitProviderTypeGitHub)
// 	gitCredentials := NewRandomGitCredentialsResponse(fabcore.GitCredentialsSourceConfiguredConnection)
// 	gitProviderDetails := gitConnection.GitProviderDetails.GetGitProviderDetails()
// 	gitInit := NewRandomGitInitializeGitConnection()

// 	fakes.FakeServer.ServerFactory.Core.GitServer.GetConnection = fakeGitGetConnection(gitConnection)
// 	fakes.FakeServer.ServerFactory.Core.GitServer.GetMyGitCredentials = fakeGitGetMyGitCredentials(gitCredentials)
// 	fakes.FakeServer.ServerFactory.Core.GitServer.Connect = fakeGitConnect()
// 	fakes.FakeServer.ServerFactory.Core.GitServer.BeginInitializeConnection = fakeGitInitializeGitConnection(gitInit)
// 	fakes.FakeServer.ServerFactory.Core.GitServer.BeginCommitToGit = fakeGitCommitToGit()
// 	fakes.FakeServer.ServerFactory.Core.GitServer.BeginUpdateFromGit = fakeGitUpdateFromGit()
// 	fakes.FakeServer.ServerFactory.Core.GitServer.Disconnect = fakeGitDisconnect()

// 	gitCredentialsResponse := gitCredentials.GitCredentialsConfigurationResponseClassification.(*fabcore.ConfiguredConnectionGitCredentialsResponse)

// 	testHelperGitProviderDetails := map[string]any{
// 		"git_provider_type": string(*gitProviderDetails.GitProviderType),
// 		"owner_name":        "TestOwner",
// 		"repository_name":   *gitProviderDetails.RepositoryName,
// 		"branch_name":       *gitProviderDetails.BranchName,
// 		"directory_name":    *gitProviderDetails.DirectoryName,
// 	}

// 	testCaseInvalidGitProviderType := testhelp.CopyMap(testHelperGitProviderDetails)
// 	testCaseInvalidGitProviderType["git_provider_type"] = "test1"

// 	testCaseInvalidDirectoryName := testhelp.CopyMap(testHelperGitProviderDetails)
// 	testCaseInvalidDirectoryName["directory_name"] = "test2"

// 	testCaseInvalidOrganizationName := testhelp.CopyMap(testHelperGitProviderDetails)
// 	testCaseInvalidOrganizationName["organization_name"] = "test3"

// 	testCaseMissingBranchName := testhelp.CopyMap(testHelperGitProviderDetails)
// 	delete(testCaseMissingBranchName, "branch_name")

// 	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceWorkspaceGitFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
// 		// error - no attributes
// 		{
// 			ResourceName: testResourceWorkspaceGitFQN,
// 			Config: at.CompileConfig(
// 				testResourceWorkspaceGitHeader,
// 				map[string]any{},
// 			),
// 			ExpectError: regexp.MustCompile(`Missing required argument`),
// 		},
// 		// error - no required git_provider_details
// 		{
// 			ResourceName: testResourceWorkspaceGitFQN,
// 			Config: at.CompileConfig(
// 				testResourceWorkspaceGitHeader,
// 				map[string]any{
// 					"workspace_id":            "00000000-0000-0000-0000-000000000000",
// 					"initialization_strategy": "PreferWorkspace",
// 				},
// 			),
// 			ExpectError: regexp.MustCompile(`Missing required argument`),
// 		},
// 		// error - no required initialization_strategy
// 		{
// 			ResourceName: testResourceWorkspaceGitFQN,
// 			Config: at.CompileConfig(
// 				testResourceWorkspaceGitHeader,
// 				map[string]any{
// 					"workspace_id":         "00000000-0000-0000-0000-000000000000",
// 					"git_provider_details": testHelperGitProviderDetails,
// 				},
// 			),
// 			ExpectError: regexp.MustCompile(`Missing required argument`),
// 		},
// 		// error - invalid initialization_strategy
// 		{
// 			ResourceName: testResourceWorkspaceGitFQN,
// 			Config: at.CompileConfig(
// 				testResourceWorkspaceGitHeader,
// 				map[string]any{
// 					"workspace_id":            "00000000-0000-0000-0000-000000000000",
// 					"initialization_strategy": "test",
// 					"git_provider_details":    testHelperGitProviderDetails,
// 				},
// 			),
// 			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
// 		},
// 		// error - invalid git_provider_type
// 		{
// 			ResourceName: testResourceWorkspaceGitFQN,
// 			Config: at.CompileConfig(
// 				testResourceWorkspaceGitHeader,
// 				map[string]any{
// 					"workspace_id":            "00000000-0000-0000-0000-000000000000",
// 					"initialization_strategy": "PreferWorkspace",
// 					"git_provider_details":    testCaseInvalidGitProviderType,
// 				},
// 			),
// 			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
// 		},
// 		// error - invalid directory_name
// 		{
// 			ResourceName: testResourceWorkspaceGitFQN,
// 			Config: at.CompileConfig(
// 				testResourceWorkspaceGitHeader,
// 				map[string]any{
// 					"workspace_id":            "00000000-0000-0000-0000-000000000000",
// 					"initialization_strategy": "PreferWorkspace",
// 					"git_provider_details":    testCaseInvalidDirectoryName,
// 				},
// 			),
// 			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
// 		},
// 		// error - invalid owner_name
// 		{
// 			ResourceName: testResourceWorkspaceGitFQN,
// 			Config: at.CompileConfig(
// 				testResourceWorkspaceGitHeader,
// 				map[string]any{
// 					"workspace_id":            "00000000-0000-0000-0000-000000000000",
// 					"initialization_strategy": "PreferWorkspace",
// 					"git_provider_details":    testCaseInvalidOrganizationName,
// 				},
// 			),
// 			ExpectError: regexp.MustCompile("Invalid configuration for attribute git_provider_details.organization_name"),
// 		},
// 		// error - missing branch_name
// 		{
// 			ResourceName: testResourceWorkspaceGitFQN,
// 			Config: at.CompileConfig(
// 				testResourceWorkspaceGitHeader,
// 				map[string]any{
// 					"workspace_id":            "00000000-0000-0000-0000-000000000000",
// 					"initialization_strategy": "PreferWorkspace",
// 					"git_provider_details":    testCaseMissingBranchName,
// 				},
// 			),
// 			ExpectError: regexp.MustCompile(`Incorrect attribute value type`),
// 		},
// 		// ok - PreferWorkspace
// 		{
// 			ResourceName: testResourceWorkspaceGitFQN,
// 			Config: at.CompileConfig(
// 				testResourceWorkspaceGitHeader,
// 				map[string]any{
// 					"workspace_id":            "00000000-0000-0000-0000-000000000000",
// 					"initialization_strategy": "PreferWorkspace",
// 					"git_provider_details":    testHelperGitProviderDetails,
// 					"git_credentials": map[string]any{
// 						"connection_id": *gitCredentialsResponse.ConnectionID,
// 					},
// 				},
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttrPtr(testResourceWorkspaceGitFQN, "git_connection_state", (*string)(gitConnection.GitConnectionState)),
// 			),
// 		},
// 		// ok - PreferRemote
// 		{
// 			ResourceName: testResourceWorkspaceGitFQN,
// 			Config: at.CompileConfig(
// 				testResourceWorkspaceGitHeader,
// 				map[string]any{
// 					"workspace_id":            "00000000-0000-0000-0000-000000000000",
// 					"initialization_strategy": "PreferRemote",
// 					"git_provider_details":    testHelperGitProviderDetails,
// 					"git_credentials": map[string]any{
// 						"connection_id": *gitCredentialsResponse.ConnectionID,
// 					},
// 				},
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttrPtr(testResourceWorkspaceGitFQN, "git_connection_state", (*string)(gitConnection.GitConnectionState)),
// 			),
// 		},
// 	}))
// }

func TestAcc_WorkspaceGitResource_GitHub(t *testing.T) {
	if testhelp.ShouldSkipTest(t) {
		t.Skip("No SPN support")
	}

	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	doPlatform := testhelp.WellKnown()["GitHub"].(map[string]any)
	ghOwner := doPlatform["ownerName"].(string)
	ghRepository := doPlatform["repositoryName"].(string)
	ghConnectionID := doPlatform["connectionId"].(string)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceWorkspaceGitFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceWorkspaceGitFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceWorkspaceGitHeader,
					map[string]any{
						"workspace_id":            testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"initialization_strategy": "PreferWorkspace",
						"git_provider_details": map[string]any{
							"git_provider_type": "GitHub",
							"owner_name":        ghOwner,
							"repository_name":   ghRepository,
							"branch_name":       "main",
							"directory_name":    "/",
						},
						"git_credentials": map[string]any{
							"connection_id": ghConnectionID,
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testResourceWorkspaceGitFQN, "git_sync_details.head"),
				resource.TestCheckResourceAttr(testResourceWorkspaceGitFQN, "git_connection_state", string(fabcore.GitConnectionStateConnectedAndInitialized)),
			),
		},
	},
	))
}
