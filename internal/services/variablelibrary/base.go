// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package variablelibrary

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const FabricItemType = fabcore.ItemTypeVariableLibrary

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "VariableLibrary",
	Type:           "variablelibrary",
	Names:          "VariableLibraries",
	Types:          "VariableLibraries",
	DocsURL:        "<TBD>",
	IsPreview:      false,
	IsSPNSupported: true,
}
