// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type principalModel struct {
	ID   customtypes.UUID `tfsdk:"id"`
	Type types.String     `tfsdk:"type"`
}

type resourceDomainRoleAssignmentsModel struct {
	ID         customtypes.UUID                                  `tfsdk:"id"`
	DomainID   customtypes.UUID                                  `tfsdk:"domain_id"`
	Role       types.String                                      `tfsdk:"role"`
	Principals supertypes.SetNestedObjectValueOf[principalModel] `tfsdk:"principals"`
	Timeouts   timeouts.Value                                    `tfsdk:"timeouts"`
}

type requestCreateDomainRoleAssignments struct {
	fabadmin.DomainRoleAssignmentRequest
}

func (to *requestCreateDomainRoleAssignments) set(ctx context.Context, from resourceDomainRoleAssignmentsModel) diag.Diagnostics {
	to.Type = (*fabadmin.DomainRole)(from.Role.ValueStringPointer())
	to.Principals = []fabadmin.Principal{}

	principals, diags := from.Principals.Get(ctx)
	if diags.HasError() {
		return diags
	}

	for _, principal := range principals {
		to.Principals = append(to.Principals, fabadmin.Principal{
			ID:   principal.ID.ValueStringPointer(),
			Type: (*fabadmin.PrincipalType)(principal.Type.ValueStringPointer()),
		})
	}

	return nil
}

type requestDeleteDomainRoleAssignments struct {
	fabadmin.DomainRoleUnassignmentRequest
}

func (to *requestDeleteDomainRoleAssignments) set(ctx context.Context, from resourceDomainRoleAssignmentsModel) diag.Diagnostics {
	to.Type = (*fabadmin.DomainRole)(from.Role.ValueStringPointer())
	to.Principals = []fabadmin.Principal{}

	principals, diags := from.Principals.Get(ctx)
	if diags.HasError() {
		return diags
	}

	for _, principal := range principals {
		to.Principals = append(to.Principals, fabadmin.Principal{
			ID:   principal.ID.ValueStringPointer(),
			Type: (*fabadmin.PrincipalType)(principal.Type.ValueStringPointer()),
		})
	}

	return nil
}
