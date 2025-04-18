// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domainwa

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseDomainWorkspaceAssignmentsModel struct {
	// ID           customtypes.UUID                         `tfsdk:"id"`
	DomainID     customtypes.UUID                        `tfsdk:"domain_id"`
	WorkspaceIDs supertypes.SetValueOf[customtypes.UUID] `tfsdk:"workspace_ids"`
}

func (to *baseDomainWorkspaceAssignmentsModel) setWorkspaceIDs(ctx context.Context, from []fabadmin.DomainWorkspace) diag.Diagnostics {
	v := supertypes.NewSetValueOfNull[customtypes.UUID](ctx)

	elements := make([]customtypes.UUID, 0, len(from))
	for _, element := range from {
		elements = append(elements, customtypes.NewUUIDPointerValue(element.ID))
	}

	if diags := v.Set(ctx, elements); diags.HasError() {
		return diags
	}

	to.WorkspaceIDs = v

	return nil
}

/*
DATA-SOURCE
*/

type dataSourceDomainWorkspaceAssignmentsModel struct {
	baseDomainWorkspaceAssignmentsModel
	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
RESOURCE
*/

type resourceDomainWorkspaceAssignmentsModel struct {
	baseDomainWorkspaceAssignmentsModel
	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
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
	elements, diags := from.WorkspaceIDs.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	values := make([]string, 0, len(elements))

	for _, element := range elements {
		values = append(values, element.ValueString())
	}

	return values, nil
}
