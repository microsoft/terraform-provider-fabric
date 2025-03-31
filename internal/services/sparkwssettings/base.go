// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sparkwssettings

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var SparkRuntimeVersionValues = []string{"1.1", "1.2", "1.3"} //nolint:gochecknoglobals

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Spark Workspace Settings",
	Type:           "spark_workspace_settings",
	DocsURL:        "https://learn.microsoft.com/fabric/data-engineering/workspace-admin-settings",
	IsPreview:      true,
	IsSPNSupported: true,
}
