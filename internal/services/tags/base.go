package tags

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Tag",
	Type:           "tag",
	Names:          "Tags",
	Types:          "tags",
	DocsURL:        "https://learn.microsoft.com/fabric/governance/tags-overview",
	IsPreview:      true,
	IsSPNSupported: true,
}
