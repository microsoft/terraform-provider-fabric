// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package monitoringdatacollection_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabmonitoring "github.com/microsoft/fabric-sdk-go/fabric/monitoring"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
	"github.com/microsoft/terraform-provider-fabric/internal/services/monitoringdatacollection"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

func configuredResource(t *testing.T, server *fabfake.ServerFactory, entity fabmonitoring.Monitoring) (resource.Resource, tfsdk.State) {
	t.Helper()

	options := &fabric.ClientOptions{}
	options.Transport = fabfake.NewServerFactoryTransport(server)
	options.Retry.MaxRetries = -1
	client, err := fabric.NewClient(&azfake.TokenCredential{}, nil, options)
	require.NoError(t, err)

	r := monitoringdatacollection.NewResourceMonitoringDataCollection()
	var configured resource.ConfigureResponse
	r.(resource.ResourceWithConfigure).Configure(t.Context(), resource.ConfigureRequest{
		ProviderData: &pconfig.ProviderData{FabricClient: client, Preview: true, Timeout: time.Second},
	}, &configured)
	require.False(t, configured.Diagnostics.HasError(), configured.Diagnostics)

	var schema resource.SchemaResponse
	r.Schema(t.Context(), resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	model := struct {
		ID           customtypes.UUID `tfsdk:"id"`
		WorkspaceID  customtypes.UUID `tfsdk:"workspace_id"`
		MonitoringID customtypes.UUID `tfsdk:"monitoring_id"`
		Enabled      types.Bool       `tfsdk:"enabled"`
		Timeouts     types.Object     `tfsdk:"timeouts"`
	}{
		ID:           customtypes.NewUUIDPointerValue(entity.ID),
		WorkspaceID:  customtypes.NewUUIDPointerValue(entity.WorkspaceID),
		MonitoringID: customtypes.NewUUIDPointerValue(entity.ID),
		Enabled:      types.BoolValue(true),
		Timeouts: types.ObjectNull(map[string]attr.Type{
			"create": types.StringType, "read": types.StringType,
			"update": types.StringType, "delete": types.StringType,
		}),
	}
	require.False(t, state.Set(t.Context(), model).HasError())

	return r, state
}

func TestUnit_MonitoringDataCollectionResource_ReadContract(t *testing.T) {
	t.Parallel()

	enabled := fabmonitoring.DataCollectionStateEnabled
	disabled := fabmonitoring.DataCollectionStateDisabled

	unknown := fabmonitoring.DataCollectionState("FutureState")
	for _, tc := range []struct {
		name       string
		properties *fabmonitoring.Properties
		wantError  bool
		enabled    bool
	}{
		{"enabled", &fabmonitoring.Properties{DataCollectionState: &enabled}, false, true},
		{"disabled_drift", &fabmonitoring.Properties{DataCollectionState: &disabled}, false, false},
		{"missing_properties", nil, true, false},
		{"missing_state", &fabmonitoring.Properties{}, true, false},
		{"unknown_state", &fabmonitoring.Properties{DataCollectionState: &unknown}, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := fakes.NewFakeServer()
			entity := fakes.NewRandomMonitoring()
			entity.Properties = tc.properties
			server.Upsert(entity)
			r, state := configuredResource(t, server.ServerFactory, entity)
			response := resource.ReadResponse{State: state}
			r.Read(t.Context(), resource.ReadRequest{State: state}, &response)
			require.Equal(t, tc.wantError, response.Diagnostics.HasError(), response.Diagnostics)

			if !tc.wantError {
				var actual types.Bool
				require.False(t, response.State.GetAttribute(t.Context(), path.Root("enabled"), &actual).HasError())
				require.Equal(t, tc.enabled, actual.ValueBool())
			}
		})
	}
}

func TestUnit_MonitoringDataCollectionResource_DeletedParent(t *testing.T) {
	t.Parallel()

	server := fakes.NewFakeServer()
	entity := fakes.NewRandomMonitoring()
	r, state := configuredResource(t, server.ServerFactory, entity)
	read := resource.ReadResponse{State: state}
	r.Read(t.Context(), resource.ReadRequest{State: state}, &read)
	require.False(t, read.Diagnostics.HasError(), read.Diagnostics)
	require.True(t, read.State.Raw.IsNull())

	deleted := resource.DeleteResponse{State: state}
	r.Delete(t.Context(), resource.DeleteRequest{State: state}, &deleted)
	require.False(t, deleted.Diagnostics.HasError(), deleted.Diagnostics)
}

