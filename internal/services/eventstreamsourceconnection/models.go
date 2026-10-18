// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package eventstreamsourceconnection

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsE "github.com/hashicorp/terraform-plugin-framework-timeouts/ephemeral/timeouts"  //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabeventstream "github.com/microsoft/fabric-sdk-go/fabric/eventstream"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceEventstreamSourceConnectionModel struct {
	baseEventstreamSourceConnectionModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

type ephemeralEventstreamSourceConnectionModel struct {
	baseEventstreamSourceConnectionModel

	Timeouts timeoutsE.Value `tfsdk:"timeouts"`
}

type baseEventstreamSourceConnectionModel struct {
	SourceID                customtypes.UUID                                      `tfsdk:"source_id"`
	EventstreamID           customtypes.UUID                                      `tfsdk:"eventstream_id"`
	WorkspaceID             customtypes.UUID                                      `tfsdk:"workspace_id"`
	EventHubName            types.String                                          `tfsdk:"event_hub_name"`
	FullyQualifiedNamespace types.String                                          `tfsdk:"fully_qualified_namespace"`
	AccessKeys              supertypes.SingleNestedObjectValueOf[accessKeysModel] `tfsdk:"access_keys"`
}

func (to *baseEventstreamSourceConnectionModel) set(ctx context.Context, workspaceID, eventstreamID, sourceID string, from fabeventstream.SourceConnectionResponse) diag.Diagnostics {
	to.SourceID = customtypes.NewUUIDValue(sourceID)
	to.EventstreamID = customtypes.NewUUIDValue(eventstreamID)
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.EventHubName = types.StringPointerValue(from.EventHubName)
	to.FullyQualifiedNamespace = types.StringPointerValue(from.FullyQualifiedNamespace)

	accessKeys := supertypes.NewSingleNestedObjectValueOfNull[accessKeysModel](ctx)

	if from.AccessKeys != nil {
		accessKeysModel := &accessKeysModel{}
		accessKeysModel.set(*from.AccessKeys)

		if diags := accessKeys.Set(ctx, accessKeysModel); diags.HasError() {
			return diags
		}
	}

	to.AccessKeys = accessKeys

	return nil
}

type accessKeysModel struct {
	PrimaryKey                types.String `tfsdk:"primary_key"`
	SecondaryKey              types.String `tfsdk:"secondary_key"`
	PrimaryConnectionString   types.String `tfsdk:"primary_connection_string"`
	SecondaryConnectionString types.String `tfsdk:"secondary_connection_string"`
}

func (to *accessKeysModel) set(from fabeventstream.AccessKeys) {
	to.PrimaryKey = types.StringPointerValue(from.PrimaryKey)
	to.SecondaryKey = types.StringPointerValue(from.SecondaryKey)
	to.PrimaryConnectionString = types.StringPointerValue(from.PrimaryConnectionString)
	to.SecondaryConnectionString = types.StringPointerValue(from.SecondaryConnectionString)
}
