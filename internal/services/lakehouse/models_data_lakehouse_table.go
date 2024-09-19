// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fablakehouse "github.com/microsoft/fabric-sdk-go/fabric/lakehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceLakehouseTableModel struct {
	LakehouseID customtypes.UUID `tfsdk:"lakehouse_id"`
	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
	lakehouseTableModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type lakehouseTableModel struct {
	Name     types.String `tfsdk:"name"`
	Location types.String `tfsdk:"location"`
	Type     types.String `tfsdk:"type"`
	Format   types.String `tfsdk:"format"`
}

func (to *lakehouseTableModel) set(from fablakehouse.Table) {
	to.Name = types.StringPointerValue(from.Name)
	to.Location = types.StringPointerValue(from.Location)
	to.Type = types.StringPointerValue((*string)(from.Type))
	to.Format = types.StringPointerValue(from.Format)
}
