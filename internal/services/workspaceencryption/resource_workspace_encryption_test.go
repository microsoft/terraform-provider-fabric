// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceencryption_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_WorkspaceEncryptionResource_Attributes(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	fakeServer := fakes.NewFakeServer()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - missing key_identifier
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - invalid workspace_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   "invalid-uuid",
					"key_identifier": NewRandomKeyIdentifier(),
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected_attr
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":    workspaceID,
					"key_identifier":  NewRandomKeyIdentifier(),
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - key_identifier is not a Key Vault URI
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"key_identifier": "not-a-key-identifier",
				},
			),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
		},
		// error - key_identifier is versioned
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"key_identifier": "https://example.vault.azure.net/keys/example/example-version",
				},
			),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
		},
	}))
}

func TestUnit_WorkspaceEncryptionResource_CRUD(t *testing.T) {
	entity := fabcore.WorkspaceEncryptionDetail{}
	workspaceID := testhelp.RandomUUID()
	keyIdentifier := NewRandomKeyIdentifier()
	keyIdentifierUpdated := NewRandomKeyIdentifier()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.GetWorkspaceEncryption = fakeGetWorkspaceEncryption(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.AssignWorkspaceEncryption = fakeAssignWorkspaceEncryption(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.ResetWorkspaceEncryption = fakeResetWorkspaceEncryption(&entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		// create and read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"key_identifier": keyIdentifier,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "key_identifier", keyIdentifier),
				resource.TestCheckResourceAttr(testResourceItemFQN, "encryption_status", string(fabcore.WorkspaceEncryptionStatusActive)),
			),
		},
		// update and read - rotate the key
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"key_identifier": keyIdentifierUpdated,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "key_identifier", keyIdentifierUpdated),
				resource.TestCheckResourceAttr(testResourceItemFQN, "encryption_status", string(fabcore.WorkspaceEncryptionStatusActive)),
			),
		},
	}))
}

func TestUnit_WorkspaceEncryptionResource_Failed(t *testing.T) {
	entity := fabcore.WorkspaceEncryptionDetail{}
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.GetWorkspaceEncryption = fakeGetWorkspaceEncryption(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.AssignWorkspaceEncryption = fakeAssignWorkspaceEncryptionWithStatus(&entity, fabcore.WorkspaceEncryptionStatusFailed)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		// error - encryption settles on Failed
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"key_identifier": NewRandomKeyIdentifier(),
				},
			),
			ExpectError: regexp.MustCompile(`Workspace Encryption failed for Workspace ID`),
		},
	}))
}

func TestUnit_WorkspaceEncryptionResource_Timeout(t *testing.T) {
	entity := fabcore.WorkspaceEncryptionDetail{}
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.GetWorkspaceEncryption = fakeGetWorkspaceEncryption(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.AssignWorkspaceEncryption = fakeAssignWorkspaceEncryptionWithStatus(&entity, fabcore.WorkspaceEncryptionStatusEnableInProgress)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		// error - encryption never leaves EnableInProgress
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"key_identifier": NewRandomKeyIdentifier(),
					"timeouts": map[string]any{
						"create": "10s",
					},
				},
			),
			ExpectError: regexp.MustCompile(`Timeout waiting for Workspace Encryption`),
		},
	}))
}

func TestUnit_WorkspaceEncryptionResource_EventuallyActive(t *testing.T) {
	entity := fabcore.WorkspaceEncryptionDetail{}
	workspaceID := testhelp.RandomUUID()
	keyIdentifier := NewRandomKeyIdentifier()

	fakeServer := fakes.NewFakeServer()
	// The first poll reports EnableInProgress, so the loop must wait one interval before the status settles.
	fakeServer.ServerFactory.Core.WorkspacesServer.GetWorkspaceEncryption = fakeGetWorkspaceEncryptionInProgress(&entity, 1)
	fakeServer.ServerFactory.Core.WorkspacesServer.AssignWorkspaceEncryption = fakeAssignWorkspaceEncryption(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.ResetWorkspaceEncryption = fakeResetWorkspaceEncryption(&entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		// create and read - the in-progress status must be polled until it becomes Active
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"key_identifier": keyIdentifier,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "key_identifier", keyIdentifier),
				resource.TestCheckResourceAttr(testResourceItemFQN, "encryption_status", string(fabcore.WorkspaceEncryptionStatusActive)),
			),
		},
	}))
}

func TestUnit_WorkspaceEncryptionResource_Disabled(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	keyIdentifier := NewRandomKeyIdentifier()

	entity := fabcore.WorkspaceEncryptionDetail{}

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.GetWorkspaceEncryption = fakeGetWorkspaceEncryption(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.AssignWorkspaceEncryption = fakeAssignWorkspaceEncryption(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.ResetWorkspaceEncryption = fakeResetWorkspaceEncryption(&entity)

	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id":   workspaceID,
			"key_identifier": keyIdentifier,
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		// error - import a workspace without a customer-managed key
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: workspaceID,
			ImportState:   true,
			ExpectError:   regexp.MustCompile(`is not enabled for Workspace ID`),
		},
		// create and read
		{
			ResourceName: testResourceItemFQN,
			Config:       testCase,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "encryption_status", string(fabcore.WorkspaceEncryptionStatusActive)),
			),
		},
		// the key was removed outside of Terraform, so the resource must be dropped from state
		{
			PreConfig: func() {
				entity.EncryptionDetail.EncryptionStatus = azto.Ptr(fabcore.WorkspaceEncryptionStatusDisabled)
			},
			ResourceName:       testResourceItemFQN,
			Config:             testCase,
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	}))
}

func TestUnit_WorkspaceEncryptionResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := NewRandomWorkspaceEncryptionDetail()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.GetWorkspaceEncryption = fakeGetWorkspaceEncryption(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.AssignWorkspaceEncryption = fakeAssignWorkspaceEncryption(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.ResetWorkspaceEncryption = fakeResetWorkspaceEncryption(&entity)

	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id":   workspaceID,
			"key_identifier": *entity.EncryptionDetail.KeyIdentifier,
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: workspaceID,
			ImportState:   true,
			ImportStateCheck: func(is []*terraform.InstanceState) error {
				if len(is) != 1 {
					return errors.New("expected one instance state")
				}

				if got := is[0].Attributes["workspace_id"]; got != workspaceID {
					return fmt.Errorf("%s: unexpected workspace_id — got %q, want %q", testResourceItemFQN, got, workspaceID)
				}

				if got, want := is[0].Attributes["key_identifier"], *entity.EncryptionDetail.KeyIdentifier; got != want {
					return fmt.Errorf("%s: unexpected key_identifier — got %q, want %q", testResourceItemFQN, got, want)
				}

				return nil
			},
		},
	}))
}

func TestAcc_WorkspaceEncryptionResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceCMK"].(map[string]any)
	workspaceID := workspace["id"].(string)

	keyVault := testhelp.WellKnown()["KeyVault"].(map[string]any)
	keyIdentifier := keyVault["keyIdentifier"].(string)
	keyIdentifierUpdated := keyVault["keyIdentifierAlt"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// create and read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"key_identifier": keyIdentifier,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "key_identifier", keyIdentifier),
				resource.TestCheckResourceAttr(testResourceItemFQN, "encryption_status", string(fabcore.WorkspaceEncryptionStatusActive)),
			),
		},
		// update and read - rotate the key
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"key_identifier": keyIdentifierUpdated,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "key_identifier", keyIdentifierUpdated),
				resource.TestCheckResourceAttr(testResourceItemFQN, "encryption_status", string(fabcore.WorkspaceEncryptionStatusActive)),
			),
		},
	}))
}
