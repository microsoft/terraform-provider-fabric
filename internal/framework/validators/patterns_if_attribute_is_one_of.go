// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package validators

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.String = PatternsIfAttributeIsOneOfValidator{}

type PatternsIfAttributeIsOneOfValidator struct {
	pathExpression path.Expression
	exceptedValues []attr.Value
	patterns       []string
	message        string
}

type PatternsIfAttributeIsOneOfRequest struct {
	Config         tfsdk.Config
	ConfigValue    types.String
	Path           path.Path
	PathExpression path.Expression
	ExceptedValues []attr.Value
}

type PatternsIfAttributeIsOneOfResponse struct {
	Diagnostics diag.Diagnostics
}

func PatternsIfAttributeIsOneOf(p path.Expression, exceptedValue []attr.Value, patterns []string, message string) PatternsIfAttributeIsOneOfValidator {
	return PatternsIfAttributeIsOneOfValidator{
		pathExpression: p,
		exceptedValues: exceptedValue,
		patterns:       patterns,
		message:        message,
	}
}

func (v PatternsIfAttributeIsOneOfValidator) Description(_ context.Context) string {
	if v.message != "" {
		return v.message
	}

	return fmt.Sprintf("value must match expression patterns '%s'", strings.Join(v.patterns, ", "))
}

func (v PatternsIfAttributeIsOneOfValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v PatternsIfAttributeIsOneOfValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateReq := PatternsIfAttributeIsOneOfRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &PatternsIfAttributeIsOneOfResponse{}

	v.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v PatternsIfAttributeIsOneOfValidator) Validate(ctx context.Context, req PatternsIfAttributeIsOneOfRequest, resp *PatternsIfAttributeIsOneOfResponse) {
	// If attribute configuration is unknown, there is nothing else to validate
	if req.ConfigValue.IsUnknown() {
		return
	}

	paths, diags := req.Config.PathMatches(ctx, req.PathExpression.Merge(v.pathExpression))
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if len(paths) == 0 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid configuration for attribute %s", req.Path),
			"Path must be set",
		)

		return
	}

	for _, p := range paths {
		var mpVal attr.Value

		diags = req.Config.GetAttribute(ctx, p, &mpVal)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Invalid configuration for attribute %s", req.Path),
				fmt.Sprintf("Unable to retrieve attribute path: %q", p),
			)

			return
		}

		// If the target attribute configuration is unknown or null, there is nothing else to validate
		if mpVal.IsNull() || mpVal.IsUnknown() {
			return
		}

		for _, expectedValue := range v.exceptedValues {
			// If the value of the target attribute is equal to one of the expected values, we need to validate the value of the current attribute
			if mpVal.Equal(expectedValue) {
				if req.ConfigValue.IsNull() {
					resp.Diagnostics.Append(
						validatordiag.InvalidAttributeValueDiagnostic(
							req.Path,
							"is empty, "+v.Description(ctx),
							req.ConfigValue.ValueString(),
						),
					)

					return
				}

				re, err := v.convertPatternsToRegexp(v.patterns)
				if err != nil {
					resp.Diagnostics.AddError(
						fmt.Sprintf("Invalid configuration for attribute %s", req.Path),
						fmt.Sprintf("Unable to compile regular expression: %q", err),
					)

					return
				}

				value := req.ConfigValue.ValueString()

				if !re.MatchString(value) {
					resp.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
						req.Path,
						v.Description(ctx),
						value,
					))
				}
			}
		}
	}
}

func (v PatternsIfAttributeIsOneOfValidator) convertPatternsToRegexp(patterns []string) (*regexp.Regexp, error) {
	p := make([]string, 0)

	p = append(p, "^(")

	for _, pattern := range patterns {
		p = append(p, regexp.QuoteMeta(pattern))
		if pattern != patterns[len(patterns)-1] {
			p = append(p, "|")
		}
	}

	p = append(p, ")$")

	out := strings.Join(p, "")
	out = strings.ReplaceAll(out, `\*`, ".+")

	re, err := regexp.Compile(out)
	if err != nil {
		return nil, err
	}

	return re, nil
}
