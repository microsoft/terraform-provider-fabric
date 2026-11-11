// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package monitoringdatacollection

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type resourceMonitoringDataCollectionModel struct {
	ID           customtypes.UUID `tfsdk:"id"`
	WorkspaceID  customtypes.UUID `tfsdk:"workspace_id"`
	MonitoringID customtypes.UUID `tfsdk:"monitoring_id"`
	Enabled      types.Bool       `tfsdk:"enabled"`
	Timeouts     timeouts.Value   `tfsdk:"timeouts"`
}
