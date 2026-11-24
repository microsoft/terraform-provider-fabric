// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func checkWorkspaceType(entity fabcore.WorkspaceInfo) diag.Diagnostics {
	var diags diag.Diagnostics

	switch *entity.Type {
	case fabcore.WorkspaceTypePersonal:
		diags.AddError(
			common.ErrorWorkspaceNotSupportedHeader,
			fmt.Sprintf(common.ErrorWorkspaceNotSupportedDetails, string(fabcore.WorkspaceTypePersonal)),
		)

		return diags
	case fabcore.WorkspaceTypeAdminWorkspace:
		diags.AddError(
			common.ErrorWorkspaceNotSupportedHeader,
			fmt.Sprintf(common.ErrorWorkspaceNotSupportedDetails, string(fabcore.WorkspaceTypeAdminWorkspace)),
		)

		return diags
	default:
		return nil
	}
}

func validateCapacityState(ctx context.Context, client *fabcore.CapacitiesClient, capacityID *string) diag.Diagnostics {
	if client == nil || capacityID == nil {
		return nil
	}

	var diags diag.Diagnostics

	pager := client.NewListCapacitiesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			if strings.EqualFold(*entity.ID, *capacityID) {
				if *entity.State != fabcore.CapacityStateActive {
					diags.AddError(
						"Fabric Capacity State",
						"Fabric Capacity is NOT in Active state. Inactive Capacity may cause unrecoverable damage. Please ensure the Capacity is in Active state before continuing.",
					)

					return diags
				}

				return nil
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		"Unable to find Capacity with 'id': "+*capacityID,
	)

	return diags
}

// syncTags synchronizes workspace tags by fetching the current applied tags from the API,
// computing the diff against the planned tags, then only unapplying removed tags and applying new ones.
func syncTags(ctx context.Context, client *fabcore.WorkspacesClient, plannedTags supertypes.SetValueOf[customtypes.UUID], workspaceID string) diag.Diagnostics {
	var plannedTagIDs []string

	if !plannedTags.IsNull() && !plannedTags.IsUnknown() {
		if diags := plannedTags.ElementsAs(ctx, &plannedTagIDs, false); diags.HasError() {
			return diags
		}
	}

	respGet, err := client.GetWorkspace(ctx, workspaceID, nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	currentTagIDs := make([]string, 0, len(respGet.Tags))

	for _, tag := range respGet.Tags {
		if tag.ID != nil {
			currentTagIDs = append(currentTagIDs, *tag.ID)
		}
	}

	toAdd, toRemove := fabricitem.DiffTags(currentTagIDs, plannedTagIDs)

	// Unapply only the tags that were removed
	if len(toRemove) > 0 {
		_, err := client.UnapplyWorkspaceTags(ctx, workspaceID, fabcore.UnapplyWorkspaceTagsRequest{Tags: toRemove}, nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil); diags.HasError() {
			return diags
		}
	}

	// Apply only the tags that are newly added
	if len(toAdd) > 0 {
		_, err := client.ApplyWorkspaceTags(ctx, workspaceID, fabcore.ApplyWorkspaceTagsRequest{Tags: toAdd}, nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil); diags.HasError() {
			return diags
		}
	}

	return nil
}
