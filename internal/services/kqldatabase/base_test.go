// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package kqldatabase_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"

	"github.com/microsoft/terraform-provider-fabric/internal/services/kqldatabase"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

const (
	itemTFName  = kqldatabase.ItemTFName
	itemsTFName = kqldatabase.ItemsTFName
	itemType    = kqldatabase.ItemType
)

func eventhouseResource(t *testing.T, workspaceID string) (resourceHCL, resourceFQN string) {
	t.Helper()

	resourceHCL = at.CompileConfig(
		at.ResourceHeader(testhelp.TypeName("fabric", "eventhouse"), "test"),
		map[string]any{
			"display_name": testhelp.RandomName(),
			"workspace_id": workspaceID,
		},
	)

	resourceFQN = testhelp.ResourceFQN("fabric", "eventhouse", "test")

	return resourceHCL, resourceFQN
}
