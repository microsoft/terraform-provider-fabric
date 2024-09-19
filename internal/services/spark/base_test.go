// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package spark_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"

	"github.com/microsoft/terraform-provider-fabric/internal/services/spark"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

const (
	sparkCustomPoolTFName           = spark.SparkCustomPoolTFName
	sparkWorkspaceSettingsTFName    = spark.SparkWorkspaceSettingsTFName
	sparkEnvironmentSettingsTFName  = spark.SparkEnvironmentSettingsTFName
	sparkEnvironmentLibrariesTFName = spark.SparkEnvironmentLibrariesTFName
)

func getSparkCustomPoolResourceAttr(t *testing.T, workspaceID, name string) map[string]any {
	t.Helper()

	return map[string]any{
		"workspace_id": workspaceID,
		"name":         name,
		"type":         "Workspace",
		"node_family":  "MemoryOptimized",
		"node_size":    "Small",
		"auto_scale": map[string]any{
			"enabled":        true,
			"min_node_count": 1,
			"max_node_count": 3,
		},
		"dynamic_executor_allocation": map[string]any{
			"enabled":       true,
			"min_executors": 1,
			"max_executors": 2,
		},
	}
}

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
