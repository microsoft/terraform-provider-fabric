// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type ResourceFabricItemConfigPropertiesModel[Ttfprop, Titemprop, Ttfconfig, Titemconfig any] struct {
	FabricItemPropertiesModel[Ttfprop, Titemprop]
	Configuration supertypes.SingleNestedObjectValueOf[Ttfconfig] `tfsdk:"configuration"`
	Timeouts      timeouts.Value                                  `tfsdk:"timeouts"`
}
