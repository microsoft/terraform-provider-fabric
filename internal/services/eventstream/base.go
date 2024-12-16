// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventstream

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName           = "Eventstream"
	ItemTFName         = "eventstream"
	ItemsName          = "Eventstreams"
	ItemsTFName        = "eventstreams"
	ItemType           = fabcore.ItemTypeEventstream
	ItemDocsSPNSupport = common.DocsSPNSupported
	ItemDocsURL        = "https://learn.microsoft.com/fabric/real-time-intelligence/event-streams/overview"
)
