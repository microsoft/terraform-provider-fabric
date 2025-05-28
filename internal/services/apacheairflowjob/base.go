// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package apacheairflowjob

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeApacheAirflowJob
	ItemDefinitionEmpty       = `{"properties":{}}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/fabric/data-factory/apache-airflow-jobs-concepts"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Apache Airflow Job",
	Type:           "apache_airflow_job",
	Names:          "Apache Airflow Jobs",
	Types:          "apache_airflow_jobs",
	DocsURL:        "https://learn.microsoft.com/fabric/data-factory/apache-airflow-jobs-concepts",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"apacheairflowjob-content.json"},
	},
}
