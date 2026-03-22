// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type ResourceFabricItemPropertiesModel[Ttfprop, Titemprop any] struct {
	FabricItemPropertiesModel[Ttfprop, Titemprop]

	SensitivityLabelSettings supertypes.SingleNestedObjectValueOf[sensitivityLabelSettingsModel] `tfsdk:"sensitivity_label_settings"`
	Timeouts                 timeouts.Value                                                      `tfsdk:"timeouts"`
}
