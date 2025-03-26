// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/fabric-sdk-go/fabric"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigure = (*DataSourceFabricItemsProperties[struct{}, struct{}])(nil)
)

type DataSourceFabricItemsProperties[Ttfprop, Titemprop any] struct {
	DataSourceFabricItems
	PropertiesAttributes map[string]schema.Attribute
	PropertiesSetter     func(ctx context.Context, from *Titemprop, to *FabricItemPropertiesModel[Ttfprop, Titemprop]) diag.Diagnostics
	ItemListGetter       func(ctx context.Context, fabricClient fabric.Client, model DataSourceFabricItemsPropertiesModel[Ttfprop, Titemprop], fabricItems *[]FabricItemProperties[Titemprop]) error
}

func NewDataSourceFabricItemsProperties[Ttfprop, Titemprop any](config DataSourceFabricItemsProperties[Ttfprop, Titemprop]) datasource.DataSource {
	return &config
}

func (d *DataSourceFabricItemsProperties[Ttfprop, Titemprop]) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.TFName
}

func (d *DataSourceFabricItemsProperties[Ttfprop, Titemprop]) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attributes := map[string]schema.Attribute{
		"workspace_id": schema.StringAttribute{
			MarkdownDescription: "The Workspace ID.",
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"id": schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s ID.", d.Name),
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s display name.", d.Name),
			Computed:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s description.", d.Name),
			Computed:            true,
		},
	}

	attributes["properties"] = getDataSourceFabricItemPropertiesNestedAttr[Ttfprop](ctx, d.Name, d.PropertiesAttributes)

	resp.Schema = schema.Schema{
		MarkdownDescription: d.MarkdownDescription,
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The Workspace ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"values": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: fmt.Sprintf("The list of %s.", d.Names),
				CustomType:          supertypes.NewSetNestedObjectTypeOf[FabricItemPropertiesModel[Ttfprop, Titemprop]](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: attributes,
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *DataSourceFabricItemsProperties[Ttfprop, Titemprop]) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	if resp.Diagnostics.Append(IsPreviewMode(d.Name, d.IsPreview, d.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}
}

func (d *DataSourceFabricItemsProperties[Ttfprop, Titemprop]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data DataSourceFabricItemsPropertiesModel[Ttfprop, Titemprop]

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

func (d *DataSourceFabricItemsProperties[Ttfprop, Titemprop]) list(ctx context.Context, model *DataSourceFabricItemsPropertiesModel[Ttfprop, Titemprop]) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %ss", d.Name))

	var fabricItems []FabricItemProperties[Titemprop]

	err := d.ItemListGetter(ctx, *d.pConfigData.FabricClient, *model, &fabricItems)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	return model.setValues(ctx, fabricItems, d.PropertiesSetter)
}
