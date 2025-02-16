// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"

	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type resourceWorkspaceManagedPrivateEndpointModel struct {
	baseWorkspaceManagedPrivateEndpointModel
	WorkspaceID    customtypes.UUID `tfsdk:"workspace_id"`
	RequestMessage types.String     `tfsdk:"request_message"`
	Timeouts       timeouts.Value   `tfsdk:"timeouts"`
}

type requestCreateWorkspaceManagedPrivateEndpoint struct {
	fabcore.CreateManagedPrivateEndpointRequest
}

func (to *requestCreateWorkspaceManagedPrivateEndpoint) set(from resourceWorkspaceManagedPrivateEndpointModel) {
	to.Name = from.Name.ValueStringPointer()
	to.TargetPrivateLinkResourceID = from.TargetPrivateLinkResourceID.ValueStringPointer()
	to.TargetSubresourceType = from.TargetSubresourceType.ValueStringPointer()
	to.RequestMessage = from.RequestMessage.ValueStringPointer()
}
