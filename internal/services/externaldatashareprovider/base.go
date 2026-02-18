// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package externaldatashareprovider

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "External Data Share",
	Type:           "external_data_share",
	Names:          "External Data Shares",
	Types:          "external_data_shares",
	DocsURL:        "https://learn.microsoft.com/fabric/governance/external-data-sharing-overview",
	IsPreview:      false,
	IsSPNSupported: true,
}
