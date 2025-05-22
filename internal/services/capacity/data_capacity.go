// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package capacity

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
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
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceCapacity() datasource.DataSource {
	return &dataSourceCapacity{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceCapacity) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

func (d *dataSourceCapacity) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = itemSchema(false).GetDataSource(ctx)
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

	var data dataSourceCapacityModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	// Set a default timeout in case pConfigData is nil
	defaultTimeout := 5 * time.Minute
	timeoutValue := defaultTimeout
	
	// Only use pConfigData.Timeout if pConfigData is not nil
	if d.pConfigData != nil {
		timeoutValue = d.pConfigData.Timeout
	}
	
	timeout, diags := data.Timeouts.Read(ctx, timeoutValue)
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
	return d.get(ctx, true, model)
}

func (d *dataSourceCapacity) getByDisplayName(ctx context.Context, model *dataSourceCapacityModel) diag.Diagnostics {
	return d.get(ctx, false, model)
}

func (d *dataSourceCapacity) get(ctx context.Context, byID bool, model *dataSourceCapacityModel) diag.Diagnostics {
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
