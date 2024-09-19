// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package capacity

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var (
	_ datasource.DataSourceWithConfigValidators = (*dataSourceCapacity)(nil)
	_ datasource.DataSourceWithConfigure        = (*dataSourceCapacity)(nil)
)

type dataSourceCapacity struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.CapacitiesClient
}

func NewDataSourceCapacity() datasource.DataSource {
	return &dataSourceCapacity{}
}

func (d *dataSourceCapacity) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemTFName
}

func (d *dataSourceCapacity) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get a Fabric " + ItemName + ".\n\n" +
			"Use this data source to fetch [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
			ItemDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("The %s ID.", ItemName),
				Optional:            true,
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("The %s display name.", ItemName),
				Optional:            true,
				Computed:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("The Azure region where the %s has been provisioned.", ItemName),
				Computed:            true,
			},
			"sku": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("The %s SKU.", ItemName),
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("The %s state.", ItemName),
				Computed:            true,
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceCapacity) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("id"),
			path.MatchRoot("display_name"),
		),
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("display_name"),
		),
	}
}

// Configure adds the provider configured client to the data source.
func (d *dataSourceCapacity) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewCapacitiesClient()
}

// Read refreshes the Terraform state with the latest data.
func (d *dataSourceCapacity) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceCapacityModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, d.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if data.ID.ValueString() != "" {
		diags = d.getByID(ctx, &data)
	} else {
		diags = d.getByDisplayName(ctx, &data)
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
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

func (d *dataSourceCapacity) getByID(ctx context.Context, model *dataSourceCapacityModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s by 'id'", ItemName))

	return d.get(ctx, true, model)
}

func (d *dataSourceCapacity) getByDisplayName(ctx context.Context, model *dataSourceCapacityModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s by 'display_name'", ItemName))

	return d.get(ctx, false, model)
}

func (d *dataSourceCapacity) get(ctx context.Context, byID bool, model *dataSourceCapacityModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting "+ItemName)

	var diags diag.Diagnostics
	var notFound string

	pager := d.client.NewListCapacitiesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			switch byID {
			case true:
				if *entity.ID == model.ID.ValueString() {
					model.set(entity)

					return nil
				}

				notFound = "Unable to find Capacity with 'id': " + model.ID.ValueString()
			default:
				if *entity.DisplayName == model.DisplayName.ValueString() {
					model.set(entity)

					return nil
				}

				notFound = "Unable to find Capacity with 'display_name': " + model.DisplayName.ValueString()
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		notFound,
	)

	return diags
}
