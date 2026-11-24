// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakediagnosticsetting

import (
	"context"
	"fmt"

	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseOneLakeDiagnosticSettingModel struct {
	WorkspaceID customtypes.UUID                                                        `tfsdk:"workspace_id"`
	Status      types.String                                                            `tfsdk:"status"`
	Destination supertypes.SingleNestedObjectValueOf[oneLakeDiagnosticDestinationModel] `tfsdk:"destination"`
}

type oneLakeDiagnosticDestinationModel struct {
	Lakehouse supertypes.SingleNestedObjectValueOf[oneLakeDiagnosticLakehouseModel] `tfsdk:"lakehouse"`
}

type oneLakeDiagnosticLakehouseModel struct {
	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
	ItemID      customtypes.UUID `tfsdk:"item_id"`
}

func (to *baseOneLakeDiagnosticSettingModel) set(ctx context.Context, workspaceID string, from fabcore.OneLakeDiagnosticSettings) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.Destination = supertypes.NewSingleNestedObjectValueOfNull[oneLakeDiagnosticDestinationModel](ctx)

	if from.Status == nil {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(common.ErrorReadHeader, "The OneLake settings response omitted the diagnostics status."),
		}
	}

	to.Status = types.StringPointerValue(from.Status)

	if from.Destination == nil {
		if *from.Status == diagnosticStatusEnabled {
			return diag.Diagnostics{
				diag.NewErrorDiagnostic(common.ErrorReadHeader, "The OneLake settings response omitted the enabled diagnostics destination."),
			}
		}

		return nil
	}

	destination, ok := from.Destination.(*fabcore.LakehouseOneLakeDiagnosticSettingsDestination)
	if !ok {
		return unsupportedDestinationDiagnostic(from.Destination.GetOneLakeDiagnosticSettingsDestinationInfo().Type)
	}

	if destination.Type == nil || *destination.Type != destinationTypeLakehouse {
		return unsupportedDestinationDiagnostic(destination.Type)
	}

	if destination.Lakehouse == nil {
		if *from.Status == diagnosticStatusEnabled {
			return diag.Diagnostics{
				diag.NewErrorDiagnostic(common.ErrorReadHeader, "The OneLake settings response omitted the enabled diagnostics Lakehouse reference."),
			}
		}

		return nil
	}

	itemReference, ok := destination.Lakehouse.(*fabcore.ItemReferenceByID)
	if !ok {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(
				common.ErrorReadHeader,
				fmt.Sprintf("Unsupported OneLake diagnostic Lakehouse reference type %T.", destination.Lakehouse),
			),
		}
	}

	if itemReference.ReferenceType == nil {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(common.ErrorReadHeader, "The OneLake settings response omitted the destination Lakehouse reference type."),
		}
	}

	if *itemReference.ReferenceType != fabcore.ItemReferenceTypeByID {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(
				common.ErrorReadHeader,
				fmt.Sprintf("Unsupported OneLake diagnostic Lakehouse reference type %q.", string(*itemReference.ReferenceType)),
			),
		}
	}

	if itemReference.WorkspaceID == nil || itemReference.ItemID == nil {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(common.ErrorReadHeader, "The OneLake settings response omitted the destination Lakehouse workspace or item ID."),
		}
	}

	destinationModel := &oneLakeDiagnosticDestinationModel{}
	if diags := destinationModel.set(ctx, *itemReference); diags.HasError() {
		return diags
	}

	return to.Destination.Set(ctx, destinationModel)
}

func (to *oneLakeDiagnosticDestinationModel) set(ctx context.Context, from fabcore.ItemReferenceByID) diag.Diagnostics {
	to.Lakehouse = supertypes.NewSingleNestedObjectValueOfNull[oneLakeDiagnosticLakehouseModel](ctx)

	lakehouseModel := &oneLakeDiagnosticLakehouseModel{}
	lakehouseModel.set(from)

	return to.Lakehouse.Set(ctx, lakehouseModel)
}

