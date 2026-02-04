// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package externaldatasharesprovider_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"

	"github.com/microsoft/terraform-provider-fabric/internal/services/externaldatasharesprovider"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var itemTypeInfo = externaldatasharesprovider.ItemTypeInfo

func lakehouseResource(t *testing.T, workspaceID string) (resourceHCL, resourceFQN string) {
	t.Helper()

	resourceHCL = at.CompileConfig(
		at.ResourceHeader(testhelp.TypeName("fabric", "lakehouse"), "test"),
		map[string]any{
			"display_name": testhelp.RandomName(),
			"workspace_id": workspaceID,
		},
	)

	resourceFQN = testhelp.ResourceFQN("fabric", "lakehouse", "test")

	return resourceHCL, resourceFQN
}
