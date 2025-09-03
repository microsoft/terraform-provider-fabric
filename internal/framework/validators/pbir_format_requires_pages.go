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

var _ validator.Map = PBIRFormatRequiresPagesValidator{}

type PBIRFormatRequiresPagesValidator struct {
	formatPath    path.Expression
	pagesPattern  string
	errorMessage  string
}

type PBIRFormatRequiresPagesRequest struct {
	Config         tfsdk.Config
	ConfigValue    types.Map
	Path           path.Path
	PathExpression path.Expression
}

type PBIRFormatRequiresPagesResponse struct {
	Diagnostics diag.Diagnostics
}

func PBIRFormatRequiresPages(formatPath path.Expression, pagesPattern string, errorMessage string) PBIRFormatRequiresPagesValidator {
	return PBIRFormatRequiresPagesValidator{
		formatPath:   formatPath,
		pagesPattern: pagesPattern,
		errorMessage: errorMessage,
	}
}

func (v PBIRFormatRequiresPagesValidator) Description(_ context.Context) string {
	if v.errorMessage != "" {
		return v.errorMessage
	}
	return fmt.Sprintf("PBIR format requires at least one page file matching '%s' pattern", v.pagesPattern)
}

func (v PBIRFormatRequiresPagesValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v PBIRFormatRequiresPagesValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	validateReq := PBIRFormatRequiresPagesRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &PBIRFormatRequiresPagesResponse{}

	v.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (v PBIRFormatRequiresPagesValidator) Validate(ctx context.Context, req PBIRFormatRequiresPagesRequest, resp *PBIRFormatRequiresPagesResponse) {
	// If map configuration is unknown, there is nothing else to validate
	if req.ConfigValue.IsUnknown() {
		return
	}

	// Get the format value
	paths, diags := req.Config.PathMatches(ctx, req.PathExpression.Merge(v.formatPath))
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if len(paths) == 0 {
		// Format not set, nothing to validate
		return
	}

	for _, p := range paths {
		var formatVal attr.Value

		diags = req.Config.GetAttribute(ctx, p, &formatVal)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Invalid configuration for attribute %s", req.Path),
				fmt.Sprintf("Unable to retrieve format attribute path: %q", p),
			)
			return
		}

		// If the format attribute configuration is unknown or null, there is nothing else to validate
		if formatVal.IsNull() || formatVal.IsUnknown() {
			return
		}

		// Check if format is PBIR
		formatString, ok := formatVal.(types.String)
		if !ok {
			return
		}

		if formatString.ValueString() == "PBIR" {
			// Now check if map contains at least one page file
			if req.ConfigValue.IsNull() {
				resp.Diagnostics.Append(
					validatordiag.InvalidAttributeValueDiagnostic(
						req.Path,
						v.Description(ctx),
						"<null>",
					),
				)
				return
			}

			// Check if any key matches the pages pattern
			hasPages := false
			elements := req.ConfigValue.Elements()
			
			// Convert pattern to regex
			pagesRegex, err := v.convertPatternToRegexp(v.pagesPattern)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Invalid configuration for attribute %s", req.Path),
					fmt.Sprintf("Unable to compile pages pattern regex: %q", err),
				)
				return
			}

			for key := range elements {
				if pagesRegex.MatchString(key) {
					hasPages = true
					break
				}
			}

			if !hasPages {
				resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
					req.Path,
					v.Description(ctx),
					"no pages found",
				))
			}
		}
	}
}

func (v PBIRFormatRequiresPagesValidator) convertPatternToRegexp(pattern string) (*regexp.Regexp, error) {
	// Escape the pattern and convert wildcards to regex
	escapedPattern := regexp.QuoteMeta(pattern)
	regexPattern := strings.ReplaceAll(escapedPattern, `\*`, ".+")
	regexPattern = "^" + regexPattern + "$"

	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, err
	}

	return re, nil
}