func (to *oneLakeDiagnosticLakehouseModel) set(from fabcore.ItemReferenceByID) {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ItemID = customtypes.NewUUIDPointerValue(from.ItemID)
}

func unsupportedDestinationDiagnostic(destinationType *string) diag.Diagnostics {
	return diag.Diagnostics{
		diag.NewErrorDiagnostic(
			common.ErrorReadHeader,
			fmt.Sprintf("Unsupported OneLake diagnostic destination type %q.", stringValue(destinationType)),
		),
	}
}

func stringValue(value *string) string {
	if value == nil {
		return "<nil>"
	}

	return *value
}

type dataSourceOneLakeDiagnosticSettingModel struct {
	baseOneLakeDiagnosticSettingModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

type resourceOneLakeDiagnosticSettingModel struct {
	baseOneLakeDiagnosticSettingModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestModifyOneLakeDiagnosticSetting struct {
	fabcore.OneLakeDiagnosticSettings
}

func (to *requestModifyOneLakeDiagnosticSetting) set(ctx context.Context, from resourceOneLakeDiagnosticSettingModel) diag.Diagnostics {
	if from.Status.IsNull() || from.Status.IsUnknown() {
		return diag.Diagnostics{
			diag.NewAttributeErrorDiagnostic(path.Root("status"), common.ErrorAttConfigMissing, "A known status is required to modify OneLake diagnostics."),
		}
	}

	status := from.Status.ValueString()
	if status != diagnosticStatusEnabled && status != diagnosticStatusDisabled {
		return diag.Diagnostics{
			diag.NewAttributeErrorDiagnostic(path.Root("status"), common.ErrorInvalidValue, "Status must be either `Enabled` or `Disabled`."),
		}
	}

	to.Status = from.Status.ValueStringPointer()
	to.Destination = nil

	if from.Destination.IsNull() {
		if status == diagnosticStatusEnabled {
			return diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(path.Root("destination"), common.ErrorAttConfigMissing, "A Lakehouse destination is required when diagnostics are enabled."),
			}
		}

		to.Destination = &fabcore.LakehouseOneLakeDiagnosticSettingsDestination{
			Type: new(destinationTypeLakehouse),
		}

		return nil
	}

	if from.Destination.IsUnknown() {
		return diag.Diagnostics{
			diag.NewAttributeErrorDiagnostic(path.Root("destination"), common.ErrorInvalidValue, "The destination must be known before modifying OneLake diagnostics."),
		}
	}

	destination, diags := from.Destination.Get(ctx)
	if diags.HasError() {
		return diags
	}

	if destination.Lakehouse.IsNull() || destination.Lakehouse.IsUnknown() {
		return diag.Diagnostics{
			diag.NewAttributeErrorDiagnostic(path.Root("destination").AtName("lakehouse"), common.ErrorAttConfigMissing, "A Lakehouse reference is required when a destination is configured."),
		}
	}

	lakehouse, diags := destination.Lakehouse.Get(ctx)
	if diags.HasError() {
		return diags
	}

	if lakehouse.WorkspaceID.IsNull() || lakehouse.WorkspaceID.IsUnknown() || lakehouse.ItemID.IsNull() || lakehouse.ItemID.IsUnknown() {
		return diag.Diagnostics{
			diag.NewAttributeErrorDiagnostic(path.Root("destination").AtName("lakehouse"), common.ErrorAttConfigMissing, "The Lakehouse workspace and item IDs must be known."),
		}
	}

	to.Destination = &fabcore.LakehouseOneLakeDiagnosticSettingsDestination{
		Type: new(destinationTypeLakehouse),
		Lakehouse: &fabcore.ItemReferenceByID{
			ReferenceType: azto.Ptr(fabcore.ItemReferenceTypeByID),
			ItemID:        lakehouse.ItemID.ValueStringPointer(),
			WorkspaceID:   lakehouse.WorkspaceID.ValueStringPointer(),
		},
	}

	return nil
}
