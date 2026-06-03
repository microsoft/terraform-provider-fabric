// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"

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

// SyncTags synchronizes item tags: unapplies current tags, then applies desired ones.
// A null or empty plannedTags means "remove all tags". CurrentTags represents the known state tags.
func SyncTags(ctx context.Context, tagsClient *fabcore.TagsClient, plannedTags, currentTags supertypes.SetValueOf[customtypes.UUID], workspaceID, itemID string) diag.Diagnostics {
	var desiredTagIDs []string

	if !plannedTags.IsNull() {
		if diags := plannedTags.ElementsAs(ctx, &desiredTagIDs, false); diags.HasError() {
			return diags
		}
	}

	var currentTagIDs []string

	if !currentTags.IsNull() {
		if diags := currentTags.ElementsAs(ctx, &currentTagIDs, false); diags.HasError() {
			return diags
		}
	}

	// Unapply current tags
	if len(currentTagIDs) > 0 {
		_, err := tagsClient.UnapplyTags(ctx, workspaceID, itemID, fabcore.UnapplyTagsRequest{Tags: currentTagIDs}, nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil); diags.HasError() {
			return diags
		}
	}

	// Apply desired tags
	if len(desiredTagIDs) > 0 {
		_, err := tagsClient.ApplyTags(ctx, workspaceID, itemID, fabcore.ApplyTagsRequest{Tags: desiredTagIDs}, nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil); diags.HasError() {
			return diags
		}
	}

	return nil
}
