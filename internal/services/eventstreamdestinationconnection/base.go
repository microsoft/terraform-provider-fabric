// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package eventstreamdestinationconnection

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Eventstream Destination Connection",
	Type:           "eventstream_destination_connection",
	DocsURL:        "https://learn.microsoft.com/fabric/real-time-intelligence/event-streams/overview",
	IsPreview:      true,
	IsSPNSupported: true,
}
