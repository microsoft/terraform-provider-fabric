// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Connection",
	Type:           "connection",
	Names:          "Connections",
	Types:          "connections",
	DocsURL:        "TODO",
	IsPreview:      true,
	IsSPNSupported: true,
}
