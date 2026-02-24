// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package kqldatabase

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabkqldatabase "github.com/microsoft/fabric-sdk-go/fabric/kqldatabase"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type kqlDatabasePropertiesModel struct {
	DatabaseType        types.String     `tfsdk:"database_type"`
	EventhouseID        customtypes.UUID `tfsdk:"eventhouse_id"`
	IngestionServiceURI customtypes.URL  `tfsdk:"ingestion_service_uri"`
	QueryServiceURI     customtypes.URL  `tfsdk:"query_service_uri"`
	// OneLakeStandardStoragePeriod types.String     `tfsdk:"onelake_standard_storage_period"`
	// OneLakeCachingPeriod         types.String     `tfsdk:"onelake_caching_period"`
}

func (to *kqlDatabasePropertiesModel) set(from fabkqldatabase.Properties) {
	to.DatabaseType = types.StringPointerValue((*string)(from.DatabaseType))
	to.EventhouseID = customtypes.NewUUIDPointerValue(from.ParentEventhouseItemID)
	to.IngestionServiceURI = customtypes.NewURLPointerValue(from.IngestionServiceURI)
	to.QueryServiceURI = customtypes.NewURLPointerValue(from.QueryServiceURI)
	// to.OneLakeStandardStoragePeriod = types.StringPointerValue(from.OneLakeStandardStoragePeriod)
	// to.OneLakeCachingPeriod = types.StringPointerValue(from.OneLakeCachingPeriod)
}
