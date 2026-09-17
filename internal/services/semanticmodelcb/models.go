// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package semanticmodelcb

import (
	"context"

	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabsemanticmodel "github.com/microsoft/fabric-sdk-go/fabric/semanticmodel"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type resourceSemanticModelConnectionBindingModel struct {
	WorkspaceID       customtypes.UUID                                             `tfsdk:"workspace_id"`
	SemanticModelID   customtypes.UUID                                             `tfsdk:"semantic_model_id"`
	ConnectivityType  types.String                                                 `tfsdk:"connectivity_type"`
	ConnectionID      customtypes.UUID                                             `tfsdk:"connection_id"`
	ConnectionDetails supertypes.SingleNestedObjectValueOf[connectionDetailsModel] `tfsdk:"connection_details"`
	Timeouts          timeoutsR.Value                                              `tfsdk:"timeouts"`
}

type connectionDetailsModel struct {
	Path types.String `tfsdk:"path"`
	Type types.String `tfsdk:"type"`
}

type requestBindSemanticModelConnection struct {
	fabsemanticmodel.BindSemanticModelConnectionRequest
}

func (to *requestBindSemanticModelConnection) set(ctx context.Context, from resourceSemanticModelConnectionBindingModel) diag.Diagnostics {
	details, diags := from.ConnectionDetails.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.ConnectionBinding = &fabsemanticmodel.ConnectionBinding{
		ConnectionDetails: &fabsemanticmodel.ListConnectionDetails{
			Path: details.Path.ValueStringPointer(),
			Type: details.Type.ValueStringPointer(),
		},
		ConnectivityType: (*fabsemanticmodel.ConnectivityType)(from.ConnectivityType.ValueStringPointer()),
		ID:               from.ConnectionID.ValueStringPointer(),
	}

	return nil
}

// setUnbind builds an unbind request reusing the data source reference from state.
func (to *requestBindSemanticModelConnection) setUnbind(ctx context.Context, from resourceSemanticModelConnectionBindingModel) diag.Diagnostics {
	details, diags := from.ConnectionDetails.Get(ctx)
	if diags.HasError() {
		return diags
	}

	none := string(fabsemanticmodel.ConnectivityTypeNone)
	to.ConnectionBinding = &fabsemanticmodel.ConnectionBinding{
		ConnectionDetails: &fabsemanticmodel.ListConnectionDetails{
			Path: details.Path.ValueStringPointer(),
			Type: details.Type.ValueStringPointer(),
		},
		ConnectivityType: (*fabsemanticmodel.ConnectivityType)(&none),
	}

	return nil
}
