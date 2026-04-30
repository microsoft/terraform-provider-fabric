// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package warehousesqlauditsetting

import (
	"context"

	timeoutsd "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsr "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabwarehouse "github.com/microsoft/fabric-sdk-go/fabric/warehouse"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseWarehouseSQLAuditSettingsModel struct {
	WorkspaceID           customtypes.UUID                    `tfsdk:"workspace_id"`
	WarehouseID           customtypes.UUID                    `tfsdk:"warehouse_id"`
	State                 types.String                        `tfsdk:"state"`
	RetentionDays         types.Int32                         `tfsdk:"retention_days"`
	AuditActionsAndGroups supertypes.SetValueOf[types.String] `tfsdk:"audit_actions_and_groups"`
}

func (to *baseWarehouseSQLAuditSettingsModel) set(ctx context.Context, from fabwarehouse.SQLAuditSettings) {
	to.State = types.StringPointerValue((*string)(from.State))
	to.RetentionDays = types.Int32PointerValue(from.RetentionDays)
	to.AuditActionsAndGroups = supertypes.NewSetValueOfNull[types.String](ctx)

	elements := make([]types.String, len(from.AuditActionsAndGroups))
	for i, v := range from.AuditActionsAndGroups {
		elements[i] = types.StringValue(v)
	}

	to.AuditActionsAndGroups = supertypes.NewSetValueOfSlice(ctx, elements)
}

/*
DATA SOURCE MODEL
*/

type dataSourceWarehouseSQLAuditSettingsModel struct {
	baseWarehouseSQLAuditSettingsModel

	Timeouts timeoutsd.Value `tfsdk:"timeouts"`
}

/*
RESOURCE MODEL
*/

type resourceWarehouseSQLAuditSettingsModel struct {
	baseWarehouseSQLAuditSettingsModel

	Timeouts timeoutsr.Value `tfsdk:"timeouts"`
}

/*
REQUEST UPDATE MODEL
*/

type requestUpdateWarehouseSQLAuditSettings struct {
	fabwarehouse.SQLAuditSettingsUpdate
}

func (to *requestUpdateWarehouseSQLAuditSettings) set(from resourceWarehouseSQLAuditSettingsModel) {
	to.State = (*fabwarehouse.AuditSettingsState)(from.State.ValueStringPointer())
	to.RetentionDays = from.RetentionDays.ValueInt32Pointer()
}

/*
REQUEST SET AUDIT ACTIONS AND GROUPS
*/

type requestSetAuditActionsAndGroups struct {
	AuditActionsAndGroups []string
}

func (to *requestSetAuditActionsAndGroups) set(ctx context.Context, from resourceWarehouseSQLAuditSettingsModel) diag.Diagnostics {
	values, diags := from.AuditActionsAndGroups.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.AuditActionsAndGroups = make([]string, len(values))
	for i, v := range values {
		to.AuditActionsAndGroups[i] = v.ValueString()
	}

	return nil
}
