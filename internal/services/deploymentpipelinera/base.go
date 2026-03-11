// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package deploymentpipelinera

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Deployment Pipeline Role Assignment",
	Type:           "deployment_pipeline_role_assignment",
	Names:          "Deployment Pipeline Assignments",
	Types:          "deployment_pipeline_role_assignments",
	DocsURL:        "https://learn.microsoft.com/fabric/cicd/deployment-pipelines/intro-to-deployment-pipelines",
	IsPreview:      true,
	IsSPNSupported: true,
}
