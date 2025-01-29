// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"

type ResourceVirtualNetworkGatewayModel struct {
	virtualNetworkGatewayModelBase

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
