// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

type ResourceFabricItemPropertiesModel[Ttfprop, Titemprop any] struct {
	FabricItemPropertiesModel[Ttfprop, Titemprop]
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
