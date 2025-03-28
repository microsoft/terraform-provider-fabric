// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package datamart

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Datamart",
	Type:           "datamart",
	Names:          "Datamarts",
	Types:          "datamarts",
	DocsURL:        "https://learn.microsoft.com/power-bi/transform-model/datamarts/datamarts-overview",
	IsPreview:      true,
	IsSPNSupported: false,
}

const FabricItemType = fabcore.ItemTypeDatamart
