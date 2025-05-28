// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package shortcut

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Shortcut",
	Type:           "shortcut",
	Names:          "Shortcuts",
	Types:          "shortcuts",
	DocsURL:        "https://learn.microsoft.com/fabric/onelake/onelake-shortcuts",
	IsPreview:      false,
	IsSPNSupported: true,
}
