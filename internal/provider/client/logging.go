// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	FabricSDKLoggerName                  = "fabric-sdk-go"
	AzureSDKLoggingEnvVar                = "AZURE_SDK_GO_LOGGING"
	AzureSDKLoggingAll                   = "all"
	FabricSDKLoggingEnvVar               = "FABRIC_SDK_GO_LOGGING"
	FabricSDKLoggingIncludeBodyEnvVar    = "FABRIC_SDK_GO_LOGGING_INCLUDE_BODY"
	FabricSDKLoggingAllowedHeadersEnvVar = "FABRIC_SDK_GO_LOGGING_ALLOWED_HEADERS"
)

// NewFabricSDKLoggerSubsystem initializes the logger subsystem for the Fabric SDK.
func NewFabricSDKLoggerSubsystem(ctx context.Context) (context.Context, hclog.Level, error) {
	targetLevel := hclog.LevelFromString(os.Getenv(FabricSDKLoggingEnvVar))

	// If the level is not set, or is set to "off", disable logging
	if targetLevel == hclog.NoLevel {
		targetLevel = hclog.Off
	}

	// Enable azcore logging if the target level is not "off"
	if targetLevel != hclog.Off {
		if err := os.Setenv(AzureSDKLoggingEnvVar, AzureSDKLoggingAll); err != nil {
			return ctx, targetLevel, err
		}
	}

	ctx = tflog.NewSubsystem(ctx, FabricSDKLoggerName, tflog.WithLevel(targetLevel))
	ctx = tflog.SubsystemMaskFieldValuesWithFieldKeys(ctx, FabricSDKLoggerName, "Authorization")

	return ctx, targetLevel, nil
}

// GetLoggingIncludeBodyOption determines if request/response bodies should be included in logs.
func GetLoggingIncludeBodyOption(ctx context.Context) (bool, error) {
	if includeBodyEnv, ok := os.LookupEnv(FabricSDKLoggingIncludeBodyEnvVar); ok && includeBodyEnv != "" {
		includeBodyBool, err := strconv.ParseBool(includeBodyEnv)
		if err != nil {
			tflog.Error(ctx, "Failed to parse include body option", map[string]any{
				"env_var": FabricSDKLoggingIncludeBodyEnvVar,
				"value":   includeBodyEnv,
				"error":   err.Error(),
			})

			return false, err
		}

		return includeBodyBool, nil
	}

	return false, nil
}

// GetLoggingAllowedHeadersOption determines which HTTP headers should be included in logs.
func GetLoggingAllowedHeadersOption(ctx context.Context) ([]string, error) {
	if allowedHeadersEnv, ok := os.LookupEnv(FabricSDKLoggingAllowedHeadersEnvVar); ok && allowedHeadersEnv != "" {
		headers := strings.Split(allowedHeadersEnv, ";")
		validHeaders := make([]string, 0, len(headers))

		// Simple validation to prevent injection: only accept standard header format
		headerRegex := regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)

		for _, header := range headers {
			h := strings.ToLower(strings.TrimSpace(header))
			if h != "" && headerRegex.MatchString(h) {
				validHeaders = append(validHeaders, h)
			} else if h != "" {
				tflog.Warn(ctx, "Skipping invalid header format", map[string]any{
					"header": header,
				})
			}
		}

		if len(validHeaders) > 0 {
			tflog.Debug(ctx, "Using custom allowed headers", map[string]any{
				"headers": validHeaders,
			})

			return validHeaders, nil
		}
	}

	return nil, nil
}

// ConfigureLoggingOptions configures the logging options for the Fabric SDK based on the environment variables.
func ConfigureLoggingOptions(ctx context.Context, logLevel hclog.Level) (*policy.LogOptions, error) {
	if logLevel == hclog.Off {
		return nil, nil //nolint:nilnil
	}

	logOptions := &policy.LogOptions{}

	includeBody, err := GetLoggingIncludeBodyOption(ctx)
	if err != nil {
		return nil, err
	}

	logOptions.IncludeBody = includeBody

	allowedHeaders, err := GetLoggingAllowedHeadersOption(ctx)
	if err != nil {
		return nil, err
	}

	if len(allowedHeaders) > 0 {
		logOptions.AllowedHeaders = allowedHeaders
	}

	return logOptions, nil
}
