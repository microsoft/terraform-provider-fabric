// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
)

type DataSourceFabricItemPropertiesModel[Ttfprop, Titemprop any] struct {
	DataSourceFabricItemPropertiesBaseModel[Ttfprop, Titemprop]

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
