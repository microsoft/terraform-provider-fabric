// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelakeshortcut

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "OneLake Shortcut",
	Type:           "onelake_shortcut",
	Names:          "OneLake Shortcuts",
	Types:          "onelake_shortcuts",
	DocsURL:        "https://learn.microsoft.com/fabric/onelake/onelake-shortcuts",
	IsPreview:      false,
	IsSPNSupported: true,
}
