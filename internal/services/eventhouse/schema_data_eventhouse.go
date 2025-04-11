// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventhouse

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getDataSourceEventhousePropertiesAttributes(ctx context.Context) map[string]schema.Attribute {
	possibleMinimumConsumptionUnitsValues := []float64{0, 2.25, 4.25, 8.5, 13, 18, 26, 34, 50}
	customMin := float64(51)
	customMax := float64(322)

	result := map[string]schema.Attribute{
		"ingestion_service_uri": schema.StringAttribute{
			MarkdownDescription: "Ingestion service URI.",
			Computed:            true,
		},
		"query_service_uri": schema.StringAttribute{
			MarkdownDescription: "Query service URI.",
			Computed:            true,
		},
		"database_ids": schema.SetAttribute{
			MarkdownDescription: "List of all KQL Database children IDs.",
			Computed:            true,
			CustomType:          supertypes.NewSetTypeOf[string](ctx),
		},
		"minimum_consumption_units": schema.Float64Attribute{
			MarkdownDescription: "Use Minimum consumption for highly time-sensitive systems to keep the service always available at a selected minimum level. " +
				"You pay for the minimum consumption level or actual consumption if above the minimum. Supported values include" +
				utils.ConvertStringSlicesToString(
					possibleMinimumConsumptionUnitsValues,
					true,
					true,
				) + " or any number between `" + fmt.Sprintf(
				"%v",
				customMin,
			) + "` and `" + fmt.Sprintf(
				"%v",
				customMax,
			) + "`. For more information, see [minimum consumption](https://learn.microsoft.com/fabric/real-time-intelligence/eventhouse#minimum-consumption)",
			Computed: true,
		},
	}

	return result
}
