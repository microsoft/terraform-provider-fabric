// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package digitaltwinbuilderflow

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabdigitaltwinbuilderflow "github.com/microsoft/fabric-sdk-go/fabric/digitaltwinbuilderflow"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type digitalTwinBuilderItemReferenceModel struct {
	ItemID        customtypes.UUID `tfsdk:"item_id"`
	ReferenceType types.String     `tfsdk:"reference_type"`
	WorkspaceID   customtypes.UUID `tfsdk:"workspace_id"`
}

type digitalTwinBuilderFlowConfigPropertiesModel struct {
	DigitalTwinBuilderItemReference supertypes.SingleNestedObjectValueOf[digitalTwinBuilderItemReferenceModel] `tfsdk:"digital_twin_builder_item_reference"`
}

func (to *digitalTwinBuilderFlowConfigPropertiesModel) set(ctx context.Context, from fabdigitaltwinbuilderflow.Properties) diag.Diagnostics {
	reference := from.DigitalTwinBuilderItemReference.GetItemReference()

	switch *reference.ReferenceType {
	case fabdigitaltwinbuilderflow.ItemReferenceTypeByID:
		digitalTwinBuilderItemReference := supertypes.NewSingleNestedObjectValueOfNull[digitalTwinBuilderItemReferenceModel](ctx)

		refByID, ok := from.DigitalTwinBuilderItemReference.(*fabdigitaltwinbuilderflow.ItemReferenceByID)
		if !ok {
			var diags diag.Diagnostics
			diags.AddError(
				"Type Assertion Failed",
				"Failed to convert DigitalTwinBuilderItemReference to ItemReferenceByID",
			)

			return diags
		}

		refType := string(*refByID.ReferenceType)

		digitalTwinBuilderItemReferenceModel := &digitalTwinBuilderItemReferenceModel{
			ItemID:        customtypes.NewUUIDPointerValue(refByID.ItemID),
			ReferenceType: types.StringPointerValue(&refType),
			WorkspaceID:   customtypes.NewUUIDPointerValue(refByID.WorkspaceID),
		}

		if diags := digitalTwinBuilderItemReference.Set(ctx, digitalTwinBuilderItemReferenceModel); diags.HasError() {
			return diags
		}

		to.DigitalTwinBuilderItemReference = digitalTwinBuilderItemReference
	default:
		var diags diag.Diagnostics
		diags.AddError(
			"Unsupported Item Reference Type",
			fmt.Sprintf("Item reference type '%s' is not supported", *reference.ReferenceType),
		)

		return diags
	}

	return nil
}
