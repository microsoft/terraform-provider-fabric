// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
)

type datasourceOnPremisesGatewayModel struct {
	onPremisesGatewayModelBase

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type datasourceVirtualNetworkGatewayModel struct {
	virtualNetworkGatewayModelBase

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type datasourceOnPremisesGatewayPersonalModel struct {
	onPremisesGatewayPersonalModelBase

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
