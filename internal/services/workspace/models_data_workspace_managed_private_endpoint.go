// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceWorkspaceManagedPrivateEndpointModel struct {
	baseWorkspaceManagedPrivateEndpointModel
	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
	Timeouts    timeouts.Value   `tfsdk:"timeouts"`
}
