// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package customtypes

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const (
	UUIDTypeErrorInvalidStringHeader  = "Invalid UUID String Value"
	UUIDTypeErrorInvalidStringDetails = "A string value was provided that is not valid UUID string format.\n\nGiven Value: %s\nError: %s"
)

var (
	_ basetypes.StringValuable                   = (*UUIDValue)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*UUIDValue)(nil)
	_ xattr.ValidateableAttribute                = (*UUIDValue)(nil)
	_ function.ValidateableParameter             = (*UUIDValue)(nil)
)

type UUID = UUIDValue

type UUIDValue struct {
	basetypes.StringValue
}

func (v UUIDValue) Type(_ context.Context) attr.Type {
	return UUIDType{}
}

func (v UUIDValue) Equal(o attr.Value) bool {
	other, ok := o.(UUIDValue)
	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v UUIDValue) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(UUIDValue)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	oldUUID, err := uuid.ParseUUID(v.ValueString())
	if err != nil {
		diags.AddError("expected old value to be a valid UUID", err.Error())
	}

	newUUID, err := uuid.ParseUUID(newValue.ValueString())
	if err != nil {
		diags.AddError("expected new value to be a valid UUID", err.Error())
	}

	if diags.HasError() {
		return false, diags
	}

	return reflect.DeepEqual(oldUUID, newUUID), diags
}

func (v UUIDValue) ValidateAttribute(_ context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	if _, err := uuid.ParseUUID(v.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			UUIDTypeErrorInvalidStringHeader,
			fmt.Sprintf(UUIDTypeErrorInvalidStringDetails, v.ValueString(), err.Error()),
		)

		return
	}
}

func (v UUIDValue) ValidateParameter(_ context.Context, req function.ValidateParameterRequest, resp *function.ValidateParameterResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	if _, err := uuid.ParseUUID(v.ValueString()); err != nil {
		resp.Error = function.NewArgumentFuncError(
			req.Position,
			UUIDTypeErrorInvalidStringHeader+": "+fmt.Sprintf(UUIDTypeErrorInvalidStringDetails, v.ValueString(), err.Error()),
		)

		return
	}
}

func (v UUIDValue) ValueUUID() (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v.IsNull() {
		diags.Append(diag.NewErrorDiagnostic(UUIDTypeErrorInvalidStringHeader, "UUID string value is null"))

		return "", diags
	}

	if v.IsUnknown() {
		diags.Append(diag.NewErrorDiagnostic(UUIDTypeErrorInvalidStringHeader, "UUID string value is unknown"))

		return "", diags
	}

	value, err := uuid.ParseUUID(v.ValueString())
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic(
			UUIDTypeErrorInvalidStringHeader,
			fmt.Sprintf(UUIDTypeErrorInvalidStringDetails, v.ValueString(), err.Error()),
		))

		return "", diags
	}

	return string(value), nil
}
