// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

type resourceWorkspaceModel struct {
	baseWorkspaceInfoModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type requestCreateWorkspace struct {
	fabcore.CreateWorkspaceRequest
}

func (to *requestCreateWorkspace) set(from resourceWorkspaceModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
	to.CapacityID = from.CapacityID.ValueStringPointer()
}

type requestUpdateWorkspace struct {
	fabcore.UpdateWorkspaceRequest
}

func (to *requestUpdateWorkspace) set(from resourceWorkspaceModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
}

type assignWorkspaceToCapacityRequest struct {
	fabcore.AssignWorkspaceToCapacityRequest
}

func (to *assignWorkspaceToCapacityRequest) set(from resourceWorkspaceModel) {
	to.CapacityID = from.CapacityID.ValueStringPointer()
}
