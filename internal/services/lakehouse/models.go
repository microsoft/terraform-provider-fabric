// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fablakehouse "github.com/microsoft/fabric-sdk-go/fabric/lakehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseLakehouseModel struct {
	WorkspaceID customtypes.UUID                                               `tfsdk:"workspace_id"`
	ID          customtypes.UUID                                               `tfsdk:"id"`
	DisplayName types.String                                                   `tfsdk:"display_name"`
	Description types.String                                                   `tfsdk:"description"`
	Properties  supertypes.SingleNestedObjectValueOf[lakehousePropertiesModel] `tfsdk:"properties"`
}

func (to *baseLakehouseModel) set(ctx context.Context, from fablakehouse.Lakehouse) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)

	properties := supertypes.NewSingleNestedObjectValueOfNull[lakehousePropertiesModel](ctx)

	if from.Properties != nil {
		propertiesModel := &lakehousePropertiesModel{}

		if diags := propertiesModel.set(ctx, from.Properties); diags.HasError() {
			return diags
		}

		if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
			return diags
		}
	}

	to.Properties = properties

	return nil
}

type lakehousePropertiesModel struct {
	OneLakeFilesPath      types.String                                                              `tfsdk:"onelake_files_path"`
	OneLakeTablesPath     types.String                                                              `tfsdk:"onelake_tables_path"`
	SQLEndpointProperties supertypes.SingleNestedObjectValueOf[lakehouseSQLEndpointPropertiesModel] `tfsdk:"sql_endpoint_properties"`
	DefaultSchema         types.String                                                              `tfsdk:"default_schema"`
}

func (to *lakehousePropertiesModel) set(ctx context.Context, from *fablakehouse.Properties) diag.Diagnostics {
	sqlEndpointProperties := supertypes.NewSingleNestedObjectValueOfNull[lakehouseSQLEndpointPropertiesModel](ctx)

	if from.SQLEndpointProperties != nil {
		sqlEndpointPropertiesModel := &lakehouseSQLEndpointPropertiesModel{}
		sqlEndpointPropertiesModel.set(from.SQLEndpointProperties)

		if diags := sqlEndpointProperties.Set(ctx, sqlEndpointPropertiesModel); diags.HasError() {
			return diags
		}
	}

	to.SQLEndpointProperties = sqlEndpointProperties
	to.OneLakeFilesPath = types.StringPointerValue(from.OneLakeFilesPath)
	to.OneLakeTablesPath = types.StringPointerValue(from.OneLakeTablesPath)
	to.DefaultSchema = types.StringPointerValue(from.DefaultSchema)

	return nil
}

type lakehouseSQLEndpointPropertiesModel struct {
	ID                 customtypes.UUID `tfsdk:"id"`
	ConnectionString   types.String     `tfsdk:"connection_string"`
	ProvisioningStatus types.String     `tfsdk:"provisioning_status"`
}

func (to *lakehouseSQLEndpointPropertiesModel) set(from *fablakehouse.SQLEndpointProperties) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.ConnectionString = types.StringPointerValue(from.ConnectionString)
	to.ProvisioningStatus = types.StringPointerValue((*string)(from.ProvisioningStatus))
}
