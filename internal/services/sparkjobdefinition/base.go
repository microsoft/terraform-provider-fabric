// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sparkjobdefinition

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName                  = "Spark Job Definition"
	ItemTFName                = "spark_job_definition"
	ItemsName                 = "Spark Job Definitions"
	ItemsTFName               = "spark_job_definitions"
	ItemType                  = fabcore.ItemTypeSparkJobDefinition
	ItemDocsSPNSupport        = common.DocsSPNSupported
	ItemDocsURL               = "https://learn.microsoft.com/fabric/data-engineering/spark-job-definition"
	ItemFormatTypeDefault     = "SparkJobDefinitionV1"
	ItemDefinitionEmpty       = `{"executableFile":null,"defaultLakehouseArtifactId":null,"mainClass":null,"additionalLakehouseIds":[],"retryPolicy":null,"commandLineArguments":null,"additionalLibraryUris":null,"language":null,"environmentArtifactId":null}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/spark-job-definition"
)

var (
	ItemFormatTypes     = []string{"SparkJobDefinitionV1"}      //nolint:gochecknoglobals
	ItemDefinitionPaths = []string{"SparkJobDefinitionV1.json"} //nolint:gochecknoglobals
)
