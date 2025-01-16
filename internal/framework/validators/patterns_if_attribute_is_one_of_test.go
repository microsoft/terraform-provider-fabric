// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package validators_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/validators"
)

func TestUnit_PatternsIfAttributeIsOneOfValidator(t *testing.T) { //nolint:maintidx
	t.Parallel()

	type testCase struct {
		req             validators.PatternsIfAttributeIsOneOfRequest
		in              path.Expression
		inPath          path.Path
		exceptedValues  []attr.Value
		patterns        []string
		message         string
		expError        bool
		expErrorMessage string
	}

	testCases := map[string]testCase{
		"multi-not-match": {
			req: validators.PatternsIfAttributeIsOneOfRequest{
				ConfigValue:    types.StringValue("foo value"),
				Path:           path.Root("foo"),
				PathExpression: path.MatchRoot("foo"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.StringAttribute{},
							"bar": schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.String,
							"bar": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.String, "foo value"),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
					}),
				},
			},
			in:     path.MatchRoot("bar"),
			inPath: path.Root("foo"),
			exceptedValues: []attr.Value{
				types.StringValue("bar value"),
			},
			patterns:        []string{"foo", "bar", "baz", "test/*.json"},
			message:         "",
			expError:        true,
			expErrorMessage: `value must match expression patterns 'foo, bar, baz, test/*.json'`,
		},
		"multi-match": {
			req: validators.PatternsIfAttributeIsOneOfRequest{
				ConfigValue:    types.StringValue("test/foo.json"),
				Path:           path.Root("foo"),
				PathExpression: path.MatchRoot("foo"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.StringAttribute{},
							"bar": schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.String,
							"bar": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.String, "foo value"),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
					}),
				},
			},
			in:     path.MatchRoot("bar"),
			inPath: path.Root("foo"),
			exceptedValues: []attr.Value{
				types.StringValue("bar value"),
			},
			patterns: []string{"foo", "bar", "baz", "test/*.json"},
			message:  "",
			expError: false,
		},
		"one-not-match": {
			req: validators.PatternsIfAttributeIsOneOfRequest{
				ConfigValue:    types.StringValue("foo value"),
				Path:           path.Root("foo"),
				PathExpression: path.MatchRoot("foo"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.StringAttribute{},
							"bar": schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.String,
							"bar": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.String, "foo value"),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
					}),
				},
			},
			in:     path.MatchRoot("bar"),
			inPath: path.Root("foo"),
			exceptedValues: []attr.Value{
				types.StringValue("bar value"),
			},
			patterns:        []string{"baz"},
			message:         "",
			expError:        true,
			expErrorMessage: `value must match expression patterns 'baz'`,
		},
		"one-match": {
			req: validators.PatternsIfAttributeIsOneOfRequest{
				ConfigValue:    types.StringValue("foo value"),
				Path:           path.Root("foo"),
				PathExpression: path.MatchRoot("foo"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.StringAttribute{},
							"bar": schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.String,
							"bar": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.String, "foo value"),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
					}),
				},
			},
			in:     path.MatchRoot("bar"),
			inPath: path.Root("foo"),
			exceptedValues: []attr.Value{
				types.StringValue("bar value"),
			},
			patterns: []string{"foo value"},
			message:  "",
			expError: false,
		},
		"custom-msg-err": {
			req: validators.PatternsIfAttributeIsOneOfRequest{
				ConfigValue:    types.StringValue("foo value"),
				Path:           path.Root("foo"),
				PathExpression: path.MatchRoot("foo"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.StringAttribute{},
							"bar": schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.String,
							"bar": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.String, "foo value"),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
					}),
				},
			},
			in:     path.MatchRoot("bar"),
			inPath: path.Root("foo"),
			exceptedValues: []attr.Value{
				types.StringValue("bar value"),
			},
			patterns:        []string{"baz"},
			message:         "message value",
			expError:        true,
			expErrorMessage: "message value",
		},
		"self-is-null": {
			req: validators.PatternsIfAttributeIsOneOfRequest{
				ConfigValue:    types.StringNull(),
				Path:           path.Root("foo"),
				PathExpression: path.MatchRoot("foo"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.StringAttribute{},
							"bar": schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.String,
							"bar": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.String, nil),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
					}),
				},
			},
			in:     path.MatchRoot("bar"),
			inPath: path.Root("foo"),
			exceptedValues: []attr.Value{
				types.StringValue("bar value"),
			},
			patterns:        []string{"baz"},
			message:         "",
			expError:        true,
			expErrorMessage: `is empty, value must match expression patterns 'baz'`,
		},
		"self-is-unknown": {
			req: validators.PatternsIfAttributeIsOneOfRequest{
				ConfigValue:    types.StringUnknown(),
				Path:           path.Root("foo"),
				PathExpression: path.MatchRoot("foo"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.StringAttribute{},
							"bar": schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.String,
							"bar": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.String, nil),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
					}),
				},
			},
			in:     path.MatchRoot("bar"),
			inPath: path.Root("foo"),
			exceptedValues: []attr.Value{
				types.StringValue("bar value"),
			},
			patterns: []string{"baz"},
			message:  "",
			expError: false,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &validators.PatternsIfAttributeIsOneOfResponse{}

			validators.PatternsIfAttributeIsOneOf(test.in, test.exceptedValues, test.patterns, test.message).Validate(context.TODO(), test.req, resp)

			if test.expError && resp.Diagnostics.HasError() {
				d1 := validatordiag.InvalidAttributeValueDiagnostic(test.inPath, test.expErrorMessage, test.req.ConfigValue.ValueString())
				d2 := validatordiag.InvalidAttributeValueMatchDiagnostic(test.inPath, test.expErrorMessage, test.req.ConfigValue.ValueString())

				if !resp.Diagnostics.Contains(d1) && !resp.Diagnostics.Contains(d2) {
					t.Fatalf("expected error(s) to contain (%s), got none. Error message is: (%s)", test.expErrorMessage, resp.Diagnostics.Errors())
				}
			}

			if !test.expError && resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error(s): %s", resp)
			}

			if test.expError && !resp.Diagnostics.HasError() {
				t.Fatal("expected error(s), got none")
			}
		})
	}
}
