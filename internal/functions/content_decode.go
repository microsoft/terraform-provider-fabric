// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package functions

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ohler55/ojg/jp"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/typeutils"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
)

var _ function.Function = (*functionContentDecode)(nil)

func NewFunctionContentDecode() function.Function {
	return &functionContentDecode{}
}

type functionContentDecode struct{}

func (f *functionContentDecode) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "content_decode"
}

func (f *functionContentDecode) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Decode a Definition JSON object.",
		MarkdownDescription: "Given a Base64 Gzip encoded content, will decode and return a Definition JSON object representation of that resource.",
		Parameters: []function.Parameter{
			function.StringParameter{
				MarkdownDescription: "The Base64 Gzip content.",
				Name:                "content",
			},
		},
		VariadicParameter: function.StringParameter{
			MarkdownDescription: "Filter JSON output using a [JSONPath](https://datatracker.ietf.org/doc/html/rfc9535) expression.",
			Name:                "expression",
			AllowNullValue:      true,
		},
		Return: function.DynamicReturn{},
	}
}

func (f *functionContentDecode) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	tflog.Debug(ctx, "CONTENT DECODE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "CONTENT DECODE", map[string]any{
		"arguments": req.Arguments,
	})

	var (
		inputContent    string
		inputExpression []string
		dResult         types.Dynamic
	)

	if resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &inputContent, &inputExpression)); resp.Error != nil {
		return
	}

	if inputContent == "" {
		resp.Error = function.NewFuncError("Parameter 'content' is required")

		return
	}

	contentDecoded, err := transforms.Base64GzipDecode(inputContent)
	if err != nil {
		resp.Error = function.NewFuncError("Failed to decode content: " + err.Error())

		return
	}

	if transforms.IsJSON(contentDecoded) { //nolint:nestif
		var contentJSON any

		if err := json.Unmarshal([]byte(contentDecoded), &contentJSON); err != nil {
			resp.Error = function.NewFuncError("Failed to unmarshal JSON: " + err.Error())

			return
		}

		var results []string

		if len(inputExpression) > 0 && inputExpression[0] != "" {
			jpExpression, err := jp.ParseString(inputExpression[0])
			if err != nil {
				resp.Error = function.NewFuncError("Failed to parse JSONPath expression: " + err.Error())

				return
			}

			jpIter := jpExpression.Get(contentJSON)

			// Add error check for empty results
			if len(jpIter) == 0 {
				resp.Error = function.NewFuncError("JSONPath expression did not match any elements")

				return
			}

			for _, v := range jpIter {
				jsonPretty, err := json.MarshalIndent(v, "", "  ")
				if err != nil {
					resp.Error = function.NewFuncError("Failed to marshal JSON: " + err.Error())

					return
				}

				results = append(results, string(jsonPretty))
			}
		} else {
			jsonPretty, err := json.MarshalIndent(contentJSON, "", "  ")
			if err != nil {
				resp.Error = function.NewFuncError("Failed to marshal JSON: " + err.Error())

				return
			}

			results = append(results, string(jsonPretty))
		}

		var err error

		dResult, err = typeutils.JSONToDynamicImplied([]byte(strings.Join(results, "\n")))
		if err != nil {
			resp.Error = function.NewFuncError("Failed to parse JSON: " + err.Error())

			return
		}
	} else {
		dResult = types.DynamicValue(types.StringValue(contentDecoded))
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, dResult))

	tflog.Debug(ctx, "RUN", map[string]any{
		"action": "end",
	})
}
