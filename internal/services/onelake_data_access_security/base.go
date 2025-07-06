// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelake_data_access_security

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "OneLake Data Access Security",
	Type:           "onelake_data_access_security",
	DocsURL:        "https://learn.microsoft.com/power-bi/consumer/end-user-dashboards",
	IsPreview:      true,
	IsSPNSupported: true,
}
