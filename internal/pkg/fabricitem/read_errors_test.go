// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

type itemReadTransport struct {
	status          int
	code            string
	message         string
	calls           int
	successfulReads int
	allowCreate     bool
}

func (tr *itemReadTransport) Do(req *http.Request) (*http.Response, error) {
	tr.calls++

	if (tr.allowCreate && req.Method == http.MethodPost) || (tr.successfulReads > 0 && req.Method == http.MethodGet) {
		status := http.StatusOK
		if req.Method == http.MethodPost {
			status = http.StatusCreated
		} else {
			tr.successfulReads--
		}

		return &http.Response{
			StatusCode: status,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(strings.NewReader(
				`{"id":"11111111-1111-4111-8111-111111111111","workspaceId":"22222222-2222-4222-8222-222222222222","displayName":"existing-item","description":"last-known description","type":"Lakehouse","tags":[]}`,
			)),
			Request: req,
		}, nil
	}

	return &http.Response{
		StatusCode: tr.status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body: io.NopCloser(strings.NewReader(fmt.Sprintf(
			`{"requestId":"read-request","errorCode":%q,"message":%q,"isRetriable":false}`,
			tr.code, tr.message,
		))),
		Request: req,
	}, nil
}

type itemReadCredential struct{}

func (itemReadCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "test-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

type itemReadProperties struct {
	Value types.String `tfsdk:"value"`
}

func itemReadResources(propertyDiags diag.Diagnostics) map[string]resource.Resource {
	attributes := map[string]schema.Attribute{
		"value": schema.StringAttribute{Computed: true},
	}
	get := func(ctx context.Context, client fabric.Client, workspaceID, itemID string, to *fabricitem.FabricItemProperties[struct{}]) error {
		resp, err := fabcore.NewClientFactoryWithClient(client).NewItemsClient().GetItem(ctx, workspaceID, itemID, nil)
		if err != nil {
			return fmt.Errorf("typed item read: %w", err)
		}

		to.Item = resp.Item

		return nil
	}

	return map[string]resource.Resource{
		"item":       fabricitem.NewResourceFabricItem(fabricitem.ResourceFabricItem{}),
		"definition": fabricitem.NewResourceFabricItemDefinition(fabricitem.ResourceFabricItemDefinition{}),
		"properties": fabricitem.NewResourceFabricItemProperties(fabricitem.ResourceFabricItemProperties[itemReadProperties, struct{}]{
			PropertiesAttributes: attributes,
			ItemGetter: func(ctx context.Context, client fabric.Client, model fabricitem.ResourceFabricItemPropertiesModel[itemReadProperties, struct{}], to *fabricitem.FabricItemProperties[struct{}]) error {
				return get(ctx, client, model.WorkspaceID.ValueString(), model.ID.ValueString(), to)
			},
			PropertiesSetter: func(context.Context, *struct{}, *fabricitem.ResourceFabricItemPropertiesModel[itemReadProperties, struct{}]) diag.Diagnostics {
				return propertyDiags
			},
		}),
		"definition_properties": fabricitem.NewResourceFabricItemDefinitionProperties(fabricitem.ResourceFabricItemDefinitionProperties[itemReadProperties, struct{}]{
			PropertiesAttributes: attributes,
			ItemGetter: func(ctx context.Context, client fabric.Client, model fabricitem.ResourceFabricItemDefinitionPropertiesModel[itemReadProperties, struct{}], to *fabricitem.FabricItemProperties[struct{}]) error {
				return get(ctx, client, model.WorkspaceID.ValueString(), model.ID.ValueString(), to)
			},
			PropertiesSetter: func(context.Context, *struct{}, *fabricitem.ResourceFabricItemDefinitionPropertiesModel[itemReadProperties, struct{}]) diag.Diagnostics {
				return propertyDiags
			},
		}),
		"config_properties": fabricitem.NewResourceFabricItemConfigProperties(fabricitem.ResourceFabricItemConfigProperties[itemReadProperties, struct{}, itemReadProperties, struct{}]{
			PropertiesAttributes: attributes,
			ConfigAttributes:     attributes,
			CreationPayloadSetter: func(context.Context, itemReadProperties) (*struct{}, diag.Diagnostics) {
				return &struct{}{}, nil
			},
			ItemGetter: func(ctx context.Context, client fabric.Client, model fabricitem.ResourceFabricItemConfigPropertiesModel[itemReadProperties, struct{}, itemReadProperties, struct{}], to *fabricitem.FabricItemProperties[struct{}]) error {
				return get(ctx, client, model.WorkspaceID.ValueString(), model.ID.ValueString(), to)
			},
			PropertiesSetter: func(context.Context, *struct{}, *fabricitem.ResourceFabricItemConfigPropertiesModel[itemReadProperties, struct{}, itemReadProperties, struct{}]) diag.Diagnostics {
				return propertyDiags
			},
		}),
		"config_definition_properties": fabricitem.NewResourceFabricItemConfigDefinitionProperties(
			fabricitem.ResourceFabricItemConfigDefinitionProperties[itemReadProperties, struct{}, itemReadProperties, struct{}]{
				PropertiesAttributes: attributes,
				ConfigAttributes:     attributes,
				CreationPayloadSetter: func(context.Context, itemReadProperties) (*struct{}, diag.Diagnostics) {
					return &struct{}{}, nil
				},
				ItemGetter: func(ctx context.Context, client fabric.Client, model fabricitem.ResourceFabricItemConfigDefinitionPropertiesModel[itemReadProperties, struct{}, itemReadProperties, struct{}], to *fabricitem.FabricItemProperties[struct{}]) error {
					return get(ctx, client, model.WorkspaceID.ValueString(), model.ID.ValueString(), to)
				},
				PropertiesSetter: func(context.Context, *struct{}, *fabricitem.ResourceFabricItemConfigDefinitionPropertiesModel[itemReadProperties, struct{}, itemReadProperties, struct{}]) diag.Diagnostics {
					return propertyDiags
				},
			},
		),
	}
}

func itemReadState(t *testing.T, r resource.Resource) tfsdk.State {
	t.Helper()

	ctx := t.Context()
	var resp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)

	objectType, ok := resp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	require.True(t, ok)

	values := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attributeType := range objectType.AttributeTypes {
		values[name] = tftypes.NewValue(attributeType, nil)
	}

	values["id"] = tftypes.NewValue(tftypes.String, "11111111-1111-4111-8111-111111111111")
	values["workspace_id"] = tftypes.NewValue(tftypes.String, "22222222-2222-4222-8222-222222222222")
	values["display_name"] = tftypes.NewValue(tftypes.String, "existing-item")
	values["description"] = tftypes.NewValue(tftypes.String, "last-known description")

	for _, name := range []string{"properties", "configuration"} {
		if attributeType, exists := objectType.AttributeTypes[name]; exists {
			values[name] = tftypes.NewValue(attributeType, map[string]tftypes.Value{
				"value": tftypes.NewValue(tftypes.String, "last-known value"),
			})
		}
	}

	return tfsdk.State{Schema: resp.Schema, Raw: tftypes.NewValue(objectType, values)}
}

