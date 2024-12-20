// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventhouse

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabeventhouse "github.com/microsoft/fabric-sdk-go/fabric/eventhouse"
)

type eventhousePropertiesModel struct {
	IngestionServiceURI types.String                   `tfsdk:"ingestion_service_uri"`
	QueryServiceURI     types.String                   `tfsdk:"query_service_uri"`
	DatabaseIDs         supertypes.ListValueOf[string] `tfsdk:"database_ids"`
}

func (to *eventhousePropertiesModel) set(ctx context.Context, from *fabeventhouse.Properties) {
	to.IngestionServiceURI = types.StringPointerValue(from.IngestionServiceURI)
	to.QueryServiceURI = types.StringPointerValue(from.QueryServiceURI)
	to.DatabaseIDs = supertypes.NewListValueOfSlice(ctx, from.DatabasesItemIDs)
}
