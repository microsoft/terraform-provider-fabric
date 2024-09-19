// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehouse

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	fabwarehouse "github.com/microsoft/fabric-sdk-go/fabric/warehouse"
)

type resourceWarehouseModel struct {
	baseWarehouseModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type requestCreateWarehouse struct {
	fabwarehouse.CreateWarehouseRequest
}

func (to *requestCreateWarehouse) set(from resourceWarehouseModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
}

type requestUpdateWarehouse struct {
	fabwarehouse.UpdateWarehouseRequest
}

func (to *requestUpdateWarehouse) set(from resourceWarehouseModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
}
