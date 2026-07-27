// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package itemjob

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Item Job",
	Type:           "item_job",
	Names:          "Item Jobs",
	Types:          "item_jobs",
	DocsURL:        "https://learn.microsoft.com/fabric/fundamentals/job-scheduler",
	IsPreview:      false,
	IsSPNSupported: true,
}
