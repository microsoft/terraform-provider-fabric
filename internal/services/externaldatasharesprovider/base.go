package externaldatasharesprovider

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Names:          "External Data Shares Provider",
	Types:          "external_data_shares_provider",
	DocsURL:        "https://learn.microsoft.com/rest/api/fabric/admin/external-data-shares-provider/list-external-data-shares",
	IsPreview:      false,
	IsSPNSupported: true,
}
