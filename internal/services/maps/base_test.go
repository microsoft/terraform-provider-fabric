// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package maps_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"

	"github.com/microsoft/terraform-provider-fabric/internal/services/maps"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var itemTypeInfo = maps.ItemTypeInfo

const fabricItemType = maps.FabricItemType

func lakehouseResource(t *testing.T, workspaceID string) (resourceHCL, resourceFQN string) {
	t.Helper()

	resourceName := testhelp.RandomName()
	resourceHCL = at.CompileConfig(
		at.ResourceHeader(testhelp.TypeName("fabric", "lakehouse"), resourceName),
		map[string]any{
			"display_name": testhelp.RandomName(),
			"workspace_id": workspaceID,
		},
	)

	resourceFQN = testhelp.ResourceFQN("fabric", "lakehouse", resourceName)

	return resourceHCL, resourceFQN
}
