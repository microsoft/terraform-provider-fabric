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
	WorkspaceID                     customtypes.UUID                                                       `tfsdk:"workspace_id"`
	KeyIdentifier                   types.String                                                           `tfsdk:"key_identifier"`
	EncryptionStatus                types.String                                                           `tfsdk:"encryption_status"`
	PreviousEncryptionDetail        supertypes.SingleNestedObjectValueOf[encryptionDetailModel]            `tfsdk:"previous_encryption_detail"`
	WorkspaceEncryptionItemsDetails supertypes.SetNestedObjectValueOf[workspaceEncryptionItemsDetailModel] `tfsdk:"workspace_encryption_items_details"`
}

type encryptionDetailModel struct {
	EncryptionStatus types.String `tfsdk:"encryption_status"`
	KeyIdentifier    types.String `tfsdk:"key_identifier"`
}

type workspaceEncryptionItemsDetailModel struct {
	EncryptionStatus types.String                                                    `tfsdk:"encryption_status"`
	Items            supertypes.SetNestedObjectValueOf[workspaceEncryptionItemModel] `tfsdk:"items"`
}

type workspaceEncryptionItemModel struct {
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
	Type        types.String     `tfsdk:"type"`
}

func (to *baseWorkspaceEncryptionModel) set(ctx context.Context, workspaceID string, from fabcore.WorkspaceEncryptionDetail) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.KeyIdentifier = types.StringNull()
	to.EncryptionStatus = types.StringValue(string(encryptionStatus(from)))
	to.PreviousEncryptionDetail = supertypes.NewSingleNestedObjectValueOfNull[encryptionDetailModel](ctx)

	if from.EncryptionDetail != nil {
		to.KeyIdentifier = types.StringPointerValue(from.EncryptionDetail.KeyIdentifier)
	}

	if from.PreviousEncryptionDetail != nil {
		previousDetail := &encryptionDetailModel{}
		previousDetail.set(*from.PreviousEncryptionDetail)

		if diags := to.PreviousEncryptionDetail.Set(ctx, previousDetail); diags.HasError() {
			return diags
		}
	}

	return to.setWorkspaceEncryptionItemsDetails(ctx, from.WorkspaceEncryptionItemsDetails)
}

func (to *encryptionDetailModel) set(from fabcore.EncryptionDetail) {
	to.EncryptionStatus = types.StringPointerValue((*string)(from.EncryptionStatus))
	to.KeyIdentifier = types.StringPointerValue(from.KeyIdentifier)
}

func (to *baseWorkspaceEncryptionModel) setWorkspaceEncryptionItemsDetails(ctx context.Context, from []fabcore.WorkspaceEncryptionItemsDetail) diag.Diagnostics {
	to.WorkspaceEncryptionItemsDetails = supertypes.NewSetNestedObjectValueOfNull[workspaceEncryptionItemsDetailModel](ctx)

	if from == nil {
		return nil
	}

	details := make([]*workspaceEncryptionItemsDetailModel, 0, len(from))

	for _, detail := range from {
		var detailModel workspaceEncryptionItemsDetailModel

		if diags := detailModel.set(ctx, detail); diags.HasError() {
			return diags
		}

		details = append(details, &detailModel)
	}

	return to.WorkspaceEncryptionItemsDetails.Set(ctx, details)
}

func (to *workspaceEncryptionItemsDetailModel) set(ctx context.Context, from fabcore.WorkspaceEncryptionItemsDetail) diag.Diagnostics {
	to.EncryptionStatus = types.StringPointerValue((*string)(from.EncryptionStatus))
	to.Items = supertypes.NewSetNestedObjectValueOfNull[workspaceEncryptionItemModel](ctx)

	if from.Items == nil {
		return nil
	}

	items := make([]*workspaceEncryptionItemModel, 0, len(from.Items))

	for _, item := range from.Items {
		var itemModel workspaceEncryptionItemModel

		itemModel.set(item)
		items = append(items, &itemModel)
	}

	return to.Items.Set(ctx, items)
}

func (to *workspaceEncryptionItemModel) set(from fabcore.WorkspaceEncryptionItem) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Type = types.StringPointerValue(from.Type)
}

// A workspace that never had a customer-managed key can omit the detail entirely, which is equivalent to Disabled.
func encryptionStatus(from fabcore.WorkspaceEncryptionDetail) fabcore.WorkspaceEncryptionStatus {
	if from.EncryptionDetail == nil || from.EncryptionDetail.EncryptionStatus == nil {
		return fabcore.WorkspaceEncryptionStatusDisabled
	}

	return *from.EncryptionDetail.EncryptionStatus
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

func (to *requestAssignWorkspaceEncryption) set(from resourceWorkspaceEncryptionModel) {
	to.KeyIdentifier = from.KeyIdentifier.ValueStringPointer()
}
