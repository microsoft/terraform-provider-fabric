// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure = (*resourceTag)(nil)
)

type resourceTag struct {
	pConfigData *pconfig.ProviderData
	client      *fabadmin.TagsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceTag() resource.Resource {
	return &resourceTag{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceTag) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceTag) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceItemSchema().GetResource(ctx)
}

func (r *resourceTag) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = fabadmin.NewClientFactoryWithClient(*pConfigData.FabricClient).NewTagsClient()
}

func (r *resourceTag) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan, state resourceTagsModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	state.Timeouts = plan.Timeouts

	var reqCreate requestCreateTags

	if resp.Diagnostics.Append(reqCreate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	respCreate, err := r.client.BulkCreateTags(ctx, reqCreate.CreateTagsRequest, nil)

	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(state.set(ctx, respCreate.Tags)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceTag) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceTagsModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Read(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	diags = r.get(ctx, &state)
	if utils.IsErrNotFound(state.ID.ValueString(), &diags, fabcore.ErrCommon.EntityNotFound) {
		resp.State.RemoveResource(ctx)

		resp.Diagnostics.Append(diags...)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceTag) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan resourceTagsModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateTags

	reqUpdate.set(plan)

	respUpdate, err := r.client.UpdateTag(ctx, plan.ID.ValueString(), reqUpdate.UpdateTagRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.setValue(ctx, respUpdate.Tag)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceTag) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	var state resourceTagsModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if state.ID.IsNull() {
		tags, diags := state.Tags.Get(ctx)
		if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
			return
		}

		for _, tag := range tags {
			_, err := r.client.DeleteTag(ctx, tag.ID.ValueString(), nil)
			if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
				return
			}
		}
	} else {
		_, err := r.client.DeleteTag(ctx, state.ID.ValueString(), nil)
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
			return
		}
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceTag) get(ctx context.Context, model *resourceTagsModel) diag.Diagnostics {
	tagsRequested, diags := model.Tags.Get(ctx)
	if diags.HasError() {
		return diags
	}

	var found []*baseTagModel

	pager := r.client.NewListTagsPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if d := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); d.HasError() {
			return d
		}

		for _, entity := range page.Value {
			for _, requestedTag := range tagsRequested {
				if *entity.DisplayName == requestedTag.DisplayName.ValueString() {
					item := &baseTagModel{}
					item.set(ctx, entity)

					found = append(found, item)
				}
			}
		}
	}

	if diags := model.Tags.Set(ctx, found); diags.HasError() {
		return diags
	}

	return nil
}
