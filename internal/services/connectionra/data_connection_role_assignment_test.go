// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connectionra_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_ConnectionRoleAssignmentDataSource(t *testing.T) {
	connectionID := testhelp.RandomUUID()
	connectionRoleAssignmentID := testhelp.RandomUUID()
	entity := NewRandomConnectionRoleAssignment()
	fakes.FakeServer.ServerFactory.Core.ConnectionsServer.GetConnectionRoleAssignment = fakeConnectionRoleAssignment(entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"connection_id":                 connectionID,
					"connection_role_assignment_id": connectionRoleAssignmentID,
					"unexpected_attr":               "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"connection_id":                 connectionID,
					"connection_role_assignment_id": connectionRoleAssignmentID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "connection_id", connectionID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "connection_role_assignment_id", connectionRoleAssignmentID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "role", (*string)(entity.Role)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "principal.id", entity.Principal.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "principal.type", (*string)(entity.Principal.Type)),
			),
		},
	}))
}
