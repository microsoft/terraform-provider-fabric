// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigure = (*dataSourceWorkspaceManagedPrivateEndpoint)(nil)
	// _ datasource.DataSourceWithConfigValidators = (*dataSourceWorkspaceManagedPrivateEndpoint)(nil)
)

type dataSourceWorkspaceManagedPrivateEndpoint struct {
	Name        string
	TFName      string
	DocsURL     string
	IsPreview   bool
	pConfigData *pconfig.ProviderData
	client      *fabcore.ManagedPrivateEndpointsClient
}

func NewDataSourceWorkspaceManagedPrivateEndpoint() datasource.DataSource {
	return &dataSourceWorkspaceManagedPrivateEndpoint{
		Name:      WorkspaceManagedPrivateEndpointName,
		TFName:    WorkspaceManagedPrivateEndpointTFName,
		DocsURL:   WorkspaceManagedPrivateEndpointDocsURL,
		IsPreview: true,
	}
}

func (d *dataSourceWorkspaceManagedPrivateEndpoint) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.TFName
}

func (d *dataSourceWorkspaceManagedPrivateEndpoint) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attributes := getDataSourceWorkspaceManagedPrivateEndpointAttributes(ctx, d.Name)
	attributes["workspace_id"] = schema.StringAttribute{
		MarkdownDescription: "The Workspace ID.",
		Required:            true,
		CustomType:          customtypes.UUIDType{},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: fabricitem.GetDataSourcePreviewNote("Get a Fabric "+d.Name+".\n\n"+
			"Use this data source to fetch a ["+d.Name+"]("+d.DocsURL+").\n\n"+
			ItemDocsSPNSupport, d.IsPreview),
		Attributes: attributes,
	}
}

func (d *dataSourceWorkspaceManagedPrivateEndpoint) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	pConfigData, ok := req.ProviderData.(*pconfig.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			common.ErrorDataSourceConfigType,
			fmt.Sprintf(common.ErrorFabricClientType, req.ProviderData),
		)

		return
	}

	d.pConfigData = pConfigData
	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewManagedPrivateEndpointsClient()

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(d.Name, d.IsPreview, d.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}
}

func (d *dataSourceWorkspaceManagedPrivateEndpoint) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceWorkspaceManagedPrivateEndpointModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, d.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(d.get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dataSourceWorkspaceManagedPrivateEndpoint) get(ctx context.Context, model *dataSourceWorkspaceManagedPrivateEndpointModel) diag.Diagnostics {
	tflog.Trace(ctx, "LIST", map[string]any{
		"action": "start",
		"model":  model,
	})

	respGet, err := d.client.GetWorkspaceManagedPrivateEndpoint(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	tflog.Trace(ctx, "LIST", map[string]any{
		"action": "end",
	})

	return model.set(ctx, respGet.ManagedPrivateEndpoint)
}
