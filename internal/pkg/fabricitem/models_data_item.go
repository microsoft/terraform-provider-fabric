// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type dataSourceFabricItemModel struct {
	fabricItemModel

	SensitivityLabel supertypes.SingleNestedObjectValueOf[sensitivityLabelModel] `tfsdk:"sensitivity_label"`
	Timeouts         timeouts.Value                                              `tfsdk:"timeouts"`
}

func (to *dataSourceFabricItemModel) set(ctx context.Context, from fabcore.Item) diag.Diagnostics {
	to.fabricItemModel.set(from)

	sl, diags := newSensitivityLabelFromAPI(ctx, from.SensitivityLabel)
	if diags.HasError() {
		return diags
	}

	to.SensitivityLabel = sl

	return nil
}
