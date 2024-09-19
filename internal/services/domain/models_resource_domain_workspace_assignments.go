// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type resourceDomainWorkspaceAssignmentsModel struct {
	ID           customtypes.UUID `tfsdk:"id"`
	DomainID     customtypes.UUID `tfsdk:"domain_id"`
	WorkspaceIDs types.Set        `tfsdk:"workspace_ids"`
	Timeouts     timeouts.Value   `tfsdk:"timeouts"`
}

func (to *resourceDomainWorkspaceAssignmentsModel) setWorkspaces(ctx context.Context, from []fabadmin.DomainWorkspace) diag.Diagnostics {
	elements := []customtypes.UUID{}

	for _, entity := range from {
		elements = append(elements, customtypes.NewUUIDPointerValue(entity.ID))
	}

	values, diags := types.SetValueFrom(ctx, customtypes.UUIDType{}, elements)
	if diags.HasError() {
		return diags
	}

	to.WorkspaceIDs = values

	return nil
}

type requestCreateDomainWorkspaceAssignments struct {
	fabadmin.AssignDomainWorkspacesByIDsRequest
}

func (to *requestCreateDomainWorkspaceAssignments) set(ctx context.Context, from resourceDomainWorkspaceAssignmentsModel) diag.Diagnostics {
	workspaceIDs, diags := getWorkspaceIDs(ctx, from)
	if diags.HasError() {
		return diags
	}

	to.WorkspacesIDs = workspaceIDs

	return nil
}

type requestDeleteDomainWorkspaceAssignments struct {
	fabadmin.UnassignDomainWorkspacesByIDsRequest
}

func (to *requestDeleteDomainWorkspaceAssignments) set(ctx context.Context, from resourceDomainWorkspaceAssignmentsModel) diag.Diagnostics {
	workspaceIDs, diags := getWorkspaceIDs(ctx, from)
	if diags.HasError() {
		return diags
	}

	to.WorkspacesIDs = workspaceIDs

	return nil
}

func getWorkspaceIDs(ctx context.Context, from resourceDomainWorkspaceAssignmentsModel) ([]string, diag.Diagnostics) {
	values := []string{}

	elements := make([]customtypes.UUID, 0, len(from.WorkspaceIDs.Elements()))

	diags := from.WorkspaceIDs.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		return nil, diags
	}

	for _, element := range elements {
		values = append(values, element.ValueString())
	}

	return values, nil
}
