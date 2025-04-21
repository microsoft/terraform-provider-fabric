// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package deploymentpipeline

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Deployment Pipeline",
	Type:           "deployment_pipeline",
	Names:          "Deployment Pipelines",
	Types:          "deployment_pipelines",
	DocsURL:        "https://learn.microsoft.com/fabric/data-factory/what-is-copy-job",
	IsPreview:      false,
	IsSPNSupported: true,
}
