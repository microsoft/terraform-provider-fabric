// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package kqldatabase

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type kqlDatabaseConfigurationModel struct {
	DatabaseType             types.String     `tfsdk:"database_type"`
	EventhouseID             customtypes.UUID `tfsdk:"eventhouse_id"`
	InvitationToken          types.String     `tfsdk:"invitation_token"`
	InvitationTokenWO        types.String     `tfsdk:"invitation_token_wo"`
	InvitationTokenWOVersion types.Int32      `tfsdk:"invitation_token_wo_version"`
	SourceClusterURI         customtypes.URL  `tfsdk:"source_cluster_uri"`
	SourceDatabaseName       types.String     `tfsdk:"source_database_name"`
}