func TestUnit_FabricItemResource_ReadErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		status    int
		code      string
		message   string
		wantError bool
		wantNull  bool
	}{
		{name: "inactive_capacity_404", status: http.StatusNotFound, code: "CapacityNotActive", message: "Capacity is not active"},
		{name: "inactive_capacity_400", status: http.StatusBadRequest, code: "CapacityNotActive", message: "Capacity is not active"},
		{name: "inactive_capacity_403", status: http.StatusForbidden, code: "CapacityNotActive", message: "Capacity is not active"},
		{name: "deleted_item", status: http.StatusNotFound, code: "ItemNotFound", message: "Item not found", wantNull: true},
		{name: "message_only_is_not_capacity_code", status: http.StatusNotFound, code: "ItemNotFound", message: "CapacityNotActive", wantNull: true},
		{name: "unrelated_error", status: http.StatusBadRequest, code: "InvalidInput", message: "Invalid input", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for name, r := range itemReadResources(nil) {
				t.Run(name, func(t *testing.T) {
					ctx := t.Context()
					transport := &itemReadTransport{status: tt.status, code: tt.code, message: tt.message}
					client, err := fabric.NewClient(itemReadCredential{}, nil, &fabric.ClientOptions{
						Transport: transport,
						Retry:     policy.RetryOptions{MaxRetries: -1},
					})
					require.NoError(t, err)

					configurable, ok := r.(resource.ResourceWithConfigure)
					require.True(t, ok)
					var configureResp resource.ConfigureResponse
					configurable.Configure(ctx, resource.ConfigureRequest{
						ProviderData: &pconfig.ProviderData{FabricClient: client, Timeout: time.Minute},
					}, &configureResp)
					require.False(t, configureResp.Diagnostics.HasError(), "%v", configureResp.Diagnostics)

					state := itemReadState(t, r)
					resp := resource.ReadResponse{State: state}
					r.Read(ctx, resource.ReadRequest{State: state}, &resp)

					require.Equal(t, 1, transport.calls)
					require.Equal(t, tt.wantError, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
					require.Equal(t, tt.wantNull, resp.State.Raw.IsNull())

					if !tt.wantNull {
						require.True(t, state.Raw.Equal(resp.State.Raw), "previous item state must remain unchanged")
					}

					if tt.code == "CapacityNotActive" {
						require.Len(t, resp.Diagnostics, 1)
						require.Equal(t, diag.SeverityWarning, resp.Diagnostics[0].Severity())
						require.Equal(t, common.WarningCapacityNotActiveHeader, resp.Diagnostics[0].Summary())
						require.Contains(t, resp.Diagnostics[0].Detail(), "11111111-1111-4111-8111-111111111111")
						require.Contains(t, resp.Diagnostics[0].Detail(), "CapacityNotActive")
						require.Contains(t, resp.Diagnostics[0].Detail(), "read-request")

						importable, ok := r.(resource.ResourceWithImportState)
						require.True(t, ok)

						importResp := resource.ImportStateResponse{State: state}
						importable.ImportState(ctx, resource.ImportStateRequest{
							ID: "22222222-2222-4222-8222-222222222222/11111111-1111-4111-8111-111111111111",
						}, &importResp)
						require.True(t, importResp.Diagnostics.HasError(), "import must not succeed with stale state")
						require.Empty(t, importResp.Diagnostics.Warnings())
						require.Contains(t, importResp.Diagnostics.Errors()[0].Detail(), "CapacityNotActive")

						transport.calls = 0
						plan := tfsdk.Plan(state)
						updateResp := resource.UpdateResponse{State: state}
						r.Update(ctx, resource.UpdateRequest{Plan: plan, State: state}, &updateResp)
						require.Equal(t, 1, transport.calls, "update must reach the item read")
						require.True(t, updateResp.Diagnostics.HasError(), "update must not succeed with stale state")
						require.Empty(t, updateResp.Diagnostics.Warnings())
						require.Contains(t, updateResp.Diagnostics.Errors()[0].Detail(), "CapacityNotActive")

						transport.calls = 0
						transport.allowCreate = true
						transport.successfulReads = 1
						createResp := resource.CreateResponse{State: state}
						r.Create(ctx, resource.CreateRequest{
							Plan:   plan,
							Config: tfsdk.Config(state),
						}, &createResp)
						require.Equal(t, 3, transport.calls, "create and tag synchronization must succeed before the item read fails")
						require.True(t, createResp.Diagnostics.HasError(), "create must not succeed with stale state")
						require.Empty(t, createResp.Diagnostics.Warnings())
						require.Contains(t, createResp.Diagnostics.Errors()[0].Detail(), "CapacityNotActive")
					} else if tt.wantNull {
						require.Equal(t, "Resource not found", resp.Diagnostics[0].Summary())
					}
				})
			}
		})
	}
}

