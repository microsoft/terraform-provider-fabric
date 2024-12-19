// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop any] struct {
	baseFabricItemModel
	Format           types.String                                                               `tfsdk:"format"`
	OutputDefinition types.Bool                                                                 `tfsdk:"output_definition"`
	Definition       supertypes.MapNestedObjectValueOf[DataSourceFabricItemDefinitionPartModel] `tfsdk:"definition"`
	Properties       supertypes.SingleNestedObjectValueOf[Ttfprop]                              `tfsdk:"properties"`
	Timeouts         timeouts.Value                                                             `tfsdk:"timeouts"`
}

func (to *DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]) set(from FabricItemProperties[Titemprop]) { //revive:disable-line:confusing-naming
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
}

func (to *DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]) setDefinition(ctx context.Context, from fabcore.ItemDefinition) diag.Diagnostics {
	defParts := make(map[string]*DataSourceFabricItemDefinitionPartModel, len(from.Parts))

	for _, part := range from.Parts {
		newPart := &DataSourceFabricItemDefinitionPartModel{}

		if diags := newPart.Set(*part.Payload); diags.HasError() {
			return diags
		}

		defParts[*part.Path] = newPart
	}

	return to.Definition.Set(ctx, defParts)
}
