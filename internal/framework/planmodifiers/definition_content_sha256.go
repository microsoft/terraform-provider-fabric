// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/params"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
)

var _ planmodifier.String = (*definitionContentSha256)(nil)

func DefinitionContentSha256(sourceAttr, processingModeAttr, tokensAttr, parametersAttr, tokensDelimiterAttr path.Expression) planmodifier.String {
	return &definitionContentSha256{
		source:          sourceAttr,
		processingMode:  processingModeAttr,
		tokens:          tokensAttr,
		parameters:      parametersAttr,
		tokensDelimiter: tokensDelimiterAttr,
	}
}

type definitionContentSha256 struct {
	source          path.Expression
	processingMode  path.Expression
	tokens          path.Expression
	parameters      path.Expression
	tokensDelimiter path.Expression
}

func (pm *definitionContentSha256) Description(_ context.Context) string {
	return "Generate SHA256 hash of the JSON normalized content of the file."
}

func (pm *definitionContentSha256) MarkdownDescription(ctx context.Context) string {
	return pm.Description(ctx)
}

func (pm *definitionContentSha256) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) { //nolint:gocognit,gocyclo
	sourcePlanPaths, diags := req.Plan.PathMatches(ctx, req.PathExpression.Merge(pm.source))
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	processingModePlanPaths, diags := req.Plan.PathMatches(ctx, req.PathExpression.Merge(pm.processingMode))
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	tokensPlanPaths, diags := req.Plan.PathMatches(ctx, req.PathExpression.Merge(pm.tokens))
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	tokensDelimiterPlanPaths, diags := req.Plan.PathMatches(ctx, req.PathExpression.Merge(pm.tokensDelimiter))
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	parametersPlanPaths, diags := req.Plan.PathMatches(ctx, req.PathExpression.Merge(pm.parameters))
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	var source types.String
	var tokens supertypes.MapValueOf[types.String]
	var tokensDelimiter types.String
	var parameters supertypes.SetNestedObjectValueOf[params.ParametersModel]
	var processingMode types.String

	if resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, sourcePlanPaths[0], &source)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, tokensPlanPaths[0], &tokens)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, tokensDelimiterPlanPaths[0], &tokensDelimiter)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, parametersPlanPaths[0], &parameters)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, processingModePlanPaths[0], &processingMode)...); resp.Diagnostics.HasError() {
		return
	}

	if source.IsNull() || source.IsUnknown() {
		resp.PlanValue = types.StringUnknown()

		return
	}

	tokensValue := make(map[string]string)

	if tokens.IsKnown() {
		for _, v := range tokens.Elements() {
			if v.IsNull() || v.IsUnknown() {
				resp.PlanValue = types.StringUnknown()

				return
			}
		}

		tokensMap, diags := tokens.Get(ctx)
		if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
			return
		}

		for k, v := range tokensMap {
			if !v.IsNull() && !v.IsUnknown() {
				tokensValue[k] = v.ValueString()
			}
		}
	}

	var parametersSlice []*params.ParametersModel

	if parameters.IsKnown() {
		parametersSlice, diags = parameters.Get(ctx)

		for _, param := range parametersSlice {
			if param.Value.IsNull() || param.Value.IsUnknown() ||
				param.Type.IsNull() || param.Type.IsUnknown() ||
				param.Find.IsNull() || param.Find.IsUnknown() {
				resp.PlanValue = types.StringUnknown()

				return
			}
		}

		if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
			return
		}
	}

	_, sha256Value, diags := transforms.SourceFileToPayload(source.ValueString(), processingMode.ValueString(), tokensValue, parametersSlice, tokensDelimiter.ValueString())
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.PlanValue = types.StringValue(sha256Value)
}
