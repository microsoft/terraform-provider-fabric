package jobscheduler

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Job Scheduler",
	Type:           "job_scheduler",
	Names:          "Job Schedulers",
	Types:          "job_schedulers",
	DocsURL:        "https://learn.microsoft.com/fabric/fundamentals/workspaces-folders",
	IsPreview:      true,
	IsSPNSupported: true,
}
