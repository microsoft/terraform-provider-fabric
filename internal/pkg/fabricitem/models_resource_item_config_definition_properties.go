// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceFabricItemConfigDefinitionPropertiesModel[Ttfprop, Titemprop, Ttfconfig, Titemconfig any] struct {
	FabricItemPropertiesModel[Ttfprop, Titemprop]
	Configuration           supertypes.SingleNestedObjectValueOf[Ttfconfig]                          `tfsdk:"configuration"`
	Format                  types.String                                                             `tfsdk:"format"`
	DefinitionUpdateEnabled types.Bool                                                               `tfsdk:"definition_update_enabled"`
	Definition              supertypes.MapNestedObjectValueOf[resourceFabricItemDefinitionPartModel] `tfsdk:"definition"`
	Timeouts                timeouts.Value                                                           `tfsdk:"timeouts"`
}