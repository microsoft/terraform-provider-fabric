// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package snowflakedatabase

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type snowflakeDatabaseConfigurationModel struct {
	ConnectionID          customtypes.UUID `tfsdk:"connection_id"`
	SnowflakeDatabaseName types.String     `tfsdk:"snowflake_database_name"`
}
