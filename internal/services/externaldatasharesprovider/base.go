// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package externaldatasharesprovider

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "External Data Share",
	Type:           "external_data_share",
	Names:          "External Data Shares",
	Types:          "external_data_shares",
	DocsURL:        "https://learn.microsoft.com/rest/api/fabric/admin/external-data-shares-provider/list-external-data-shares",
	IsPreview:      false,
	IsSPNSupported: true,
}
