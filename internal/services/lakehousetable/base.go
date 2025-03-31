// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehousetable

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Lakehouse Table",
	Type:           "lakehouse_table",
	Names:          "Lakehouse Tables",
	Types:          "lakehouse_tables",
	DocsURL:        "https://learn.microsoft.com/fabric/data-engineering/lakehouse-and-delta-tables",
	IsPreview:      true,
	IsSPNSupported: true,
}
