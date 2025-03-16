// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package client_test

import (
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pclient "github.com/microsoft/terraform-provider-fabric/internal/provider/client"
)

func TestUnit_GetLoggingIncludeBodyOption(t *testing.T) {
	ctx := t.Context()

	testCases := []struct {
		name           string
		envValue       string
		expectedResult bool
		expectError    bool
	}{
		{
			name:           "No environment variable set",
			envValue:       "",
			expectedResult: false,
			expectError:    false,
		},
		{
			name:           "Environment variable set to true",
			envValue:       "true",
			expectedResult: true,
			expectError:    false,
		},
		{
			name:           "Environment variable set to false",
			envValue:       "false",
			expectedResult: false,
			expectError:    false,
		},
		{
			name:           "Environment variable set to 1",
			envValue:       "1",
			expectedResult: true,
			expectError:    false,
		},
		{
			name:           "Environment variable set to 0",
			envValue:       "0",
			expectedResult: false,
			expectError:    false,
		},
		{
			name:           "Invalid boolean value",
			envValue:       "not-a-bool",
			expectedResult: false,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			err := os.Unsetenv(pclient.FabricSDKLoggingIncludeBodyEnvVar)
			require.NoError(t, err)

			if tc.envValue != "" {
				t.Setenv(pclient.FabricSDKLoggingIncludeBodyEnvVar, tc.envValue)
			}

			// Execute
			result, err := pclient.GetLoggingIncludeBodyOption(ctx)

			// Verify
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestUnit_GetLoggingAllowedHeadersOption(t *testing.T) {
	ctx := t.Context()

	testCases := []struct {
		name           string
		envValue       string
		expectedResult []string
		expectError    bool
	}{
		{
			name:           "No environment variable set",
			envValue:       "",
			expectedResult: nil,
			expectError:    false,
		},
		{
			name:           "Single valid header",
			envValue:       "location",
			expectedResult: []string{"location"},
			expectError:    false,
		},
		{
			name:           "Multiple valid headers",
			envValue:       "location;requestid;x-ms-client-request-id",
			expectedResult: []string{"location", "requestid", "x-ms-client-request-id"},
			expectError:    false,
		},
		{
			name:           "Headers with whitespace",
			envValue:       "  location  ;  requestid  ",
			expectedResult: []string{"location", "requestid"},
			expectError:    false,
		},
		{
			name:           "Mixed valid and invalid headers",
			envValue:       "location;invalid@header;x-ms-client-request-id",
			expectedResult: []string{"location", "x-ms-client-request-id"},
			expectError:    false,
		},
		{
			name:           "All invalid headers",
			envValue:       "invalid@header;another<invalid>;!not-valid!",
			expectedResult: nil,
			expectError:    false,
		},
		{
			name:           "Empty header value",
			envValue:       ";",
			expectedResult: nil,
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			err := os.Unsetenv(pclient.FabricSDKLoggingAllowedHeadersEnvVar)
			require.NoError(t, err)

			if tc.envValue != "" {
				t.Setenv(pclient.FabricSDKLoggingAllowedHeadersEnvVar, tc.envValue)
			}

			// Execute
			result, err := pclient.GetLoggingAllowedHeadersOption(ctx)

			// Verify
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestUnit_GetLoggingAllowedQueryParamsOption(t *testing.T) {
	ctx := t.Context()

	testCases := []struct {
		name           string
		envValue       string
		expectedResult []string
		expectError    bool
	}{
		{
			name:           "No environment variable set",
			envValue:       "",
			expectedResult: nil,
			expectError:    false,
		},
		{
			name:           "Single valid query param",
			envValue:       "filter",
			expectedResult: []string{"filter"},
			expectError:    false,
		},
		{
			name:           "Multiple valid query params",
			envValue:       "filter;page;limit",
			expectedResult: []string{"filter", "page", "limit"},
			expectError:    false,
		},
		{
			name:           "Query params with whitespace",
			envValue:       "  filter  ;  page  ",
			expectedResult: []string{"filter", "page"},
			expectError:    false,
		},
		{
			name:           "Mixed valid and invalid query params",
			envValue:       "filter;invalid@param;limit",
			expectedResult: []string{"filter", "limit"},
			expectError:    false,
		},
		{
			name:           "All invalid query params",
			envValue:       "invalid@param;another<invalid>;!not-valid!",
			expectedResult: nil,
			expectError:    false,
		},
		{
			name:           "Empty query param value",
			envValue:       ";",
			expectedResult: nil,
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			err := os.Unsetenv(pclient.FabricSDKLoggingAllowedQueryParamsEnvVar)
			require.NoError(t, err)

			if tc.envValue != "" {
				t.Setenv(pclient.FabricSDKLoggingAllowedQueryParamsEnvVar, tc.envValue)
			}

			// Execute
			result, err := pclient.GetLoggingAllowedQueryParamsOption(ctx)

			// Verify
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestUnit_ConfigureLoggingOptions(t *testing.T) {
	ctx := t.Context()

	testCases := []struct {
		name                  string
		logLevel              hclog.Level
		includeBodyEnv        string
		allowedHeadersEnv     string
		allowedQueryParamsEnv string
		expectLogOptions      bool
		expectIncludeBody     bool
		expectHeaders         []string
		expectQueryParams     []string
		expectError           bool
	}{
		{
			name:             "Logging disabled",
			logLevel:         hclog.Off,
			expectLogOptions: false,
			expectError:      false,
		},
		{
			name:              "Logging enabled without options",
			logLevel:          hclog.Debug,
			expectLogOptions:  true,
			expectIncludeBody: false,
			expectHeaders:     nil,
			expectQueryParams: nil,
			expectError:       false,
		},
		{
			name:              "Logging with includeBody enabled",
			logLevel:          hclog.Debug,
			includeBodyEnv:    "true",
			expectLogOptions:  true,
			expectIncludeBody: true,
			expectHeaders:     nil,
			expectQueryParams: nil,
			expectError:       false,
		},
		{
			name:              "Logging with headers enabled",
			logLevel:          hclog.Debug,
			allowedHeadersEnv: "location;requestid",
			expectLogOptions:  true,
			expectIncludeBody: false,
			expectHeaders:     []string{"location", "requestid"},
			expectQueryParams: nil,
			expectError:       false,
		},
		{
			name:                  "Logging with query params enabled",
			logLevel:              hclog.Debug,
			allowedQueryParamsEnv: "filter;page",
			expectLogOptions:      true,
			expectIncludeBody:     false,
			expectHeaders:         nil,
			expectQueryParams:     []string{"filter", "page"},
			expectError:           false,
		},
		{
			name:                  "Logging with all options enabled",
			logLevel:              hclog.Debug,
			includeBodyEnv:        "true",
			allowedHeadersEnv:     "location;requestid",
			allowedQueryParamsEnv: "filter;page",
			expectLogOptions:      true,
			expectIncludeBody:     true,
			expectHeaders:         []string{"location", "requestid"},
			expectQueryParams:     []string{"filter", "page"},
			expectError:           false,
		},
		{
			name:             "Invalid includeBody value",
			logLevel:         hclog.Debug,
			includeBodyEnv:   "not-a-bool",
			expectLogOptions: false,
			expectError:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			err := os.Unsetenv(pclient.FabricSDKLoggingIncludeBodyEnvVar)
			require.NoError(t, err)

			err = os.Unsetenv(pclient.FabricSDKLoggingAllowedHeadersEnvVar)
			require.NoError(t, err)

			err = os.Unsetenv(pclient.FabricSDKLoggingAllowedQueryParamsEnvVar)
			require.NoError(t, err)

			// if err := os.Unsetenv(pclient.AzureSDKLoggingEnvVar); err != nil {
			// 	t.Fatalf("Failed to unset environment variable: %v", err)
			// }

			if tc.includeBodyEnv != "" {
				t.Setenv(pclient.FabricSDKLoggingIncludeBodyEnvVar, tc.includeBodyEnv)
			}

			if tc.allowedHeadersEnv != "" {
				t.Setenv(pclient.FabricSDKLoggingAllowedHeadersEnvVar, tc.allowedHeadersEnv)
			}

			if tc.allowedQueryParamsEnv != "" {
				t.Setenv(pclient.FabricSDKLoggingAllowedQueryParamsEnvVar, tc.allowedQueryParamsEnv)
			}

			// Execute
			result, err := pclient.ConfigureLoggingOptions(ctx, tc.logLevel)

			// Verify
			if tc.expectError {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)

				if tc.expectLogOptions {
					assert.NotNil(t, result)
					assert.Equal(t, tc.expectIncludeBody, result.IncludeBody)
					assert.Equal(t, tc.expectHeaders, result.AllowedHeaders)
					assert.Equal(t, tc.expectQueryParams, result.AllowedQueryParams)
				} else {
					assert.Nil(t, result)
				}
			}
		})
	}
}

func TestUnit_NewFabricSDKLoggerSubsystem(t *testing.T) {
	testCases := []struct {
		name          string
		logLevelEnv   string
		expectedLevel hclog.Level
	}{
		{
			name:          "No log level set",
			logLevelEnv:   "",
			expectedLevel: hclog.Off,
		},
		{
			name:          "Debug level set",
			logLevelEnv:   "debug",
			expectedLevel: hclog.Debug,
		},
		{
			name:          "Info level set",
			logLevelEnv:   "info",
			expectedLevel: hclog.Info,
		},
		{
			name:          "Warn level set",
			logLevelEnv:   "warn",
			expectedLevel: hclog.Warn,
		},
		{
			name:          "Error level set",
			logLevelEnv:   "error",
			expectedLevel: hclog.Error,
		},
		{
			name:          "Off level set",
			logLevelEnv:   "off",
			expectedLevel: hclog.Off,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			err := os.Unsetenv(pclient.FabricSDKLoggingEnvVar)
			require.NoError(t, err)

			ctx := t.Context()

			if tc.logLevelEnv != "" {
				t.Setenv(pclient.FabricSDKLoggingEnvVar, tc.logLevelEnv)
			}

			// Test that we get the expected Azure SDK logging env var when not off
			if tc.expectedLevel != hclog.Off {
				err := os.Unsetenv(pclient.AzureSDKLoggingEnvVar)
				require.NoError(t, err)
			}

			// Execute
			newCtx, level, err := pclient.NewFabricSDKLoggerSubsystem(ctx)

			// Verify
			require.NoError(t, err)
			assert.Equal(t, tc.expectedLevel, level)
			assert.NotNil(t, newCtx)

			// If logging is enabled, verify Azure SDK logging env var is set
			if tc.expectedLevel != hclog.Off {
				azLogEnv, exists := os.LookupEnv(pclient.AzureSDKLoggingEnvVar)
				assert.True(t, exists)
				assert.Equal(t, pclient.AzureSDKLoggingAll, azLogEnv)
			}
		})
	}
}
