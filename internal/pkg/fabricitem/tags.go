// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

type DataSourceTagModel struct {
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
}

func SetResourceTagsFromItem(_ context.Context, tags *types.Set, from []fabcore.ItemTag) diag.Diagnostics {
	elements := make([]attr.Value, 0, len(from))

	for _, tag := range from {
		elements = append(elements, customtypes.NewUUIDPointerValue(tag.ID))
	}

	setValue, diags := types.SetValue(customtypes.UUIDType{}, elements)
	if diags.HasError() {
		return diags
	}

	*tags = setValue

	return nil
}

func SetDataSourceTagsFromItem(ctx context.Context, tags *supertypes.SetNestedObjectValueOf[DataSourceTagModel], from []fabcore.ItemTag) diag.Diagnostics {
	if from == nil {
		*tags = supertypes.NewSetNestedObjectValueOfNull[DataSourceTagModel](ctx)

		return nil
	}

	result := make([]*DataSourceTagModel, 0, len(from))

	for _, tag := range from {
		result = append(result, &DataSourceTagModel{
			ID:          customtypes.NewUUIDPointerValue(tag.ID),
			DisplayName: types.StringPointerValue(tag.DisplayName),
		})
	}

	if diags := tags.Set(ctx, result); diags.HasError() {
		return diags
	}

	return nil
}

// SyncTags synchronizes item tags: lists current tags, unapplies all existing, then applies desired ones.
// A null set means "remove all tags" (attribute omitted from config).
func SyncTags(ctx context.Context, tagsClient *fabcore.TagsClient, itemsClient *fabcore.ItemsClient, tags types.Set, workspaceID, itemID string) diag.Diagnostics {
	var desiredTagIDs []string

	if !tags.IsNull() {
		if diags := tags.ElementsAs(ctx, &desiredTagIDs, false); diags.HasError() {
			return diags
		}
	}

	// Get current tags from the item
	respGet, err := itemsClient.GetItem(ctx, workspaceID, itemID, nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	// Unapply all current tags
	if len(respGet.Tags) > 0 {
		currentTagIDs := make([]string, 0, len(respGet.Tags))
		for _, tag := range respGet.Tags {
			if tag.ID != nil {
				currentTagIDs = append(currentTagIDs, *tag.ID)
			}
		}

		if len(currentTagIDs) > 0 {
			_, err := tagsClient.UnapplyTags(ctx, workspaceID, itemID, fabcore.UnapplyTagsRequest{Tags: currentTagIDs}, nil)
			if diags := utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil); diags.HasError() {
				return diags
			}
		}
	}

	// Apply desired tags
	if len(desiredTagIDs) > 0 {
		_, err := tagsClient.ApplyTags(ctx, workspaceID, itemID, fabcore.ApplyTagsRequest{Tags: desiredTagIDs}, nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil); diags.HasError() {
			return diags
		}
	}

	return nil
}
