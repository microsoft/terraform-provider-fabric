// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func SetTags(ctx context.Context, tags *supertypes.SetValueOf[customtypes.UUID], from []fabcore.ItemTag) diag.Diagnostics {
	elements := make([]customtypes.UUID, 0, len(from))

	for _, tag := range from {
		elements = append(elements, customtypes.NewUUIDPointerValue(tag.ID))
	}

	v := supertypes.NewSetValueOfNull[customtypes.UUID](ctx)

	if diags := v.Set(ctx, elements); diags.HasError() {
		return diags
	}

	*tags = v

	return nil
}

// SyncTags synchronizes item tags by fetching the current applied tags from the API,
// computing the diff against the planned tags, then only unapplying removed tags and applying new ones.
func SyncTags(ctx context.Context, itemsClient *fabcore.ItemsClient, tagsClient *fabcore.TagsClient, plannedTags supertypes.SetValueOf[customtypes.UUID], workspaceID, itemID string) diag.Diagnostics {
	var plannedTagIDs []string

	if !plannedTags.IsNull() {
		if diags := plannedTags.ElementsAs(ctx, &plannedTagIDs, false); diags.HasError() {
			return diags
		}
	}

	// Fetch current tags from the API
	respGet, err := itemsClient.GetItem(ctx, workspaceID, itemID, nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	var currentTagIDs []string

	for _, tag := range respGet.Tags {
		if tag.ID != nil {
			currentTagIDs = append(currentTagIDs, *tag.ID)
		}
	}

	toAdd, toRemove := diffTags(currentTagIDs, plannedTagIDs)

	// Unapply only the tags that were removed
	if len(toRemove) > 0 {
		_, err := tagsClient.UnapplyTags(ctx, workspaceID, itemID, fabcore.UnapplyTagsRequest{Tags: toRemove}, nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil); diags.HasError() {
			return diags
		}
	}

	// Apply only the tags that are newly added
	if len(toAdd) > 0 {
		_, err := tagsClient.ApplyTags(ctx, workspaceID, itemID, fabcore.ApplyTagsRequest{Tags: toAdd}, nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil); diags.HasError() {
			return diags
		}
	}

	return nil
}

// diffTags computes which tags to add and remove by comparing current and planned tag IDs.
func diffTags(current, planned []string) (toAdd, toRemove []string) { //nolint:nonamedreturns
	for _, id := range planned {
		if !slices.Contains(current, id) {
			toAdd = append(toAdd, id)
		}
	}

	for _, id := range current {
		if !slices.Contains(planned, id) {
			toRemove = append(toRemove, id)
		}
	}

	return toAdd, toRemove
}
