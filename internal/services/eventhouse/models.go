// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package eventhouse

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	fabeventhouse "github.com/microsoft/fabric-sdk-go/fabric/eventhouse"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type eventhousePropertiesModel struct {
	IngestionServiceURI     types.String                  `tfsdk:"ingestion_service_uri"`
	QueryServiceURI         types.String                  `tfsdk:"query_service_uri"`
	DatabaseIDs             supertypes.SetValueOf[string] `tfsdk:"database_ids"`
	MinimumConsumptionUnits types.Float64                 `tfsdk:"minimum_consumption_units"`
}

func (to *eventhousePropertiesModel) set(ctx context.Context, from fabeventhouse.Properties) {
	to.IngestionServiceURI = types.StringPointerValue(from.IngestionServiceURI)
	to.QueryServiceURI = types.StringPointerValue(from.QueryServiceURI)
	to.DatabaseIDs = supertypes.NewSetValueOfSlice(ctx, from.DatabasesItemIDs)
	to.MinimumConsumptionUnits = types.Float64PointerValue(from.MinimumConsumptionUnits)
}
