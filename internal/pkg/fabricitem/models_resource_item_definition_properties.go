// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type ResourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop any] struct {
	FabricItemPropertiesModel[Ttfprop, Titemprop]
	Format                  types.String                                                             `tfsdk:"format"`
	DefinitionUpdateEnabled types.Bool                                                               `tfsdk:"definition_update_enabled"`
	Definition              supertypes.MapNestedObjectValueOf[resourceFabricItemDefinitionPartModel] `tfsdk:"definition"`
	Timeouts                timeouts.Value                                                           `tfsdk:"timeouts"`
}
