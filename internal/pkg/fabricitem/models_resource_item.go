// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

type resourceFabricItemModel struct {
	baseFabricItemModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type requestCreateFabricItem struct {
	fabcore.CreateItemRequest
}

func (to *requestCreateFabricItem) set(from resourceFabricItemModel, itemType fabcore.ItemType) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
	to.Type = azto.Ptr(itemType)
}

type requestUpdateFabricItem struct {
	fabcore.UpdateItemRequest
}

func (to *requestUpdateFabricItem) set(from resourceFabricItemModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
}
