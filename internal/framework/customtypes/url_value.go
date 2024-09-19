// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package customtypes

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const (
	URLTypeErrorInvalidStringHeader  = "Invalid absolute URL String Value"
	URLTypeErrorInvalidStringDetails = "A string value was provided that is not valid absolute URL string format.\n\nGiven Value: %s\nError: %s"
)

var (
	_ basetypes.StringValuable                   = (*URLValue)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*URLValue)(nil)
	_ xattr.ValidateableAttribute                = (*URLValue)(nil)
	_ function.ValidateableParameter             = (*URLValue)(nil)
)

type URL = URLValue

type URLValue struct {
	basetypes.StringValue
}

func (v URLValue) Type(_ context.Context) attr.Type {
	return URLType{}
}

func (v URLValue) Equal(o attr.Value) bool {
	other, ok := o.(URLValue)
	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v URLValue) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(URLValue)
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

	oldURL, err := v.check(v.ValueString())
	if err != nil {
		diags.AddError("expected old value to be a valid absolute URL", err.Error())
	}

	newURL, err := v.check(newValue.ValueString())
	if err != nil {
		diags.AddError("expected new value to be a valid absolute URL", err.Error())
	}

	if diags.HasError() {
		return false, diags
	}

	return oldURL == newURL, diags
}

func (v URLValue) ValidateAttribute(_ context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	_, err := v.check(v.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			URLTypeErrorInvalidStringHeader,
			fmt.Sprintf(URLTypeErrorInvalidStringDetails, v.ValueString(), err.Error()),
		)

		return
	}
}

func (v URLValue) ValidateParameter(_ context.Context, req function.ValidateParameterRequest, resp *function.ValidateParameterResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	_, err := v.check(v.ValueString())
	if err != nil {
		resp.Error = function.NewArgumentFuncError(
			req.Position,
			URLTypeErrorInvalidStringHeader+": "+fmt.Sprintf(URLTypeErrorInvalidStringDetails, v.ValueString(), err.Error()),
		)

		return
	}
}

func (v URLValue) ValueURL() (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v.IsNull() {
		diags.Append(diag.NewErrorDiagnostic(URLTypeErrorInvalidStringHeader, "URL string value is null"))

		return "", diags
	}

	if v.IsUnknown() {
		diags.Append(diag.NewErrorDiagnostic(URLTypeErrorInvalidStringHeader, "URL string value is unknown"))

		return "", diags
	}

	value, err := v.check(v.ValueString())
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic(
			URLTypeErrorInvalidStringHeader,
			fmt.Sprintf(URLTypeErrorInvalidStringDetails, v.ValueString(), err.Error()),
		))

		return "", diags
	}

	return value, nil
}

func (v URLValue) check(input string) (string, error) {
	value, err := url.ParseRequestURI(input)
	if err != nil || value.Host == "" {
		urlErr := &url.Error{}
		if errors.As(err, &urlErr) {
			err = urlErr.Unwrap()
		} else {
			err = errors.New("invalid URI for request")
		}

		return "", err
	}

	return value.String(), nil
}
