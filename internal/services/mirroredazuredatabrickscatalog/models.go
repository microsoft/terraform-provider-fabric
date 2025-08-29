// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroredazuredatabrickscatalog

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabmirroredazuredatabrickscatalog "github.com/microsoft/fabric-sdk-go/fabric/mirroredazuredatabrickscatalog"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type mirroredAzureDatabricksCatalogConfigurationModel struct {
	CatalogName                     types.String     `tfsdk:"catalog_name"`
	DatabricksWorkspaceConnectionID customtypes.UUID `tfsdk:"databricks_workspace_connection_id"`
	MirroringMode                   types.String     `tfsdk:"mirroring_mode"`
	StorageConnectionID             customtypes.UUID `tfsdk:"storage_connection_id"`
}

type mirroredAzureDatabricksCatalogPropertiesModel struct {
	AutoSync                        types.String                                                     `tfsdk:"auto_sync"`
	CatalogName                     types.String                                                     `tfsdk:"catalog_name"`
	DatabricksWorkspaceConnectionID customtypes.UUID                                                 `tfsdk:"databricks_workspace_connection_id"`
	MirrorStatus                    types.String                                                     `tfsdk:"mirror_status"`
	MirroringMode                   types.String                                                     `tfsdk:"mirroring_mode"`
	OneLakeTablesPath               types.String                                                     `tfsdk:"onelake_tables_path"`
	SQLEndpointProperties           supertypes.SingleNestedObjectValueOf[sqlEndpointPropertiesModel] `tfsdk:"sql_endpoint_properties"`
	StorageConnectionID             customtypes.UUID                                                 `tfsdk:"storage_connection_id"`
	SyncDetails                     supertypes.SingleNestedObjectValueOf[syncDetailsModel]           `tfsdk:"sync_details"`
}

type sqlEndpointPropertiesModel struct {
	ID               customtypes.UUID `tfsdk:"id"`
	ConnectionString types.String     `tfsdk:"connection_string"`
}

type syncDetailsModel struct {
	LastSyncDateTime timetypes.RFC3339                                    `tfsdk:"last_sync_date_time"`
	Status           types.String                                         `tfsdk:"status"`
	ErrorInfo        supertypes.SingleNestedObjectValueOf[errorInfoModel] `tfsdk:"error_info"`
}

type errorInfoModel struct {
	ErrorCode    types.String `tfsdk:"error_code"`
	ErrorDetails types.String `tfsdk:"error_details"`
	ErrorMessage types.String `tfsdk:"error_message"`
}

func (to *mirroredAzureDatabricksCatalogPropertiesModel) set(ctx context.Context, from fabmirroredazuredatabrickscatalog.Properties) diag.Diagnostics {
	to.AutoSync = types.StringPointerValue((*string)(from.AutoSync))
	to.CatalogName = types.StringPointerValue(from.CatalogName)
	to.DatabricksWorkspaceConnectionID = customtypes.NewUUIDPointerValue(from.DatabricksWorkspaceConnectionID)
	to.MirrorStatus = types.StringPointerValue((*string)(from.MirrorStatus))
	to.MirroringMode = types.StringPointerValue((*string)(from.MirroringMode))
	to.OneLakeTablesPath = types.StringPointerValue(from.OneLakeTablesPath)
	to.StorageConnectionID = customtypes.NewUUIDPointerValue(from.StorageConnectionID)
	sqlEndpointProperties := supertypes.NewSingleNestedObjectValueOfNull[sqlEndpointPropertiesModel](ctx)
	syncDetails := supertypes.NewSingleNestedObjectValueOfNull[syncDetailsModel](ctx)

	if from.SQLEndpointProperties != nil {
		sqlEndpointPropertiesModel := &sqlEndpointPropertiesModel{}

		sqlEndpointPropertiesModel.set(*from.SQLEndpointProperties)

		if diags := sqlEndpointProperties.Set(ctx, sqlEndpointPropertiesModel); diags.HasError() {
			return diags
		}
	}

	to.SQLEndpointProperties = sqlEndpointProperties

	if from.SyncDetails != nil {
		syncDetailsModel := &syncDetailsModel{}

		syncDetailsModel.set(ctx, *from.SyncDetails)

		if diags := syncDetails.Set(ctx, syncDetailsModel); diags.HasError() {
			return diags
		}
	}

	to.SyncDetails = syncDetails

	return nil
}

func (to *sqlEndpointPropertiesModel) set(from fabmirroredazuredatabrickscatalog.SQLEndpointProperties) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.ConnectionString = types.StringPointerValue(from.ConnectionString)
}

func (to *syncDetailsModel) set(ctx context.Context, from fabmirroredazuredatabrickscatalog.SyncDetails) {
	to.LastSyncDateTime = timetypes.NewRFC3339TimePointerValue(from.LastSyncDateTime)
	to.Status = types.StringPointerValue((*string)(from.Status))

	errorInfo := supertypes.NewSingleNestedObjectValueOfNull[errorInfoModel](ctx)
	if from.ErrorInfo != nil {
		errorInfoModel := &errorInfoModel{
			ErrorCode:    types.StringPointerValue(from.ErrorInfo.ErrorCode),
			ErrorDetails: types.StringPointerValue(from.ErrorInfo.ErrorDetails),
			ErrorMessage: types.StringPointerValue(from.ErrorInfo.ErrorMessage),
		}

		if diags := errorInfo.Set(ctx, errorInfoModel); diags.HasError() {
			return
		}
	}

	to.ErrorInfo = errorInfo
}
