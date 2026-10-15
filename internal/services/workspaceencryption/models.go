// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceencryption

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseWorkspaceEncryptionModel struct {
	WorkspaceID               customtypes.UUID                                             `tfsdk:"workspace_id"`
	EncryptionDetails         supertypes.SingleNestedObjectValueOf[encryptionDetailsModel] `tfsdk:"encryption_details"`
	PreviousEncryptionDetails supertypes.SingleNestedObjectValueOf[encryptionDetailsModel] `tfsdk:"previous_encryption_details"`
}

type encryptionDetailsModel struct {
	EncryptionStatus types.String `tfsdk:"encryption_status"`
	KeyIdentifier    types.String `tfsdk:"key_identifier"`
}

func (to *baseWorkspaceEncryptionModel) set(ctx context.Context, workspaceID string, from fabcore.WorkspaceEncryptionDetail) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.EncryptionDetails = supertypes.NewSingleNestedObjectValueOfNull[encryptionDetailsModel](ctx)
	to.PreviousEncryptionDetails = supertypes.NewSingleNestedObjectValueOfNull[encryptionDetailsModel](ctx)

	if from.EncryptionDetail != nil {
		detail := &encryptionDetailsModel{}
		detail.set(*from.EncryptionDetail)

		if diags := to.EncryptionDetails.Set(ctx, detail); diags.HasError() {
			return diags
		}
	}

	if from.PreviousEncryptionDetail != nil {
		previousDetail := &encryptionDetailsModel{}
		previousDetail.set(*from.PreviousEncryptionDetail)

		if diags := to.PreviousEncryptionDetails.Set(ctx, previousDetail); diags.HasError() {
			return diags
		}
	}

	return nil
}

func (to *encryptionDetailsModel) set(from fabcore.EncryptionDetail) {
	to.EncryptionStatus = types.StringPointerValue((*string)(from.EncryptionStatus))
	to.KeyIdentifier = types.StringPointerValue(from.KeyIdentifier)
}

/*
DATA-SOURCE
*/

type dataSourceWorkspaceEncryptionModel struct {
	baseWorkspaceEncryptionModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
RESOURCE
*/

type resourceWorkspaceEncryptionModel struct {
	baseWorkspaceEncryptionModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestAssignWorkspaceEncryption struct {
	fabcore.AssignWorkspaceEncryptionRequest
}

func (to *requestAssignWorkspaceEncryption) set(ctx context.Context, from resourceWorkspaceEncryptionModel) diag.Diagnostics {
	detail, diags := from.EncryptionDetails.Get(ctx)
	if diags.HasError() {
		return diags
	}

	if detail != nil {
		to.KeyIdentifier = detail.KeyIdentifier.ValueStringPointer()
	}

	return nil
}
