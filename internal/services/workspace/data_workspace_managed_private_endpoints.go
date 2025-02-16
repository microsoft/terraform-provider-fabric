// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigure = (*dataSourceWorkspaceManagedPrivateEndpoints)(nil)
)

type dataSourceWorkspaceManagedPrivateEndpoints struct {
	Name        string
	Names       string
	TFName      string
	DocsURL     string
	IsPreview   bool
	pConfigData *pconfig.ProviderData
	client      *fabcore.ManagedPrivateEndpointsClient
}

func NewDataSourceWorkspaceManagedPrivateEndpoints() datasource.DataSource {
	return &dataSourceWorkspaceManagedPrivateEndpoints{
		Name:      WorkspaceManagedPrivateEndpointName,
		Names:     WorkspaceManagedPrivateEndpointsName,
		TFName:    WorkspaceManagedPrivateEndpointsTFName,
		DocsURL:   WorkspaceManagedPrivateEndpointDocsURL,
		IsPreview: true,
	}
}

func (d *dataSourceWorkspaceManagedPrivateEndpoints) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.TFName
}

func (d *dataSourceWorkspaceManagedPrivateEndpoints) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fabricitem.GetDataSourcePreviewNote("List a Fabric "+d.Names+".\n\n"+
			"Use this data source to list ["+d.Names+"]("+d.DocsURL+") for more information.\n\n"+
			ItemDocsSPNSupport, d.IsPreview),
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The Workspace ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"values": schema.ListNestedAttribute{
				MarkdownDescription: fmt.Sprintf("The list of %s.", d.Names),
				Computed:            true,
				CustomType:          supertypes.NewListNestedObjectTypeOf[baseWorkspaceManagedPrivateEndpointModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: getDataSourceWorkspaceManagedPrivateEndpointAttributes(ctx, d.Name),
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceWorkspaceManagedPrivateEndpoints) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dataSourceWorkspaceManagedPrivateEndpoints) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceWorkspaceManagedPrivateEndpointsModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, d.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(d.list(ctx, &data)...); resp.Diagnostics.HasError() {
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

func (d *dataSourceWorkspaceManagedPrivateEndpoints) list(ctx context.Context, model *dataSourceWorkspaceManagedPrivateEndpointsModel) diag.Diagnostics {
	tflog.Trace(ctx, "LIST", map[string]any{
		"action": "start",
		"model":  model,
	})

	respList, err := d.client.ListWorkspaceManagedPrivateEndpoints(ctx, model.WorkspaceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	tflog.Trace(ctx, "LIST", map[string]any{
		"action": "end",
	})

	return model.setValues(ctx, respList)
}
