// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
)

var _ planmodifier.String = (*definitionContentSha256)(nil)

func DefinitionContentSha256(sourceFileAttr, tokensAttr path.Expression) planmodifier.String {
	return &definitionContentSha256{
		sourceFileAttr: sourceFileAttr,
		tokensAttr:     tokensAttr,
	}
}

type definitionContentSha256 struct {
	sourceFileAttr path.Expression
	tokensAttr     path.Expression
}

func (pm *definitionContentSha256) Description(_ context.Context) string {
	return "Generate SHA256 hash of the JSON normalized content of the file."
}

func (pm *definitionContentSha256) MarkdownDescription(ctx context.Context) string {
	return pm.Description(ctx)
}

func (pm *definitionContentSha256) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	sourceFilePlanPaths, diags := req.Plan.PathMatches(ctx, req.PathExpression.Merge(pm.sourceFileAttr))
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	tokensPlanPaths, diags := req.Plan.PathMatches(ctx, req.PathExpression.Merge(pm.tokensAttr))
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	var sourceFile types.String
	var defPartTokens customtypes.MapOfString

	if resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, sourceFilePlanPaths[0], &sourceFile)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, tokensPlanPaths[0], &defPartTokens)...); resp.Diagnostics.HasError() {
		return
	}

	if sourceFile.IsNull() || sourceFile.IsUnknown() {
		resp.PlanValue = types.StringUnknown()

		return
	}

	tokensValue := make(map[string]string)

	if defPartTokens.IsKnown() {
		for _, v := range defPartTokens.Elements() {
			if v.IsNull() || v.IsUnknown() {
				resp.PlanValue = types.StringUnknown()

				return
			}
		}

		tokens, diags := defPartTokens.Get(ctx)
		if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
			return
		}

		for k, v := range tokens {
			if !v.IsNull() && !v.IsUnknown() {
				tokensValue[k] = v.ValueString()
			}
		}
	}

	_, sha256Value, diags := transforms.SourceFileToPayload(sourceFile.ValueString(), tokensValue)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.PlanValue = types.StringValue(sha256Value)
}
