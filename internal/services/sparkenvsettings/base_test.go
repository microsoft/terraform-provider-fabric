// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sparkenvsettings_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"

	"github.com/microsoft/terraform-provider-fabric/internal/services/sparkenvsettings"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var itemTypeInfo = sparkenvsettings.ItemTypeInfo

func environmentResource(t *testing.T, workspaceID string) (resourceHCL, resourceFQN string) {
	t.Helper()

	resourceHCL = at.CompileConfig(
		at.ResourceHeader(testhelp.TypeName("fabric", "environment"), "test"),
		map[string]any{
			"display_name": testhelp.RandomName(),
			"workspace_id": workspaceID,
		},
	)

	resourceFQN = testhelp.ResourceFQN("fabric", "environment", "test")

	return resourceHCL, resourceFQN
}
