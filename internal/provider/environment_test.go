// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"context"
	"errors"
	"maps"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	tfprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/microsoft/fabric-sdk-go/fabric"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/provider"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

func clearProviderEnvironment(t *testing.T) {
	t.Helper()

	for _, env := range os.Environ() {
		key, _, _ := strings.Cut(env, "=")
		if strings.HasPrefix(key, "FABRIC_") || strings.HasPrefix(key, "ARM_") {
			t.Setenv(key, "")
			require.NoError(t, os.Unsetenv(key))
		}
	}
}

func providerCloudConfig(t *testing.T, prov tfprovider.Provider, values map[string]tftypes.Value) tfsdk.Config {
	t.Helper()

	var schemaResp tfprovider.SchemaResponse
	prov.Schema(t.Context(), tfprovider.SchemaRequest{}, &schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError(), "%v", schemaResp.Diagnostics)

	objectType, ok := schemaResp.Schema.Type().TerraformType(t.Context()).(tftypes.Object)
	require.True(t, ok)

	fields := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attrType := range objectType.AttributeTypes {
		fields[name] = tftypes.NewValue(attrType, nil)
	}

	maps.Copy(fields, values)

	return tfsdk.Config{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objectType, fields)}
}

func TestUnit_Provider_Environment(t *testing.T) {
	testCases := []struct {
		name        string
		values      map[string]tftypes.Value
		env         map[string]string
		environment fabric.Environment
		endpoint    string
		errorText   string
	}{
		{name: "default Public"},
		{
			name:   "explicit Public",
			values: map[string]tftypes.Value{"environment": tftypes.NewValue(tftypes.String, "public")},
		},
		{
			name:        "GCC High",
			values:      map[string]tftypes.Value{"environment": tftypes.NewValue(tftypes.String, "gcchigh")},
			environment: fabric.EnvironmentGCCHigh,
		},
		{
			name:        "mixed case",
			values:      map[string]tftypes.Value{"environment": tftypes.NewValue(tftypes.String, "GcCHigh")},
			environment: fabric.EnvironmentGCCHigh,
		},
		{
			name:        "Fabric environment variable",
			env:         map[string]string{"FABRIC_ENVIRONMENT": " GCCHIGH "},
			environment: fabric.EnvironmentGCCHigh,
		},
		{
			name:        "legacy environment variable",
			env:         map[string]string{"ARM_ENVIRONMENT": "gcchigh"},
			environment: fabric.EnvironmentGCCHigh,
		},
		{
			name:   "provider argument takes precedence",
			values: map[string]tftypes.Value{"environment": tftypes.NewValue(tftypes.String, "public")},
			env:    map[string]string{"FABRIC_ENVIRONMENT": "gcchigh", "ARM_ENVIRONMENT": "usgovernment"},
		},
		{
			name:        "Fabric variable takes precedence over Azure authority",
			env:         map[string]string{"FABRIC_ENVIRONMENT": "gcchigh", "ARM_ENVIRONMENT": "usgovernment"},
			environment: fabric.EnvironmentGCCHigh,
		},
		{
			name: "endpoint override",
			values: map[string]tftypes.Value{
				"environment": tftypes.NewValue(tftypes.String, "gcchigh"),
				"endpoint":    tftypes.NewValue(tftypes.String, "https://custom.example.com"),
			},
			env:         map[string]string{"FABRIC_ENDPOINT": "https://ignored.example.com"},
			environment: fabric.EnvironmentGCCHigh,
			endpoint:    "https://custom.example.com",
		},
		{
			name:        "endpoint environment variable",
			env:         map[string]string{"FABRIC_ENVIRONMENT": "gcchigh", "FABRIC_ENDPOINT": "https://custom.example.com"},
			environment: fabric.EnvironmentGCCHigh,
			endpoint:    "https://custom.example.com",
		},
		{
			name:      "legacy government authority is ambiguous",
			env:       map[string]string{"ARM_ENVIRONMENT": "usgovernment"},
			errorText: `Unsupported Fabric environment "usgovernment". Supported values are 'public' and 'gcchigh'.`,
		},
		{
			name:      "unsupported China",
			env:       map[string]string{"FABRIC_ENVIRONMENT": "china"},
			errorText: `Unsupported Fabric environment "china". Supported values are 'public' and 'gcchigh'.`,
		},
		{
			name:      "unsupported DoD",
			env:       map[string]string{"FABRIC_ENVIRONMENT": "dod"},
			errorText: "Unsupported Fabric environment",
		},
		{
			name:      "unknown name",
			env:       map[string]string{"FABRIC_ENVIRONMENT": "typo"},
			errorText: "Unsupported Fabric environment",
		},
		{
			name:      "usgov was never supported",
			values:    map[string]tftypes.Value{"environment": tftypes.NewValue(tftypes.String, "usgov")},
			errorText: "Unsupported Fabric environment",
		},
		{
			name:      "unknown Terraform environment",
			values:    map[string]tftypes.Value{"environment": tftypes.NewValue(tftypes.String, tftypes.UnknownValue)},
			errorText: "unknown configuration value for the Fabric environment",
		},
		{
			name:      "unknown Terraform endpoint",
			values:    map[string]tftypes.Value{"endpoint": tftypes.NewValue(tftypes.String, tftypes.UnknownValue)},
			errorText: "unknown configuration value for the Microsoft Fabric API endpoint",
		},
		{
			name:      "invalid endpoint environment variable",
			env:       map[string]string{"FABRIC_ENDPOINT": "not-a-url"},
			errorText: "not valid absolute URL",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			clearProviderEnvironment(t)

			for key, value := range testCase.env {
				t.Setenv(key, value)
			}

			prov := provider.New("testUnit")
			var configured *pconfig.ProviderConfig

			prov.ConfigureAffirmProviderConfig(func(cfg *pconfig.ProviderConfig) {
				configured = cfg
			})

			var resp tfprovider.ConfigureResponse
			prov.Configure(t.Context(), tfprovider.ConfigureRequest{Config: providerCloudConfig(t, prov, testCase.values)}, &resp)

			if testCase.errorText != "" {
				require.True(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
				require.Contains(t, resp.Diagnostics.Errors()[0].Detail(), testCase.errorText)
				require.Nil(t, configured, "invalid configuration must not create a client")
				require.Nil(t, resp.ResourceData)

				return
			}

			require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
			require.NotNil(t, configured)
			require.Equal(t, testCase.environment, configured.Environment)

			authority := "https://login.microsoftonline.com/"
			endpoint := "https://api.fabric.microsoft.com"

			if testCase.environment == fabric.EnvironmentGCCHigh {
				authority = "https://login.microsoftonline.us/"
				endpoint = "https://highapi.fabric.microsoft.us"
			}

			if testCase.endpoint != "" {
				endpoint = testCase.endpoint
			}

			require.Equal(t, authority, configured.Auth.Environment.ActiveDirectoryAuthorityHost)
			require.Equal(t, endpoint, configured.Endpoint)
			require.Equal(t, endpoint, configured.FabricClient.Endpoint)
			require.Same(t, configured.ProviderData, resp.ResourceData)
			require.Same(t, resp.ResourceData, resp.DataSourceData)
			require.Same(t, resp.ResourceData, resp.EphemeralResourceData)
		})
	}
}

