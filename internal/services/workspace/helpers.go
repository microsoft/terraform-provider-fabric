// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
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

func getCapacity(ctx context.Context, client *fabcore.CapacitiesClient, capacityID *string, asDiagErr bool) diag.Diagnostics { //revive:disable-line:flag-parameter
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
