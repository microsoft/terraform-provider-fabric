package jobscheduler

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
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceJobSchedules)(nil)

type dataSourceJobSchedules struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.JobSchedulerClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceJobSchedules() datasource.DataSource {
	return &dataSourceJobSchedules{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceJobSchedules) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(true)
}

func (d *dataSourceJobSchedules) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	s := itemSchema(true).GetDataSource(ctx)

	resp.Schema = schema.Schema{
		MarkdownDescription: s.GetMarkdownDescription(),
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The Workspace ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"item_id": schema.StringAttribute{
				MarkdownDescription: "The Item ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"job_type": schema.StringAttribute{
				MarkdownDescription: "The type of job to filter by. If not specified, all",
				Required:            true,
			},
			"values": schema.SetNestedAttribute{
				MarkdownDescription: "The set of " + d.TypeInfo.Names + ".",
				Computed:            true,
				CustomType:          supertypes.NewSetNestedObjectTypeOf[baseJobScheduleModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: s.Attributes,
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceJobSchedules) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewJobSchedulerClient()
}

func (d *dataSourceJobSchedules) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceJobSchedulesModel

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

func (d *dataSourceJobSchedules) list(ctx context.Context, model *dataSourceJobSchedulesModel) diag.Diagnostics {
	tflog.Trace(ctx, "LIST", map[string]any{
		"action": "start",
		"model":  model,
	})

	respList, err := d.client.ListItemSchedules(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), model.JobType.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	tflog.Trace(ctx, "LIST", map[string]any{
		"action": "end",
	})

	return model.setValues(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), model.JobType.ValueString(), respList)
}
