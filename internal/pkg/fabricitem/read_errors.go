// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

type itemReadResult struct {
	Diagnostics       diag.Diagnostics
	CapacityNotActive bool
}

func getItemReadResult(ctx context.Context, err error) itemReadResult {
	var azureErr *azcore.ResponseError
	if errors.As(err, &azureErr) && azureErr.RawResponse != nil {
		err = fabcore.NewResponseError(azureErr.RawResponse)
	}

	var fabricErr *fabcore.ResponseError
	if errors.As(err, &fabricErr) && fabricErr.ErrorCode == "CapacityNotActive" {
		// Keep this an error for create, update, and import callers of get.
		// Only resource refresh may retain an existing state with a warning.
		return itemReadResult{
			Diagnostics:       utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil),
			CapacityNotActive: true,
		}
	}

	return itemReadResult{
		Diagnostics: utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound),
		CapacityNotActive: false,
	}
}

func preserveItemStateOnInactiveCapacity(
	req resource.ReadRequest,
	resp *resource.ReadResponse,
	resourceID string,
	result itemReadResult,
) bool {
	if !result.CapacityNotActive || len(result.Diagnostics) != 1 || !result.Diagnostics.HasError() {
		return false
	}

	resp.State = req.State
	resp.Diagnostics.AddWarning(
		common.WarningCapacityNotActiveHeader,
		fmt.Sprintf(common.WarningCapacityNotActiveDetails, resourceID)+"\n\n"+result.Diagnostics[0].Detail(),
	)

	return true
}
