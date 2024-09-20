// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package capacity

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseCapacityModel struct {
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
	SKU         types.String     `tfsdk:"sku"`
	Region      types.String     `tfsdk:"region"`
	State       types.String     `tfsdk:"state"`
}

func (to *baseCapacityModel) set(from fabcore.Capacity) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.SKU = types.StringPointerValue(from.SKU)
	to.Region = types.StringPointerValue(from.Region)
	to.State = types.StringPointerValue((*string)(from.State))
}