func TestUnit_MonitoringDataCollectionResource_PartialCreate(t *testing.T) {
	t.Parallel()

	server := fakes.NewFakeServer()
	entity := fakes.NewRandomMonitoring()
	server.Upsert(entity)
	server.ServerFactory.Monitoring.ItemsServer.GetMonitoring = func(context.Context, string, string, *fabmonitoring.ItemsClientGetMonitoringOptions) (azfake.Responder[fabmonitoring.ItemsClientGetMonitoringResponse], azfake.ErrorResponder) {
		var response azfake.Responder[fabmonitoring.ItemsClientGetMonitoringResponse]
		var errResponse azfake.ErrorResponder
		errResponse.SetError(fabfake.SetResponseError(http.StatusForbidden, "InsufficientPrivileges", "Read denied"))

		return response, errResponse
	}
	r, state := configuredResource(t, server.ServerFactory, entity)
	response := resource.CreateResponse{State: tfsdk.State{Schema: state.Schema}}
	r.Create(t.Context(), resource.CreateRequest{Plan: tfsdk.Plan(state)}, &response)
	require.True(t, response.Diagnostics.HasError())
	var id customtypes.UUID
	require.False(t, response.State.GetAttribute(t.Context(), path.Root("id"), &id).HasError())
	require.Equal(t, *entity.ID, id.ValueString(), "successful enablement must remain tracked after a failed read")
}

func TestUnit_MonitoringDataCollectionResource_DeleteError(t *testing.T) {
	t.Parallel()

	server := fakes.NewFakeServer()
	entity := fakes.NewRandomMonitoring()
	server.Upsert(entity)
	server.ServerFactory.Monitoring.DataCollectionServer.Disable = func(context.Context, string, string, *fabmonitoring.DataCollectionClientDisableOptions) (azfake.Responder[fabmonitoring.DataCollectionClientDisableResponse], azfake.ErrorResponder) {
		var response azfake.Responder[fabmonitoring.DataCollectionClientDisableResponse]
		var errResponse azfake.ErrorResponder
		errResponse.SetError(fabfake.SetResponseError(http.StatusForbidden, "InsufficientPrivileges", "Disable denied"))

		return response, errResponse
	}
	r, state := configuredResource(t, server.ServerFactory, entity)
	response := resource.DeleteResponse{State: state}
	r.Delete(t.Context(), resource.DeleteRequest{State: state}, &response)
	require.True(t, response.Diagnostics.HasError())
	require.False(t, response.State.Raw.IsNull())
}

func TestUnit_MonitoringDataCollectionResource_Configure(t *testing.T) {
	for _, input := range []any{nil, "invalid", &pconfig.ProviderData{Preview: false}} {
		r := monitoringdatacollection.NewResourceMonitoringDataCollection()
		var response resource.ConfigureResponse
		r.(resource.ResourceWithConfigure).Configure(t.Context(), resource.ConfigureRequest{ProviderData: input}, &response)
		require.Equal(t, input != nil, response.Diagnostics.HasError())
	}
}

func TestUnit_MonitoringDataCollectionResource_UnresolvedEnabled(t *testing.T) {
	for _, enabled := range []types.Bool{types.BoolNull(), types.BoolUnknown()} {
		server := fakes.NewFakeServer()
		entity := fakes.NewRandomMonitoring()
		r, state := configuredResource(t, server.ServerFactory, entity)
		require.False(t, state.SetAttribute(t.Context(), path.Root("enabled"), enabled).HasError())
		response := resource.CreateResponse{State: tfsdk.State{Schema: state.Schema}}
		r.Create(t.Context(), resource.CreateRequest{Plan: tfsdk.Plan(state)}, &response)
		require.True(t, response.Diagnostics.HasError())
		require.Contains(t, response.Diagnostics.Errors()[0].Summary(), "Invalid collection state")
	}
}
