package itemjobscheduler

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Item Job Scheduler",
	Type:           "item_job_scheduler",
	Names:          "Item Job Schedulers",
	Types:          "item_job_schedulers",
	DocsURL:        "https://learn.microsoft.com/rest/api/fabric/articles/",
	IsPreview:      false,
	IsSPNSupported: true,
}