func TestUnit_Provider_ClientCreationError(t *testing.T) {
	clearProviderEnvironment(t)

	prov := provider.New("testUnit")
	prov.ConfigureCreateClient(func(_ context.Context, _ *pconfig.ProviderConfig) (*fabric.Client, error) {
		return nil, errors.New("test client creation failed")
	})

	var resp tfprovider.ConfigureResponse
	prov.Configure(t.Context(), tfprovider.ConfigureRequest{Config: providerCloudConfig(t, prov, nil)}, &resp)
	require.True(t, resp.Diagnostics.HasError())
	require.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "test client creation failed")
	require.Nil(t, resp.ResourceData)
}

func TestUnit_Provider_EnvironmentSchema(t *testing.T) {
	var resp tfprovider.SchemaResponse
	provider.New("testUnit").Schema(t.Context(), tfprovider.SchemaRequest{}, &resp)
	require.False(t, resp.Diagnostics.HasError())

	attribute, ok := resp.Schema.Attributes["environment"].(schema.StringAttribute)
	require.True(t, ok)

	for _, value := range []string{"public", "PUBLIC", "gcchigh", "GcCHigh", "usgovernment", "usgov", "china", "dod", "typo", ""} {
		t.Run(value, func(t *testing.T) {
			valid := strings.EqualFold(value, "public") || strings.EqualFold(value, "gcchigh")

			for _, attributeValidator := range attribute.Validators {
				var validationResp validator.StringResponse
				attributeValidator.ValidateString(t.Context(), validator.StringRequest{
					Path:        path.Root("environment"),
					ConfigValue: types.StringValue(value),
				}, &validationResp)
				require.Equal(t, !valid, validationResp.Diagnostics.HasError(), "%v", validationResp.Diagnostics)
			}
		})
	}
}
