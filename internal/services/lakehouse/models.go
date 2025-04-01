// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fablakehouse "github.com/microsoft/fabric-sdk-go/fabric/lakehouse"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type lakehouseConfigurationModel struct {
	EnableSchemas types.Bool `tfsdk:"enable_schemas"`
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
		sqlEndpointPropertiesModel.set(*from.SQLEndpointProperties)

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

func (to *lakehouseSQLEndpointPropertiesModel) set(from fablakehouse.SQLEndpointProperties) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.ConnectionString = types.StringPointerValue(from.ConnectionString)
	to.ProvisioningStatus = types.StringPointerValue((*string)(from.ProvisioningStatus))
}
