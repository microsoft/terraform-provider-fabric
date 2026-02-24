// Copyright Microsoft Corporation 2026
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
	DocsURL:        "https://learn.microsoft.com/fabric/cicd/deployment-pipelines/intro-to-deployment-pipelines",
	IsPreview:      true,
	IsSPNSupported: true,
}
