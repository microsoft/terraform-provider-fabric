// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package mirroreddatabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabmirroreddatabase "github.com/microsoft/fabric-sdk-go/fabric/mirroreddatabase"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type mirroredDatabasePropertiesModel struct {
	DefaultSchema         types.String                                                                     `tfsdk:"default_schema"`
	OneLakeTablesPath     types.String                                                                     `tfsdk:"onelake_tables_path"`
	SQLEndpointProperties supertypes.SingleNestedObjectValueOf[mirroredDatabaseSQLEndpointPropertiesModel] `tfsdk:"sql_endpoint_properties"`
}

func (to *mirroredDatabasePropertiesModel) set(ctx context.Context, from fabmirroreddatabase.Properties) diag.Diagnostics {
	to.DefaultSchema = types.StringPointerValue(from.DefaultSchema)
	to.OneLakeTablesPath = types.StringPointerValue(from.OneLakeTablesPath)

	sqlEndpointProperties := supertypes.NewSingleNestedObjectValueOfNull[mirroredDatabaseSQLEndpointPropertiesModel](ctx)

	if from.SQLEndpointProperties != nil {
		sqlEndpointPropertiesModel := &mirroredDatabaseSQLEndpointPropertiesModel{}
		sqlEndpointPropertiesModel.set(*from.SQLEndpointProperties)

		if diags := sqlEndpointProperties.Set(ctx, sqlEndpointPropertiesModel); diags.HasError() {
			return diags
		}
	}

	to.SQLEndpointProperties = sqlEndpointProperties

	return nil
}

type mirroredDatabaseSQLEndpointPropertiesModel struct {
	ProvisioningStatus types.String     `tfsdk:"provisioning_status"` // PossibleSQLEndpointProvisioningStatusValues
	ConnectionString   types.String     `tfsdk:"connection_string"`
	ID                 customtypes.UUID `tfsdk:"id"`
}

func (to *mirroredDatabaseSQLEndpointPropertiesModel) set(from fabmirroreddatabase.SQLEndpointProperties) {
	to.ProvisioningStatus = types.StringPointerValue((*string)(from.ProvisioningStatus))
	to.ConnectionString = types.StringPointerValue(from.ConnectionString)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
}
