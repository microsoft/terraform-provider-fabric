// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceonelakeimmutability_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/services/workspaceonelakeimmutability"
)

func TestUnit_WorkspaceOneLakeImmutabilityPolicy_Configure(t *testing.T) {
	resourceWithConfigure := workspaceonelakeimmutability.NewResourceWorkspaceOneLakeImmutabilityPolicy().(resource.ResourceWithConfigure)

	var resourceMissingDataResponse resource.ConfigureResponse
	resourceWithConfigure.Configure(context.Background(), resource.ConfigureRequest{}, &resourceMissingDataResponse)

	if resourceMissingDataResponse.Diagnostics.HasError() {
		t.Fatalf("unexpected resource configure diagnostic with no provider data: %v", resourceMissingDataResponse.Diagnostics)
	}

	var resourceInvalidDataResponse resource.ConfigureResponse
	resourceWithConfigure.Configure(context.Background(), resource.ConfigureRequest{ProviderData: "invalid"}, &resourceInvalidDataResponse)

	if !resourceInvalidDataResponse.Diagnostics.HasError() {
		t.Fatal("expected a resource configure diagnostic for invalid provider data")
	}

	dataSourceWithConfigure := workspaceonelakeimmutability.NewDataSourceWorkspaceOneLakeImmutabilityPolicy().(datasource.DataSourceWithConfigure)

	var dataSourceMissingDataResponse datasource.ConfigureResponse
	dataSourceWithConfigure.Configure(context.Background(), datasource.ConfigureRequest{}, &dataSourceMissingDataResponse)

	if dataSourceMissingDataResponse.Diagnostics.HasError() {
		t.Fatalf("unexpected data source configure diagnostic with no provider data: %v", dataSourceMissingDataResponse.Diagnostics)
	}

	var dataSourceInvalidDataResponse datasource.ConfigureResponse
	dataSourceWithConfigure.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: "invalid"}, &dataSourceInvalidDataResponse)

	if !dataSourceInvalidDataResponse.Diagnostics.HasError() {
		t.Fatal("expected a data source configure diagnostic for invalid provider data")
	}
}
