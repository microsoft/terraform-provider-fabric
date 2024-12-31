// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceDomainWorkspaceAssignments)(nil)

type dataSourceDomainWorkspaceAssignments struct {
	pConfigData *pconfig.ProviderData
	client      *fabadmin.DomainsClient
}

func NewDataSourceDomainWorkspaceAssignments() datasource.DataSource {
	return &dataSourceDomainWorkspaceAssignments{}
}

func (d *dataSourceDomainWorkspaceAssignments) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DomainWorkspaceAssignmentsTFName
}

func (d *dataSourceDomainWorkspaceAssignments) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List a Fabric " + DomainWorkspaceAssignmentsName + ".\n\n" +
			"See [" + ItemName + "](" + ItemDocsURL + ") for more information.\n\n" +
			ItemDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.StringAttribute{
				MarkdownDescription: "The Domain ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"values": schema.ListNestedAttribute{
				MarkdownDescription: "The list of " + DomainWorkspaceAssignmentsName + ".",
				Computed:            true,
				CustomType:          supertypes.NewListNestedObjectTypeOf[workspaceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The Workspace ID.",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The Workspace display name.",
							Computed:            true,
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceDomainWorkspaceAssignments) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabadmin.NewClientFactoryWithClient(*pConfigData.FabricClient).NewDomainsClient()
}

func (d *dataSourceDomainWorkspaceAssignments) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceDomainWorkspaceAssignmentsModel

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

func (d *dataSourceDomainWorkspaceAssignments) list(ctx context.Context, model *dataSourceDomainWorkspaceAssignmentsModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting "+DomainWorkspaceAssignmentsName)

	respList, err := d.client.ListDomainWorkspaces(ctx, model.DomainID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	return model.setValues(ctx, respList)
}
