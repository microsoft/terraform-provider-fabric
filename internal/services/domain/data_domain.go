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

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceDomain)(nil)

type dataSourceDomain struct {
	pConfigData         *pconfig.ProviderData
	client              *fabadmin.DomainsClient
	Name                string
	TFName              string
	MarkdownDescription string
	IsPreview           bool
}

func NewDataSourceDomain() datasource.DataSource {
	markdownDescription := "Get a Fabric " + ItemName + ".\n\n" +
		"Use this data source to get [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
		ItemDocsSPNSupport

	return &dataSourceDomain{
		Name:                ItemName,
		TFName:              ItemTFName,
		MarkdownDescription: fabricitem.GetDataSourcePreviewNote(markdownDescription, ItemPreview),
		IsPreview:           ItemPreview,
	}
}

func (d *dataSourceDomain) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.TFName
}

func (d *dataSourceDomain) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: d.MarkdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The " + d.Name + " ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The " + d.Name + " display name.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The " + d.Name + " description.",
				Computed:            true,
			},
			"parent_domain_id": schema.StringAttribute{
				MarkdownDescription: "The " + d.Name + " parent ID.",
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"contributors_scope": schema.StringAttribute{
				MarkdownDescription: "The " + d.Name + " contributors scope. Possible values: " + utils.ConvertStringSlicesToString(fabadmin.PossibleContributorsScopeTypeValues(), true, true) + ".",
				Computed:            true,
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceDomain) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	diags := fabricitem.IsPreviewMode(d.Name, d.IsPreview, d.pConfigData.Preview)
	if diags != nil {
		resp.Diagnostics.Append(diags...)

		if diags.HasError() {
			return
		}
	}
}

func (d *dataSourceDomain) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceDomainModel

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

func (d *dataSourceDomain) get(ctx context.Context, model *dataSourceDomainModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting "+ItemName)

	respGet, err := d.client.GetDomain(ctx, model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	model.set(respGet.Domain)

	return nil
}
