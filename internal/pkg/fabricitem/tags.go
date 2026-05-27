// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

// TagsAttrPath is the attribute path of the `tags` set on item resources and
// data sources.
const TagsAttrPath = "tags"

// itemTagDataSourceObjectType is the object type used by the data-source side
// `tags` nested attribute (id + display_name).
func itemTagDataSourceObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":           customtypes.UUIDType{},
			"display_name": types.StringType,
		},
	}
}

// resourceTagsValueFromTagInfos converts a slice of ItemTagInfo into the
// resource-side `Tags` model field type (a typed Set of UUIDs).
func resourceTagsValueFromTagInfos(ctx context.Context, tags []ItemTagInfo) supertypes.SetValueOf[customtypes.UUID] {
	values := make([]customtypes.UUID, 0, len(tags))

	for _, t := range tags {
		values = append(values, customtypes.NewUUIDPointerValue(t.ID))
	}

	return supertypes.NewSetValueOfSlice(ctx, values)
}

// dataSourceTagsValueFromTagInfos converts a slice of ItemTagInfo into a
// types.Set value carrying nested {id, display_name} objects, suitable for the
// data-source side `tags` schema attribute.
func dataSourceTagsValueFromTagInfos(ctx context.Context, tags []ItemTagInfo) (types.Set, diag.Diagnostics) {
	objType := itemTagDataSourceObjectType()

	if len(tags) == 0 {
		return types.SetValue(objType, []attr.Value{})
	}

	values := make([]attr.Value, 0, len(tags))

	for _, t := range tags {
		obj, diags := types.ObjectValue(objType.AttrTypes, map[string]attr.Value{
			"id":           customtypes.NewUUIDPointerValue(t.ID),
			"display_name": types.StringPointerValue(t.DisplayName),
		})
		if diags.HasError() {
			return types.SetNull(objType), diags
		}

		values = append(values, obj)
	}

	return types.SetValue(objType, values)
}

// extractTagIDs reads UUID elements out of a typed Set of customtypes.UUID. Null
// or unknown sets yield an empty slice (used both for plan and state extraction).
func extractTagIDs(ctx context.Context, set supertypes.SetValueOf[customtypes.UUID]) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if set.IsNull() || set.IsUnknown() {
		return []string{}, diags
	}

	uuids := make([]customtypes.UUID, 0, len(set.Elements()))
	if d := set.ElementsAs(ctx, &uuids, false); d.HasError() {
		return nil, d
	}

	out := make([]string, 0, len(uuids))
	for _, u := range uuids {
		out = append(out, u.ValueString())
	}

	return out, diags
}

// diffTagSets returns the IDs that are in `plan` but not in `state` (added) and
// the IDs that are in `state` but not in `plan` (removed).
func diffTagSets(plan, state []string) (added, removed []string) {
	planSet := make(map[string]struct{}, len(plan))
	for _, id := range plan {
		planSet[id] = struct{}{}
	}

	stateSet := make(map[string]struct{}, len(state))
	for _, id := range state {
		stateSet[id] = struct{}{}
	}

	for id := range planSet {
		if _, ok := stateSet[id]; !ok {
			added = append(added, id)
		}
	}

	for id := range stateSet {
		if _, ok := planSet[id]; !ok {
			removed = append(removed, id)
		}
	}

	return added, removed
}

// applyTagDiff issues an unapplyTags call for removed IDs followed by an
// applyTags call for added IDs against the supplied workspace+item. Each side is
// a single batched request — well under the API's 25 req/min/principal cap.
//
// Order is unapply-then-apply to avoid hitting any per-item tag cap mid-update.
func applyTagDiff(
	ctx context.Context,
	tagsClient *fabcore.TagsClient,
	workspaceID, itemID string,
	added, removed []string,
) error {
	if len(removed) > 0 {
		if _, err := tagsClient.UnapplyTags(ctx, workspaceID, itemID, fabcore.UnapplyTagsRequest{Tags: removed}, nil); err != nil {
			return fmt.Errorf("unapplyTags: %w", err)
		}
	}

	if len(added) > 0 {
		if _, err := tagsClient.ApplyTags(ctx, workspaceID, itemID, fabcore.ApplyTagsRequest{Tags: added}, nil); err != nil {
			return fmt.Errorf("applyTags: %w", err)
		}
	}

	return nil
}
