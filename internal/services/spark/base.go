// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package spark

import (
	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	SparkCustomPoolName           = "Spark Custom Pool"
	SparkCustomPoolTFName         = "spark_custom_pool"
	SparkCustomPoolDocsSPNSupport = common.DocsSPNSupported
	SparkCustomPoolDocsURL        = "https://learn.microsoft.com/fabric/data-engineering/create-custom-spark-pools"

	SparkWorkspaceSettingsName           = "Spark Workspace Settings"
	SparkWorkspaceSettingsTFName         = "spark_workspace_settings"
	SparkWorkspaceSettingsDocsSPNSupport = common.DocsSPNSupported
	SparkWorkspaceSettingsDocsURL        = "https://learn.microsoft.com/fabric/data-engineering/workspace-admin-settings"

	SparkEnvironmentPublicationStatusPublished = "Published"
	SparkEnvironmentSettingsName               = "Spark Environment Settings"
	SparkEnvironmentSettingsTFName             = "spark_environment_settings"
	SparkEnvironmentSettingsDocsSPNSupport     = common.DocsSPNSupported
	SparkEnvironmentSettingsDocsURL            = "https://learn.microsoft.com/fabric/data-engineering/environment-manage-compute"

	SparkEnvironmentLibrariesName           = "Spark Environment Libraries"
	SparkEnvironmentLibrariesTFName         = "spark_environment_libraries"
	SparkEnvironmentLibrariesDocsSPNSupport = common.DocsSPNSupported
	SparkEnvironmentLibrariesDocsURL        = "https://learn.microsoft.com/fabric/data-engineering/environment-manage-library"
)

var (
	SparkRuntimeVersionValues               = []string{"1.1", "1.2", "1.3"}                  //nolint:gochecknoglobals
	SparkEnvironmentPublicationStatusValues = []string{"Published", "Staging"}               //nolint:gochecknoglobals
	SparkEnvironmentDriverCoresValues       = []int32{4, 8, 16, 32, 64}                      //nolint:gochecknoglobals
	SparkEnvironmentDriverMemoryValues      = []string{"28g", "56g", "112g", "224g", "400g"} //nolint:gochecknoglobals
	SparkEnvironmentExecutorCoresValues     = []int32{4, 8, 16, 32, 64}                      //nolint:gochecknoglobals
	SparkEnvironmentExecutorMemoryValues    = []string{"28g", "56g", "112g", "224g", "400g"} //nolint:gochecknoglobals

)
