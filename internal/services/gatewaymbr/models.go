// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gatewaymbr

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseGatewayMemberModel struct {
	ID          customtypes.UUID                                     `tfsdk:"id"`
	GatewayID   customtypes.UUID                                     `tfsdk:"gateway_id"`
	DisplayName types.String                                         `tfsdk:"display_name"`
	Version     types.String                                         `tfsdk:"version"`
	Enabled     types.Bool                                           `tfsdk:"enabled"`
	PublicKey   supertypes.SingleNestedObjectValueOf[publicKeyModel] `tfsdk:"public_key"`
}

func (to *baseGatewayMemberModel) set(ctx context.Context, gatewayID string, from fabcore.OnPremisesGatewayMember) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.GatewayID = customtypes.NewUUIDValue(gatewayID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Version = types.StringPointerValue(from.Version)
	to.Enabled = types.BoolPointerValue(from.Enabled)

	publicKey := supertypes.NewSingleNestedObjectValueOfNull[publicKeyModel](ctx)

	if from.PublicKey != nil {
		publicKeyModel := &publicKeyModel{}
		publicKeyModel.set(*from.PublicKey)

		if diags := publicKey.Set(ctx, publicKeyModel); diags.HasError() {
			return diags
		}
	}

	to.PublicKey = publicKey

	return nil
}

/*
DATA-SOURCE
*/

type dataSourceGatewayMemberModel struct {
	baseGatewayMemberModel
	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceGatewayMembersModel struct {
	GatewayID customtypes.UUID                                          `tfsdk:"gateway_id"`
	Values    supertypes.SetNestedObjectValueOf[baseGatewayMemberModel] `tfsdk:"values"`
	Timeouts  timeoutsD.Value                                           `tfsdk:"timeouts"`
}

func (to *dataSourceGatewayMembersModel) setValues(ctx context.Context, gatewayID string, from []fabcore.OnPremisesGatewayMember) diag.Diagnostics {
	slice := make([]*baseGatewayMemberModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseGatewayMemberModel

		if diags := entityModel.set(ctx, gatewayID, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

/*
RESOURCE
*/

type resourceGatewayMemberModel struct {
	baseGatewayMemberModel
	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateGatewayMember struct {
	fabcore.UpdateGatewayMemberRequest
}

func (to *requestCreateGatewayMember) set(ctx context.Context, from resourceGatewayMemberModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Enabled = from.Enabled.ValueBoolPointer()
}

type requestUpdateMemberGateway struct {
	fabcore.UpdateGatewayMemberRequest
}

func (to *requestUpdateMemberGateway) set(from resourceGatewayMemberModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Enabled = from.Enabled.ValueBoolPointer()
}

/*
HELPER MODELS
*/

type publicKeyModel struct {
	Exponent types.String `tfsdk:"exponent"`
	Modulus  types.String `tfsdk:"modulus"`
}

func (to *publicKeyModel) set(from fabcore.PublicKey) {
	to.Exponent = types.StringPointerValue(from.Exponent)
	to.Modulus = types.StringPointerValue(from.Modulus)
}
