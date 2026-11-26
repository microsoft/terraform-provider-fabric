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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

type itemReadErrorTransport struct {
	code  string
	calls int
}

func (tr *itemReadErrorTransport) Do(req *http.Request) (*http.Response, error) {
	tr.calls++

	return &http.Response{
		StatusCode: http.StatusNotFound,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body: io.NopCloser(strings.NewReader(fmt.Sprintf(
			`{"requestId":"read-request","errorCode":%q,"message":"Fabric item read failed"}`,
			tr.code,
		))),
		Request: req,
	}, nil
}

type itemReadErrorCredential struct{}

func (itemReadErrorCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "test-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

type itemReadErrorProperties struct {
	Value types.String `tfsdk:"value"`
}

func itemReadErrorResources() map[string]resource.Resource {
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
		"properties": fabricitem.NewResourceFabricItemProperties(fabricitem.ResourceFabricItemProperties[itemReadErrorProperties, struct{}]{
			PropertiesAttributes: attributes,
			ItemGetter: func(ctx context.Context, client fabric.Client, model fabricitem.ResourceFabricItemPropertiesModel[itemReadErrorProperties, struct{}], to *fabricitem.FabricItemProperties[struct{}]) error {
				return get(ctx, client, model.WorkspaceID.ValueString(), model.ID.ValueString(), to)
			},
			PropertiesSetter: func(context.Context, *struct{}, *fabricitem.ResourceFabricItemPropertiesModel[itemReadErrorProperties, struct{}]) diag.Diagnostics {
				return nil
			},
		}),
		"definition_properties": fabricitem.NewResourceFabricItemDefinitionProperties(fabricitem.ResourceFabricItemDefinitionProperties[itemReadErrorProperties, struct{}]{
			PropertiesAttributes: attributes,
			ItemGetter: func(ctx context.Context, client fabric.Client, model fabricitem.ResourceFabricItemDefinitionPropertiesModel[itemReadErrorProperties, struct{}], to *fabricitem.FabricItemProperties[struct{}]) error {
				return get(ctx, client, model.WorkspaceID.ValueString(), model.ID.ValueString(), to)
			},
			PropertiesSetter: func(context.Context, *struct{}, *fabricitem.ResourceFabricItemDefinitionPropertiesModel[itemReadErrorProperties, struct{}]) diag.Diagnostics {
				return nil
			},
		}),
		"config_properties": fabricitem.NewResourceFabricItemConfigProperties(fabricitem.ResourceFabricItemConfigProperties[itemReadErrorProperties, struct{}, itemReadErrorProperties, struct{}]{
			PropertiesAttributes: attributes,
			ConfigAttributes:     attributes,
			CreationPayloadSetter: func(context.Context, itemReadErrorProperties) (*struct{}, diag.Diagnostics) {
				return &struct{}{}, nil
			},
			ItemGetter: func(ctx context.Context, client fabric.Client, model fabricitem.ResourceFabricItemConfigPropertiesModel[itemReadErrorProperties, struct{}, itemReadErrorProperties, struct{}], to *fabricitem.FabricItemProperties[struct{}]) error {
				return get(ctx, client, model.WorkspaceID.ValueString(), model.ID.ValueString(), to)
			},
			PropertiesSetter: func(context.Context, *struct{}, *fabricitem.ResourceFabricItemConfigPropertiesModel[itemReadErrorProperties, struct{}, itemReadErrorProperties, struct{}]) diag.Diagnostics {
				return nil
			},
		}),
		"config_definition_properties": fabricitem.NewResourceFabricItemConfigDefinitionProperties(
			fabricitem.ResourceFabricItemConfigDefinitionProperties[itemReadErrorProperties, struct{}, itemReadErrorProperties, struct{}]{
				PropertiesAttributes: attributes,
				ConfigAttributes:     attributes,
				CreationPayloadSetter: func(context.Context, itemReadErrorProperties) (*struct{}, diag.Diagnostics) {
					return &struct{}{}, nil
				},
				ItemGetter: func(ctx context.Context, client fabric.Client, model fabricitem.ResourceFabricItemConfigDefinitionPropertiesModel[itemReadErrorProperties, struct{}, itemReadErrorProperties, struct{}], to *fabricitem.FabricItemProperties[struct{}]) error {
					return get(ctx, client, model.WorkspaceID.ValueString(), model.ID.ValueString(), to)
				},
				PropertiesSetter: func(
					context.Context,
					*struct{},
					*fabricitem.ResourceFabricItemConfigDefinitionPropertiesModel[itemReadErrorProperties, struct{}, itemReadErrorProperties, struct{}],
				) diag.Diagnostics {
					return nil
				},
			},
		),
	}
}

func itemReadErrorState(t *testing.T, r resource.Resource) tfsdk.State {
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
		code      string
		wantError bool
		wantNull  bool
	}{
		{name: "capacity_not_active", code: "CapacityNotActive", wantError: true},
		{name: "ordinary_not_found", code: "ItemNotFound", wantNull: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for name, r := range itemReadErrorResources() {
				t.Run(name, func(t *testing.T) {
					ctx := t.Context()
					transport := &itemReadErrorTransport{code: tt.code}
					client, err := fabric.NewClient(itemReadErrorCredential{}, nil, &fabric.ClientOptions{
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

					state := itemReadErrorState(t, r)
					resp := resource.ReadResponse{State: state}
					r.Read(ctx, resource.ReadRequest{State: state}, &resp)

					require.Equal(t, 1, transport.calls)
					require.Equal(t, tt.wantError, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
					require.Equal(t, tt.wantNull, resp.State.Raw.IsNull())

					if tt.wantError {
						require.Empty(t, resp.Diagnostics.Warnings())
						require.Len(t, resp.Diagnostics.Errors(), 1)
						require.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "CapacityNotActive")
						require.True(t, state.Raw.Equal(resp.State.Raw), "failed refresh should not report the resource as deleted")
					} else {
						require.Len(t, resp.Diagnostics.Warnings(), 1)
						require.Equal(t, "Resource not found", resp.Diagnostics.Warnings()[0].Summary())
					}
				})
			}
		})
	}
}
