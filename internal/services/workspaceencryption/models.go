// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceencryption

import (
	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseWorkspaceEncryptionModel struct {
	WorkspaceID      customtypes.UUID `tfsdk:"workspace_id"`
	KeyIdentifier    types.String     `tfsdk:"key_identifier"`
	EncryptionStatus types.String     `tfsdk:"encryption_status"`
}

// Only the current encryption detail is mapped. PreviousEncryptionDetail and workspaceEncryptionItemsDetails
// will not be supported in the provider.
func (to *baseWorkspaceEncryptionModel) set(workspaceID string, from fabcore.WorkspaceEncryptionDetail) {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.KeyIdentifier = types.StringNull()
	to.EncryptionStatus = types.StringValue(string(encryptionStatus(from)))

	if from.EncryptionDetail != nil {
		to.KeyIdentifier = types.StringPointerValue(from.EncryptionDetail.KeyIdentifier)
	}
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
