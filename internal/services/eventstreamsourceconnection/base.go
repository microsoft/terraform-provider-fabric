// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventstreamsourceconnection

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Eventstream Source Connection",
	Type:           "eventstream_source_connection",
	DocsURL:        "https://learn.microsoft.com/fabric/real-time-intelligence/event-streams/overview",
	IsPreview:      true,
	IsSPNSupported: true,
}
