// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"fmt"

	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigValidators = (*DataSourceFabricItemDefinition)(nil)
	_ datasource.DataSourceWithConfigure        = (*DataSourceFabricItemDefinition)(nil)
)

type DataSourceFabricItemDefinition struct {
	pConfigData         *pconfig.ProviderData
	client              *fabcore.ItemsClient
	Type                fabcore.ItemType
	Name                string
	TFName              string
	MarkdownDescription string
	IsDisplayNameUnique bool
	FormatTypeDefault   string
	FormatTypes         []string
	DefinitionPathKeys  []string
	IsPreview           bool
}

func NewDataSourceFabricItemDefinition(config DataSourceFabricItemDefinition) datasource.DataSource {
	return &config
}

func (d *DataSourceFabricItemDefinition) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.TFName
}

func (d *DataSourceFabricItemDefinition) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getDataSourceFabricItemDefinitionSchema(ctx, *d)
}

func (d *DataSourceFabricItemDefinition) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	if d.IsDisplayNameUnique {
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

	return []datasource.ConfigValidator{}
}

func (d *DataSourceFabricItemDefinition) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewItemsClient()

	diags := IsPreviewModeEnabled(d.Name, d.IsPreview, d.pConfigData.Preview)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)

		return
	} else if diags != nil {
		resp.Diagnostics.Append(diags...)
	}
}

func (d *DataSourceFabricItemDefinition) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceFabricItemDefinitionModel

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

	data.Format = types.StringNull()

	if d.FormatTypeDefault != "" {
		data.Format = types.StringValue(d.FormatTypeDefault)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	if data.OutputDefinition.IsNull() || data.OutputDefinition.IsUnknown() {
		data.OutputDefinition = types.BoolValue(false)
	}

	if data.OutputDefinition.ValueBool() {
		if resp.Diagnostics.Append(d.getDefinition(ctx, &data)...); resp.Diagnostics.HasError() {
			return
		}

		tflog.Debug(ctx, "Definition parts content is gzip base64. Use `provider::fabric::content_decode` function to decode content.")

		resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
	}

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *DataSourceFabricItemDefinition) getByID(ctx context.Context, model *dataSourceFabricItemDefinitionModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s by ID: %s", d.Name, model.ID.ValueString()))

	respGet, err := d.client.GetItem(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	model.set(respGet.Item)

	return nil
}

func (d *DataSourceFabricItemDefinition) getByDisplayName(ctx context.Context, model *dataSourceFabricItemDefinitionModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s by Display Name: %s", d.Name, model.DisplayName.ValueString()))

	var diags diag.Diagnostics

	opts := &fabcore.ItemsClientListItemsOptions{
		Type: azto.Ptr(string(d.Type)),
	}

	pager := d.client.NewListItemsPager(model.WorkspaceID.ValueString(), opts)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			if *entity.DisplayName == model.DisplayName.ValueString() {
				model.set(entity)

				return nil
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		fmt.Sprintf("Unable to find %s with 'display_name': %s in the Workspace ID: %s", d.Name, model.DisplayName.ValueString(), model.WorkspaceID.ValueString()),
	)

	return diags
}

func (d *DataSourceFabricItemDefinition) getDefinition(ctx context.Context, model *dataSourceFabricItemDefinitionModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s definition (WorkspaceID: %s ItemID: %s)", d.Name, model.WorkspaceID.ValueString(), model.ID.ValueString()))

	respGetOpts := &fabcore.ItemsClientBeginGetItemDefinitionOptions{}

	if !model.Format.IsNull() {
		respGetOpts.Format = model.Format.ValueStringPointer()
	}

	respGet, err := d.client.GetItemDefinition(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), respGetOpts)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	definition, diags := getDataSourceDefinitionModel(ctx, *respGet.Definition)
	if diags.HasError() {
		return diags
	}

	model.setDefinition(definition)

	return nil
}
