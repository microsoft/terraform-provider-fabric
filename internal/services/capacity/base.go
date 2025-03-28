// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package capacity

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Capacity",
	Type:           "capacity",
	Names:          "Capacities",
	Types:          "capacities",
	DocsURL:        "https://learn.microsoft.com/fabric/enterprise/licenses#capacity",
	IsPreview:      false,
	IsSPNSupported: true,
}
