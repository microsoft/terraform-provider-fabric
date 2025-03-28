// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sparkjobdefinition

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeSparkJobDefinition
	ItemDefinitionEmpty       = `{"executableFile":null,"defaultLakehouseArtifactId":null,"mainClass":null,"additionalLakehouseIds":[],"retryPolicy":null,"commandLineArguments":null,"additionalLibraryUris":null,"language":null,"environmentArtifactId":null}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/spark-job-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Spark Job Definition",
	Type:           "spark_job_definition",
	Names:          "Spark Job Definitions",
	Types:          "spark_job_definitions",
	DocsURL:        "https://learn.microsoft.com/fabric/data-engineering/spark-job-definition",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  "SparkJobDefinitionV1",
		API:   "SparkJobDefinitionV1",
		Paths: []string{"SparkJobDefinitionV1.json"},
	},
}
