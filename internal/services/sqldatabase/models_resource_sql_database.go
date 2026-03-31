// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sqldatabase

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabsqldatabase "github.com/microsoft/fabric-sdk-go/fabric/sqldatabase"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type sqlDatabaseConfigurationModel struct {
	CreationMode            types.String                                                       `tfsdk:"creation_mode"`
	BackupRetentionDays     types.Int32                                                        `tfsdk:"backup_retention_days"`
	Collation               types.String                                                       `tfsdk:"collation"`
	RestorePointInTime      timetypes.RFC3339                                                  `tfsdk:"restore_point_in_time"`
	SourceDatabaseReference supertypes.SingleNestedObjectValueOf[sourceDatabaseReferenceModel] `tfsdk:"source_database_reference"`
}

type sourceDatabaseReferenceModel struct {
	ItemID        customtypes.UUID `tfsdk:"item_id"`
	ReferenceType types.String     `tfsdk:"reference_type"`
	WorkspaceID   customtypes.UUID `tfsdk:"workspace_id"`
}

type requestCreateSQLDatabasePayload struct {
	fabsqldatabase.CreationPayloadClassification
}

func (to *requestCreateSQLDatabasePayload) set(ctx context.Context, from sqlDatabaseConfigurationModel) diag.Diagnostics {
	creationMode := fabsqldatabase.CreationMode(from.CreationMode.ValueString())

	switch creationMode {
	case fabsqldatabase.CreationModeNew:
		creationPayload := fabsqldatabase.NewSQLDatabaseCreationPayload{
			CreationMode: &creationMode,
		}

		if !from.BackupRetentionDays.IsNull() && !from.BackupRetentionDays.IsUnknown() {
			creationPayload.BackupRetentionDays = from.BackupRetentionDays.ValueInt32Pointer()
		}

		if !from.Collation.IsNull() && !from.Collation.IsUnknown() {
			creationPayload.Collation = from.Collation.ValueStringPointer()
		}

		to.CreationPayloadClassification = &creationPayload
	case fabsqldatabase.CreationModeRestore:
		creationPayload := fabsqldatabase.RestoreSQLDatabaseCreationPayload{
			CreationMode: &creationMode,
		}

		if !from.RestorePointInTime.IsNull() && !from.RestorePointInTime.IsUnknown() {
			t, diags := from.RestorePointInTime.ValueRFC3339Time()
			if diags.HasError() {
				return diags
			}

			creationPayload.RestorePointInTime = new(t)
		}

		if !from.SourceDatabaseReference.IsNull() && !from.SourceDatabaseReference.IsUnknown() {
			ref, diags := getSourceDatabaseReference(ctx, from.SourceDatabaseReference)
			if diags.HasError() {
				return diags
			}

			creationPayload.SourceDatabaseReference = ref
		}

		to.CreationPayloadClassification = &creationPayload
	default:
		var diags diag.Diagnostics

		diags.AddError(
			"Unsupported SQL Database creation mode",
			fmt.Sprintf("The creation mode '%s' is not supported.", string(creationMode)),
		)

		return diags
	}

	return nil
}

func getSourceDatabaseReference(
	ctx context.Context,
	sourceDatabaseReference supertypes.SingleNestedObjectValueOf[sourceDatabaseReferenceModel],
) (fabsqldatabase.ItemReferenceClassification, diag.Diagnostics) {
	model, diags := sourceDatabaseReference.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	refType := fabsqldatabase.ItemReferenceType(model.ReferenceType.ValueString())

	switch refType {
	case fabsqldatabase.ItemReferenceTypeByID:
		ref := fabsqldatabase.ItemReferenceByID{
			ReferenceType: &refType,
			ItemID:        model.ItemID.ValueStringPointer(),
			WorkspaceID:   model.WorkspaceID.ValueStringPointer(),
		}

		return &ref, nil

	default:
		var diags diag.Diagnostics

		diags.AddError(
			"Unsupported reference type",
			fmt.Sprintf("The reference type '%s' is not supported.", string(refType)),
		)

		return nil, diags
	}
}
