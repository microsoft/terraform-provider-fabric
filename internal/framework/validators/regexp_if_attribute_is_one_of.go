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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = RegexpIfAttributeIsOneOfValidator{}

type RegexpIfAttributeIsOneOfValidator struct {
	pathExpression path.Expression
	exceptedValues []attr.Value
	patterns       []string
	message        string
}

func RegexpIfAttributeIsOneOf(p path.Expression, exceptedValue []attr.Value, patterns []string, message string) RegexpIfAttributeIsOneOfValidator {
	return RegexpIfAttributeIsOneOfValidator{
		pathExpression: p,
		exceptedValues: exceptedValue,
		patterns:       patterns,
		message:        message,
	}
}

func (v RegexpIfAttributeIsOneOfValidator) Description(_ context.Context) string {
	if v.message != "" {
		return v.message
	}

	return fmt.Sprintf("value must match pattern expression '%s'", strings.Join(v.patterns, ", "))
}

func (v RegexpIfAttributeIsOneOfValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v RegexpIfAttributeIsOneOfValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	paths, diags := req.Config.PathMatches(ctx, req.PathExpression.Merge(v.pathExpression))
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)

		return
	}

	if len(paths) == 0 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid configuration for attribute %s", req.Path),
			"Path must be set",
		)

		return
	}

	p := paths[0]

	// mpVal is the value of the attribute in the path
	var mpVal attr.Value
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, p, &mpVal)...)

	if resp.Diagnostics.HasError() {
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
		if mpVal.Equal(expectedValue) || mpVal.String() == expectedValue.String() {
			if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					p,
					fmt.Sprintf("Invalid configuration for attribute %s", req.Path),
					"Value is empty. "+v.Description(ctx),
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

func (v RegexpIfAttributeIsOneOfValidator) convertPatternsToRegexp(patterns []string) (*regexp.Regexp, error) {
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
