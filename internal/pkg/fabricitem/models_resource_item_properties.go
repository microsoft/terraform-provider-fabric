// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

type ResourceFabricItemPropertiesModel[Ttfprop, Titemprop any] struct {
	FabricItemPropertiesModel[Ttfprop, Titemprop]
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type requestCreateFabricItemProperties[Ttfprop, Titemprop any] struct {
	fabcore.CreateItemRequest
}

func (to *requestCreateFabricItemProperties[Ttfprop, Titemprop]) set(from ResourceFabricItemPropertiesModel[Ttfprop, Titemprop], itemType fabcore.ItemType) { //revive:disable-line:confusing-naming
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
	to.Type = azto.Ptr(itemType)
}

type requestUpdateFabricItemProperties[Ttfprop, Titemprop any] struct {
	fabcore.UpdateItemRequest
}

func (to *requestUpdateFabricItemProperties[Ttfprop, Titemprop]) set(from ResourceFabricItemPropertiesModel[Ttfprop, Titemprop]) { //revive:disable-line:confusing-naming
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
}
