// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package snowflakedatabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabsnowflakedatabase "github.com/microsoft/fabric-sdk-go/fabric/snowflakedatabase"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type snowflakeDatabasePropertiesModel struct {
	ConnectionID          customtypes.UUID                                                                  `tfsdk:"connection_id"`
	DefaultSchema         types.String                                                                      `tfsdk:"default_schema"`
	OnelakeTablesPath     types.String                                                                      `tfsdk:"onelake_tables_path"`
	SnowflakeAccountURL   types.String                                                                      `tfsdk:"snowflake_account_url"`
	SnowflakeDatabaseName types.String                                                                      `tfsdk:"snowflake_database_name"`
	SnowflakeVolumePath   types.String                                                                      `tfsdk:"snowflake_volume_path"`
	SQLEndpointProperties supertypes.SingleNestedObjectValueOf[snowflakeDatabaseSQLEndpointPropertiesModel] `tfsdk:"sql_endpoint_properties"`
}

func (to *snowflakeDatabasePropertiesModel) set(ctx context.Context, from fabsnowflakedatabase.Properties) diag.Diagnostics {
	sqlEndpointProperties := supertypes.NewSingleNestedObjectValueOfNull[snowflakeDatabaseSQLEndpointPropertiesModel](ctx)

	if from.SQLEndpointProperties != nil {
		sqlEndpointPropertiesModel := &snowflakeDatabaseSQLEndpointPropertiesModel{}
		sqlEndpointPropertiesModel.set(*from.SQLEndpointProperties)

		if diags := sqlEndpointProperties.Set(ctx, sqlEndpointPropertiesModel); diags.HasError() {
			return diags
		}
	}

	to.ConnectionID = customtypes.NewUUIDPointerValue(from.ConnectionID)
	to.DefaultSchema = types.StringPointerValue(from.DefaultSchema)
	to.OnelakeTablesPath = types.StringPointerValue(from.OnelakeTablesPath)
	to.SnowflakeAccountURL = types.StringPointerValue(from.SnowflakeAccountURL)
	to.SnowflakeDatabaseName = types.StringPointerValue(from.SnowflakeDatabaseName)
	to.SnowflakeVolumePath = types.StringPointerValue(from.SnowflakeVolumePath)
	to.SQLEndpointProperties = sqlEndpointProperties

	return nil
}

type snowflakeDatabaseSQLEndpointPropertiesModel struct {
	ID                 customtypes.UUID `tfsdk:"id"`
	ConnectionString   types.String     `tfsdk:"connection_string"`
	ProvisioningStatus types.String     `tfsdk:"provisioning_status"`
}

func (to *snowflakeDatabaseSQLEndpointPropertiesModel) set(from fabsnowflakedatabase.SQLEndpointProperties) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.ConnectionString = types.StringPointerValue(from.ConnectionString)
	to.ProvisioningStatus = types.StringPointerValue((*string)(from.ProvisioningStatus))
}
