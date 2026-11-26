// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/provider"
)

func TestUnit_Provider_FabricItemSensitivityLabelSchemas(t *testing.T) {
	ctx := t.Context()
	prov := provider.New("testUnit")
	itemPackage := reflect.TypeFor[fabricitem.ResourceFabricItem]().PkgPath()

	t.Run("resources", func(t *testing.T) {
		count := 0

		for _, factory := range prov.Resources(ctx) {
			item := factory()
			if reflect.TypeOf(item).Elem().PkgPath() != itemPackage {
				continue
			}

			count++

			var metadata resource.MetadataResponse
			item.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: common.ProviderTypeName}, &metadata)

			t.Run(metadata.TypeName, func(t *testing.T) {
				var resp resource.SchemaResponse
				item.Schema(ctx, resource.SchemaRequest{}, &resp)
				require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
				require.Contains(t, resp.Schema.Attributes, "tags")

				require.NotContains(t, resp.Schema.Attributes, "sensitivity_label_settings")

				label, ok := resp.Schema.Attributes["sensitivity_label"].(resourceschema.SingleNestedAttribute)
				require.True(t, ok, "missing sensitivity label")
				require.True(t, label.Optional)
				require.True(t, label.Computed)
				require.NotEmpty(t, label.PlanModifiers)
				require.Len(t, label.Attributes, 1)
				require.NotEmpty(t, label.Validators)
				require.True(t, label.Attributes["id"].IsOptional())
				require.True(t, label.Attributes["id"].IsComputed())
				require.Equal(t, customtypes.UUIDType{}, label.Attributes["id"].GetType())
			})
		}

		require.NotZero(t, count)
		t.Logf("Checked %d Fabric item resources", count)
	})

	t.Run("data_sources", func(t *testing.T) {
		count := 0

		for _, factory := range prov.DataSources(ctx) {
			item := factory()
			if reflect.TypeOf(item).Elem().PkgPath() != itemPackage {
				continue
			}

			count++

			var metadata datasource.MetadataResponse
			item.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: common.ProviderTypeName}, &metadata)

			t.Run(metadata.TypeName, func(t *testing.T) {
				var resp datasource.SchemaResponse
				item.Schema(ctx, datasource.SchemaRequest{}, &resp)
				require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)

				attributes := resp.Schema.Attributes
				if values, exists := attributes["values"]; exists {
					nested, ok := values.(datasourceschema.SetNestedAttribute)
					require.True(t, ok, "unexpected item list schema")

					attributes = nested.NestedObject.Attributes
				}

				require.Contains(t, attributes, "tags")

				label, ok := attributes["sensitivity_label"].(datasourceschema.SingleNestedAttribute)
				require.True(t, ok, "missing sensitivity label")
				require.True(t, label.Computed)
				require.NotContains(t, label.Attributes, "label_id")
				require.True(t, label.Attributes["id"].IsComputed())
				require.Equal(t, customtypes.UUIDType{}, label.Attributes["id"].GetType())
			})
		}

		require.NotZero(t, count)
		t.Logf("Checked %d Fabric item data sources", count)
	})
}
