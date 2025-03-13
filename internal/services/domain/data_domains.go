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
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceDomains)(nil)

type dataSourceDomains struct {
	pConfigData         *pconfig.ProviderData
	client              *fabadmin.DomainsClient
	Name                string
	TFName              string
	MarkdownDescription string
	IsPreview           bool
}

func NewDataSourceDomains() datasource.DataSource {
	markdownDescription := "List a Fabric " + ItemsName + ".\n\n" +
		"Use this data source to list [" + ItemsName + "](" + ItemDocsURL + ").\n\n" +
		ItemDocsSPNSupport

	return &dataSourceDomains{
		Name:                ItemsName,
		TFName:              ItemsTFName,
		MarkdownDescription: fabricitem.GetDataSourcePreviewNote(markdownDescription, ItemPreview),
		IsPreview:           ItemPreview,
	}
}

func (d *dataSourceDomains) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.TFName
}

func (d *dataSourceDomains) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List a Fabric " + d.Name + ".\n\n" +
			"Use this data source to list [" + d.Name + "](" + ItemDocsURL + ").\n\n" +
			ItemDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"values": schema.ListNestedAttribute{
				MarkdownDescription: "The list of " + d.Name + ".",
				Computed:            true,
				CustomType:          supertypes.NewListNestedObjectTypeOf[baseDomainModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The " + ItemName + " ID.",
							Computed:            true,
							CustomType:          customtypes.UUIDType{},
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The " + ItemName + " display name.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The " + ItemName + " description.",
							Computed:            true,
						},
						"parent_domain_id": schema.StringAttribute{
							MarkdownDescription: "The " + ItemName + " parent ID.",
							Computed:            true,
							CustomType:          customtypes.UUIDType{},
						},
						"contributors_scope": schema.StringAttribute{
							MarkdownDescription: "The " + ItemName + " contributors scope. Possible values: " + utils.ConvertStringSlicesToString(fabadmin.PossibleContributorsScopeTypeValues(), true, true) + ".",
							Computed:            true,
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceDomains) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(d.Name, d.IsPreview, d.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}
}

func (d *dataSourceDomains) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceDomainsModel

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

func (d *dataSourceDomains) list(ctx context.Context, model *dataSourceDomainsModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting "+ItemsName)

	respList, err := d.client.ListDomains(ctx, nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	return model.setValues(ctx, respList.Domains)
}
