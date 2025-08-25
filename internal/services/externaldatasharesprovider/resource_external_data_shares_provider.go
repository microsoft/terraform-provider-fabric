package externaldatasharesprovider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure = (*resourceExternalDataSharesProvider)(nil)
)

type resourceExternalDataSharesProvider struct {
	pConfigData *pconfig.ProviderData
	client      *fabadmin.ExternalDataSharesProviderClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceExternalDataSharesProvider() resource.Resource {
	return &resourceExternalDataSharesProvider{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceExternalDataSharesProvider) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceExternalDataSharesProvider) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema().GetResource(ctx)
}

func (r *resourceExternalDataSharesProvider) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	pConfigData, ok := req.ProviderData.(*pconfig.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			common.ErrorResourceConfigType,
			fmt.Sprintf(common.ErrorFabricClientType, req.ProviderData),
		)

		return
	}

	r.pConfigData = pConfigData

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(r.TypeInfo.Name, r.TypeInfo.IsPreview, r.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}

	r.client = fabadmin.NewClientFactoryWithClient(*pConfigData.FabricClient).NewExternalDataSharesProviderClient()
}

func (r *resourceExternalDataSharesProvider) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddWarning(
		"delete operation not supported",
		fmt.Sprintf(
			"Resource %s does not support creation. It will be removed from Terraform state, but no action will be taken in the Fabric. All current settings will remain.",
			r.TypeInfo.Names,
		),
	)
	resp.State.RemoveResource(ctx)
}

func (r *resourceExternalDataSharesProvider) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	resp.Diagnostics.AddWarning(
		"read operation not supported",
		fmt.Sprintf(
			"Resource %s does not support Read. It will be removed from Terraform state, but no action will be taken in the Fabric. All current settings will remain.",
			r.TypeInfo.Names,
		),
	)
	resp.State.RemoveResource(ctx)
}

func (r *resourceExternalDataSharesProvider) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"update operation not supported",
		fmt.Sprintf(
			"Resource %s does not support Update. It will be removed from Terraform state, but no action will be taken in the Fabric. All current settings will remain.",
			r.TypeInfo.Names,
		),
	)
	resp.State.RemoveResource(ctx)
}

func (r *resourceExternalDataSharesProvider) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	var state resourceExternalDataSharesProviderModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.RevokeExternalDataShare(ctx, state.WorkspaceID.ValueString(), state.ItemID.ValueString(), state.ExternalDataShareID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}
