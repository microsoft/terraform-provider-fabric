// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package folder

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
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
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceFolders)(nil)

type dataSourceFolders struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.FoldersClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceFolders() datasource.DataSource {
	return &dataSourceFolders{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceFolders) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(true)
}

func (d *dataSourceFolders) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	s := itemSchema(true).GetDataSource(ctx)

	resp.Schema = schema.Schema{
		MarkdownDescription: s.GetMarkdownDescription(),
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The Workspace ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"root_folder_id": schema.StringAttribute{
				MarkdownDescription: "This parameter allows users to filter folders based on a specific root folder. If not provided, the workspace is used as the root folder.",
				Optional:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"recursive": schema.BoolAttribute{
				MarkdownDescription: "Lists folders in a folder and its nested folders, or just a folder only. True - All folders in the folder and its nested folders are listed, False - Only folders in the folder are listed. The default value is true.",
				Optional:            true,
			},
			"values": schema.SetNestedAttribute{
				MarkdownDescription: "The set of " + d.TypeInfo.Names + ".",
				Computed:            true,
				CustomType:          supertypes.NewSetNestedObjectTypeOf[baseFolderModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: s.Attributes,
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceFolders) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(d.TypeInfo.Name, d.TypeInfo.IsPreview, d.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}

	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewFoldersClient()
}

func (d *dataSourceFolders) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceFoldersModel

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

func (d *dataSourceFolders) list(ctx context.Context, model *dataSourceFoldersModel) diag.Diagnostics {
	recursive := true
	if !model.Recursive.IsNull() && !model.Recursive.IsUnknown() {
		recursive = model.Recursive.ValueBool()
	}

	options := &fabcore.FoldersClientListFoldersOptions{
		Recursive: to.Ptr(recursive),
	}

	if !model.RootFolderID.IsNull() && !model.RootFolderID.IsUnknown() {
		options.RootFolderID = model.RootFolderID.ValueStringPointer()
	}

	respList, err := d.client.ListFolders(ctx, model.WorkspaceID.ValueString(), options)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	model.setValues(ctx, respList)

	return nil
}
