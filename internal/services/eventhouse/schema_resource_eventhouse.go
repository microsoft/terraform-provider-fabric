// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventhouse

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getResourceEventhousePropertiesAttributes(ctx context.Context) map[string]schema.Attribute {
	result := map[string]schema.Attribute{
		"ingestion_service_uri": schema.StringAttribute{
			MarkdownDescription: "Ingestion service URI.",
			Computed:            true,
		},
		"query_service_uri": schema.StringAttribute{
			MarkdownDescription: "Query service URI.",
			Computed:            true,
		},
		"database_ids": schema.ListAttribute{
			MarkdownDescription: "List of all KQL Database children IDs.",
			Computed:            true,
			CustomType:          supertypes.NewListTypeOf[string](ctx),
		},
	}

	return result
}

func getResourceEventhouseConfigurationAttributes() map[string]schema.Attribute {
	possibleMinimumConsumptionUnitsValues := []float64{0, 2.25, 4.25, 8.5, 13, 18, 26, 34, 50}
	customMin := float64(51)
	customMax := float64(322)

	return map[string]schema.Attribute{
		"minimum_consumption_units": schema.Float64Attribute{
			MarkdownDescription: "When activated, the eventhouse is always available at the selected minimum level and you pay at least the minimum compute selected. Accepted values: " + utils.ConvertStringSlicesToString(possibleMinimumConsumptionUnitsValues, true, true) + " or any number between `" + fmt.Sprintf("%v", customMin) + "` and `" + fmt.Sprintf("%v", customMax) + "`. For more information, see [minimum consumption](https://learn.microsoft.com/fabric/real-time-intelligence/eventhouse#minimum-consumption)",
			Required:            true,
			Validators: []validator.Float64{
				float64validator.Any(
					float64validator.OneOf(possibleMinimumConsumptionUnitsValues...),
					float64validator.Between(customMin, customMax),
				),
			},
			PlanModifiers: []planmodifier.Float64{
				float64planmodifier.RequiresReplace(),
			},
		},
	}
}
