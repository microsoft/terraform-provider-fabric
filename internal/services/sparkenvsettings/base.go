// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sparkenvsettings

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var (
	SparkEnvironmentPublicationStatusValues = []string{"Published", "Staging"} //nolint:gochecknoglobals
	SparkEnvironmentDriverCoresValues       = []int32{4, 8, 16, 32, 64}        //nolint:gochecknoglobals
	SparkEnvironmentExecutorCoresValues     = []int32{4, 8, 16, 32, 64}        //nolint:gochecknoglobals
)

const SparkEnvironmentPublicationStatusPublished = "Published"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Spark Environment Settings",
	Type:           "spark_environment_settings",
	DocsURL:        "https://learn.microsoft.com/fabric/data-engineering/environment-manage-compute",
	IsPreview:      true,
	IsSPNSupported: true,
}
