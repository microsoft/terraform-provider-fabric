// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakediagnosticsetting_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/services/onelakediagnosticsetting"
)

func TestUnit_OneLakeDiagnosticSettingResource_Configure(t *testing.T) {
	testResource := onelakediagnosticsetting.NewResourceOneLakeDiagnosticSetting()
	configurable := testResource.(resource.ResourceWithConfigure)

	var nilProviderDataResp resource.ConfigureResponse
	configurable.Configure(context.Background(), resource.ConfigureRequest{}, &nilProviderDataResp)

	if nilProviderDataResp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics with nil provider data: %v", nilProviderDataResp.Diagnostics)
	}

	var invalidProviderDataResp resource.ConfigureResponse
	configurable.Configure(context.Background(), resource.ConfigureRequest{ProviderData: "invalid"}, &invalidProviderDataResp)

	if !invalidProviderDataResp.Diagnostics.HasError() {
		t.Fatal("expected diagnostics for invalid provider data")
	}

	if got := invalidProviderDataResp.Diagnostics.Errors()[0].Summary(); got != common.ErrorResourceConfigType {
		t.Fatalf("unexpected diagnostic summary: got %q, want %q", got, common.ErrorResourceConfigType)
	}
}

func TestUnit_OneLakeDiagnosticSettingDataSource_Configure(t *testing.T) {
	testDataSource := onelakediagnosticsetting.NewDataSourceOneLakeDiagnosticSetting()
	configurable := testDataSource.(datasource.DataSourceWithConfigure)

	var nilProviderDataResp datasource.ConfigureResponse
	configurable.Configure(context.Background(), datasource.ConfigureRequest{}, &nilProviderDataResp)

	if nilProviderDataResp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics with nil provider data: %v", nilProviderDataResp.Diagnostics)
	}

	var invalidProviderDataResp datasource.ConfigureResponse
	configurable.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: "invalid"}, &invalidProviderDataResp)

	if !invalidProviderDataResp.Diagnostics.HasError() {
		t.Fatal("expected diagnostics for invalid provider data")
	}

	if got := invalidProviderDataResp.Diagnostics.Errors()[0].Summary(); got != common.ErrorDataSourceConfigType {
		t.Fatalf("unexpected diagnostic summary: got %q, want %q", got, common.ErrorDataSourceConfigType)
	}
}
