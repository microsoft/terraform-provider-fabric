// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceencryption_test

import (
	"regexp"
	"testing"

	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_WorkspaceEncryptionDataSource(t *testing.T) {
	entity := NewRandomWorkspaceEncryptionDetail()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.GetWorkspaceEncryption = fakeGetWorkspaceEncryption(&entity)

	workspaceID := testhelp.RandomUUID()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
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
		// error - invalid workspace_id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": "invalid-uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
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
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "key_identifier", *entity.EncryptionDetail.KeyIdentifier),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "encryption_status", string(*entity.EncryptionDetail.EncryptionStatus)),
			),
		},
	}))
}

func TestUnit_WorkspaceEncryptionDataSource_Disabled(t *testing.T) {
	entity := fabcore.WorkspaceEncryptionDetail{
		EncryptionDetail: &fabcore.EncryptionDetail{
			EncryptionStatus: azto.Ptr(fabcore.WorkspaceEncryptionStatusDisabled),
		},
	}

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.GetWorkspaceEncryption = fakeGetWorkspaceEncryption(&entity)

	workspaceID := testhelp.RandomUUID()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakeServer.ServerFactory, nil, []resource.TestStep{
		// read - customer-managed key is not configured
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "encryption_status", string(fabcore.WorkspaceEncryptionStatusDisabled)),
				resource.TestCheckNoResourceAttr(testDataSourceItemFQN, "key_identifier"),
			),
		},
	}))
}

func TestAcc_WorkspaceEncryptionDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceCMK"].(map[string]any)
	workspaceID := workspace["id"].(string)

	// Not parallel: the resource acceptance test mutates the encryption state of the same Workspace.
	resource.Test(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read by workspace id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "encryption_status"),
			),
		},
	}))
}
