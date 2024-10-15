// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fablakehouse "github.com/microsoft/fabric-sdk-go/fabric/lakehouse"
)

type resourceLakehouseModel struct {
	baseLakehouseModel
	Configuration supertypes.SingleNestedObjectValueOf[lakehouseConfigurationModel] `tfsdk:"configuration"`
	Timeouts      timeouts.Value                                                    `tfsdk:"timeouts"`
}
type lakehouseConfigurationModel struct {
	EnableSchemas types.Bool `tfsdk:"enable_schemas"`
}
type requestCreateLakehouse struct {
	fablakehouse.CreateLakehouseRequest
}

func (to *requestCreateLakehouse) set(ctx context.Context, from resourceLakehouseModel) diag.Diagnostics {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()

	if !from.Configuration.IsNull() && !from.Configuration.IsUnknown() {
		configuration, diags := from.Configuration.Get(ctx)
		if diags.HasError() {
			return diags
		}

		if configuration.EnableSchemas.ValueBool() {
			to.CreationPayload = &fablakehouse.CreationPayload{
				EnableSchemas: configuration.EnableSchemas.ValueBoolPointer(),
			}
		}
	}

	return nil
}

type requestUpdateLakehouse struct {
	fablakehouse.UpdateLakehouseRequest
}

func (to *requestUpdateLakehouse) set(from resourceLakehouseModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
}
