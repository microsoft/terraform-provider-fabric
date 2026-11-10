// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package monitoring_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/services/monitoring"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

func TestUnit_MonitoringResource_DisplayName(t *testing.T) {
	t.Parallel()
	var response frameworkresource.SchemaResponse
	monitoring.NewResourceMonitoring().Schema(t.Context(), frameworkresource.SchemaRequest{}, &response)
	attribute := response.Schema.Attributes["display_name"].(schema.StringAttribute)

	for _, tc := range []struct {
		name  string
		valid bool
	}{
		{"ab", true},
		{"_a", true},
		{"a_", true},
		{"a__b", true},
		{"a.b-c_d", true},
		{strings.Repeat("a", 128), true},
		{"", false},
		{"a", false},
		{"a b", false},
		{"__a", false},
		{"a__", false},
		{"a\nb", false},
		{"caf\u00e9", false},
		{strings.Repeat("a", 129), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			hasError := false

			for _, v := range attribute.Validators {
				var validation validator.StringResponse
				v.ValidateString(t.Context(), validator.StringRequest{
					Path: path.Root("display_name"), ConfigValue: types.StringValue(tc.name),
				}, &validation)
				hasError = hasError || validation.Diagnostics.HasError()
			}

			require.Equal(t, !tc.valid, hasError)
		})
	}
}

func TestUnit_MonitoringDataSources_MissingProperties(t *testing.T) {
	server := fakes.NewFakeServer()
	entity := fakes.NewRandomMonitoring()
	entity.Properties = nil
	server.Upsert(entity)
	// Each lookup gets its own step because the SDK fake shares pager state by request path.
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, server.ServerFactory, nil, []resource.TestStep{
		{
			Config: fmt.Sprintf(`
data "fabric_monitoring" "by_id" {
  workspace_id = %q
  id           = %q
}
`, *entity.WorkspaceID, *entity.ID),
			Check: resource.TestCheckNoResourceAttr("data.fabric_monitoring.by_id", "properties.data_collection_state"),
		},
		{
			Config: fmt.Sprintf(`
data "fabric_monitoring" "by_name" {
  workspace_id = %q
  display_name = %q
}
`, *entity.WorkspaceID, *entity.DisplayName),
			Check: resource.TestCheckNoResourceAttr("data.fabric_monitoring.by_name", "properties.data_collection_state"),
		},
		{
			Config: fmt.Sprintf(`
data "fabric_monitoring_items" "all" {
  workspace_id = %q
}
`, *entity.WorkspaceID),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.fabric_monitoring_items.all", "values.#", "1"),
				resource.TestCheckNoResourceAttr("data.fabric_monitoring_items.all", "values.0.properties.data_collection_state"),
			),
		},
	}))
}
