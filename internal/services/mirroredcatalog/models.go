// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package mirroredcatalog

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabmirroredcatalog "github.com/microsoft/fabric-sdk-go/fabric/mirroredcatalog"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type mirroredCatalogPropertiesModel struct {
	ConnectionID          customtypes.UUID                                                                `tfsdk:"connection_id"`
	OneLakeTablesPath     types.String                                                                    `tfsdk:"onelake_tables_path"`
	Scope                 supertypes.ListValueOf[string]                                                  `tfsdk:"scope"`
	SourceType            types.String                                                                    `tfsdk:"source_type"`
	SQLEndpointProperties supertypes.SingleNestedObjectValueOf[mirroredCatalogSQLEndpointPropertiesModel] `tfsdk:"sql_endpoint_properties"`
}

func (to *mirroredCatalogPropertiesModel) set(ctx context.Context, from fabmirroredcatalog.Properties) diag.Diagnostics {
	to.ConnectionID = customtypes.NewUUIDPointerValue(from.ConnectionID)
	to.OneLakeTablesPath = types.StringPointerValue(from.OneLakeTablesPath)
	to.SourceType = types.StringPointerValue(from.SourceType)
	to.Scope = supertypes.NewListValueOfSlice(ctx, from.Scope)

	sqlEndpointProperties := supertypes.NewSingleNestedObjectValueOfNull[mirroredCatalogSQLEndpointPropertiesModel](ctx)

	if from.SQLEndpointProperties != nil {
		sqlEndpointPropertiesModel := &mirroredCatalogSQLEndpointPropertiesModel{}
		sqlEndpointPropertiesModel.set(*from.SQLEndpointProperties)

		if diags := sqlEndpointProperties.Set(ctx, sqlEndpointPropertiesModel); diags.HasError() {
			return diags
		}
	}

	to.SQLEndpointProperties = sqlEndpointProperties

	return nil
}

type mirroredCatalogSQLEndpointPropertiesModel struct {
	ProvisioningStatus types.String     `tfsdk:"provisioning_status"`
	ConnectionString   types.String     `tfsdk:"connection_string"`
	ID                 customtypes.UUID `tfsdk:"id"`
}

func (to *mirroredCatalogSQLEndpointPropertiesModel) set(from fabmirroredcatalog.SQLEndpointProperties) {
	to.ProvisioningStatus = types.StringPointerValue((*string)(from.ProvisioningStatus))
	to.ConnectionString = types.StringPointerValue(from.ConnectionString)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
}
