// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	FabricSDKLoggerName   = "fabric-sdk-go"
	AzureSDKLoggingEnvVar = "AZURE_SDK_GO_LOGGING"
	AzureSDKLoggingAll    = "all"
)

func NewFabricSDKLoggerSubsystem(ctx context.Context) (context.Context, hclog.Level, error) {
	targetLevel := hclog.LevelFromString(os.Getenv("FABRIC_SDK_GO_LOGGING"))

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