func TestUnit_FabricItemResource_ReadResults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		diags diag.Diagnostics
	}{
		{name: "success"},
		{name: "property_error", diags: diag.Diagnostics{diag.NewErrorDiagnostic("Property read failed", "Unable to map item properties")}},
		{name: "property_warning", diags: diag.Diagnostics{diag.NewWarningDiagnostic("Property read warning", "Some item properties are unavailable")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for name, r := range itemReadResources(tt.diags) {
				t.Run(name, func(t *testing.T) {
					ctx := t.Context()
					transport := &itemReadTransport{successfulReads: 1}
					client, err := fabric.NewClient(itemReadCredential{}, nil, &fabric.ClientOptions{
						Transport: transport,
						Retry:     policy.RetryOptions{MaxRetries: -1},
					})
					require.NoError(t, err)

					configurable, ok := r.(resource.ResourceWithConfigure)
					require.True(t, ok)
					var configureResp resource.ConfigureResponse
					configurable.Configure(ctx, resource.ConfigureRequest{
						ProviderData: &pconfig.ProviderData{FabricClient: client, Timeout: time.Minute},
					}, &configureResp)
					require.False(t, configureResp.Diagnostics.HasError(), "%v", configureResp.Diagnostics)

					state := itemReadState(t, r)
					resp := resource.ReadResponse{State: state}
					r.Read(ctx, resource.ReadRequest{State: state}, &resp)

					require.Equal(t, 1, transport.calls)
					require.False(t, resp.State.Raw.IsNull())

					if name == "item" || name == "definition" {
						require.Empty(t, resp.Diagnostics)
					} else {
						require.Equal(t, tt.diags, resp.Diagnostics)
					}

					if resp.Diagnostics.HasError() {
						require.True(t, state.Raw.Equal(resp.State.Raw), "property errors must not overwrite previous state")
					} else {
						var displayName types.String
						require.Empty(t, resp.State.GetAttribute(ctx, path.Root("display_name"), &displayName))
						require.Equal(t, "existing-item", displayName.ValueString())
						require.False(t, state.Raw.Equal(resp.State.Raw), "successful reads must update state from the API")
					}
				})
			}
		})
	}
}
