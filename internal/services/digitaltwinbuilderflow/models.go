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
)

type digitalTwinBuilderItemReferenceModel struct {
	ItemID        types.String `tfsdk:"item_id"`
	ReferenceType types.String `tfsdk:"reference_type"`
	WorkspaceID   types.String `tfsdk:"workspace_id"`
}

type digitalTwinBuilderFlowConfigPropertiesModel struct {
	digitalTwinBuilderItemReference supertypes.SingleNestedObjectValueOf[digitalTwinBuilderItemReferenceModel] `tfsdk:"digital_twin_builder_item_reference"`
}

func (to *digitalTwinBuilderFlowConfigPropertiesModel) set(ctx context.Context, from fabdigitaltwinbuilderflow.Properties) diag.Diagnostics {
	reference := from.DigitalTwinBuilderItemReference.GetItemReference()

	switch *reference.ReferenceType {
	case fabdigitaltwinbuilderflow.ItemReferenceTypeByID:
		digitalTwinBuilderItemReference := supertypes.NewSingleNestedObjectValueOfNull[digitalTwinBuilderItemReferenceModel](ctx)
		refByID := from.DigitalTwinBuilderItemReference.(*fabdigitaltwinbuilderflow.ItemReferenceByID)
		refType := string(*refByID.ReferenceType)

		digitalTwinBuilderItemReferenceModel := &digitalTwinBuilderItemReferenceModel{
			ItemID:        types.StringPointerValue(refByID.ItemID),
			ReferenceType: types.StringPointerValue(&refType),
			WorkspaceID:   types.StringPointerValue(refByID.WorkspaceID),
		}

		if diags := digitalTwinBuilderItemReference.Set(ctx, digitalTwinBuilderItemReferenceModel); diags.HasError() {
			return diags
		}
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
