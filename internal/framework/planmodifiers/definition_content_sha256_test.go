// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package planmodifiers_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/planmodifiers"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/params"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
)

func testDefinitionPartSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true},
			"description": schema.StringAttribute{Optional: true},
			"part": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"source":          schema.StringAttribute{Required: true},
					"processing_mode": schema.StringAttribute{Optional: true, Computed: true},
					"parameters": schema.SetNestedAttribute{
						Optional:   true,
						CustomType: supertypes.NewSetNestedObjectTypeOf[params.ParametersModel](ctx),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type":  schema.StringAttribute{Required: true},
								"find":  schema.StringAttribute{Required: true},
								"value": schema.StringAttribute{Required: true},
							},
						},
					},
					"tokens": schema.MapAttribute{
						Optional:    true,
						CustomType:  supertypes.MapTypeOf[types.String]{MapType: types.MapType{ElemType: types.StringType}},
						ElementType: types.StringType,
					},
					"tokens_delimiter":      schema.StringAttribute{Optional: true, Computed: true},
					"source_content_sha256": schema.StringAttribute{Computed: true},
				},
			},
		},
	}
}

func testDefinitionPartRaw(ctx context.Context, t *testing.T, id, description, source, sha tftypes.Value) tftypes.Value {
	t.Helper()

	rootType, ok := testDefinitionPartSchema(ctx).Type().TerraformType(ctx).(tftypes.Object)
	require.True(t, ok)

	partType, ok := rootType.AttributeTypes["part"].(tftypes.Object)
	require.True(t, ok)

	part := tftypes.NewValue(partType, map[string]tftypes.Value{
		"source":                source,
		"processing_mode":       tftypes.NewValue(tftypes.String, transforms.ProcessingModeGoTemplate),
		"parameters":            tftypes.NewValue(partType.AttributeTypes["parameters"], nil),
		"tokens":                tftypes.NewValue(partType.AttributeTypes["tokens"], nil),
		"tokens_delimiter":      tftypes.NewValue(tftypes.String, transforms.TokensDelimiterCurlyBraces),
		"source_content_sha256": sha,
	})

	return tftypes.NewValue(rootType, map[string]tftypes.Value{
		"id":          id,
		"description": description,
		"part":        part,
	})
}

func TestUnit_DefinitionContentSha256(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	sourcePath := filepath.Join(t.TempDir(), "definition.json")
	require.NoError(t, os.WriteFile(sourcePath, []byte(`{"name":"dev"}`), 0o600))

	pm := planmodifiers.DefinitionContentSha256(
		path.MatchRelative().AtParent().AtName("source"),
		path.MatchRelative().AtParent().AtName("processing_mode"),
		path.MatchRelative().AtParent().AtName("tokens"),
		path.MatchRelative().AtParent().AtName("parameters"),
		path.MatchRelative().AtParent().AtName("tokens_delimiter"),
	)

	knownDescription := tftypes.NewValue(tftypes.String, "example")
	nullValue := tftypes.NewValue(tftypes.String, nil)
	unknownValue := tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	knownSource := tftypes.NewValue(tftypes.String, sourcePath)

	type testCase struct {
		description tftypes.Value
		source      tftypes.Value
		expUnknown  bool
	}

	testCases := map[string]testCase{
		// `id` is unknown in the plan, so a known hash here also pins the check to the config.
		"config fully known": {
			description: knownDescription,
			source:      knownSource,
			expUnknown:  false,
		},
		// The resource is deferred, so another resource may still rewrite the file before apply.
		"config not fully known": {
			description: unknownValue,
			source:      knownSource,
			expUnknown:  true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.StringRequest{
				Path:           path.Root("part").AtName("source_content_sha256"),
				PathExpression: path.MatchRoot("part").AtName("source_content_sha256"),
				Config: tfsdk.Config{
					Schema: testDefinitionPartSchema(ctx),
					Raw:    testDefinitionPartRaw(ctx, t, nullValue, testCase.description, testCase.source, nullValue),
				},
				Plan: tfsdk.Plan{
					Schema: testDefinitionPartSchema(ctx),
					Raw:    testDefinitionPartRaw(ctx, t, unknownValue, testCase.description, testCase.source, unknownValue),
				},
				PlanValue: types.StringUnknown(),
			}

			resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}

			pm.PlanModifyString(ctx, req, resp)

			require.False(t, resp.Diagnostics.HasError())
			require.Equal(t, testCase.expUnknown, resp.PlanValue.IsUnknown())

			if !testCase.expUnknown {
				require.Regexp(t, "^[0-9a-f]{64}$", resp.PlanValue.ValueString())
			}
		})
	}
}
