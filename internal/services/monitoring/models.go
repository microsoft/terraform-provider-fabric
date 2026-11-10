// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package monitoring

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabmonitoring "github.com/microsoft/fabric-sdk-go/fabric/monitoring"
)

type monitoringPropertiesModel struct {
	DataCollectionState types.String `tfsdk:"data_collection_state"`
}

func (to *monitoringPropertiesModel) set(_ context.Context, from fabmonitoring.Properties) diag.Diagnostics {
	to.DataCollectionState = types.StringPointerValue((*string)(from.DataCollectionState))

	return nil
}
