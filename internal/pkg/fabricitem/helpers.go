// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

func NewResourceMarkdownDescription(typeInfo tftypeinfo.TFTypeInfo, plural bool) string { //revive:disable-line:flag-parameter
	md := fmt.Sprintf("The %s resource allows you to manage a Fabric", typeInfo.Name)

	if plural {
		md = fmt.Sprintf("The %s resource allows you to manage a Fabric", typeInfo.Names)
	}

	if typeInfo.DocsURL != "" {
		name := typeInfo.Name
		if plural {
			name = typeInfo.Names
		}

		md += fmt.Sprintf(" [%s](%s).", name, typeInfo.DocsURL)
	} else {
		md += fmt.Sprintf(" %s.", typeInfo.Name)
	}

	if typeInfo.IsSPNSupported {
		md += SPNSupportedResource
	} else {
		md += SPNNotSupportedResource
	}

	if typeInfo.IsPreview {
		md += PreviewResource
	}

	return md
}

func NewDataSourceMarkdownDescription(typeInfo tftypeinfo.TFTypeInfo, plural bool) string { //revive:disable-line:flag-parameter
	md := fmt.Sprintf("The %s data-source allows you to retrieve details about a Fabric", typeInfo.Name)

	if plural {
		md = fmt.Sprintf("The %s data-source allows you to retrieve a list of Fabric", typeInfo.Names)
	}

	if typeInfo.DocsURL != "" {
		name := typeInfo.Name
		if plural {
			name = typeInfo.Names
		}

		md += fmt.Sprintf(" [%s](%s).", name, typeInfo.DocsURL)
	} else {
		md += fmt.Sprintf(" %s.", typeInfo.Name)
	}

	if typeInfo.IsSPNSupported {
		md += SPNSupportedDataSource
	} else {
		md += SPNNotSupportedDataSource
	}

	if typeInfo.IsPreview {
		md += PreviewDataSource
	}

	return md
}

func NewEphemeralResourceMarkdownDescription(typeInfo tftypeinfo.TFTypeInfo, plural bool) string { //revive:disable-line:flag-parameter
	md := fmt.Sprintf("The %s ephemeral resource allows you to manage a temporary Fabric", typeInfo.Name)

	if plural {
		md = fmt.Sprintf("The %s ephemeral resources allow you to manage temporary Fabric", typeInfo.Names)
	}

	if typeInfo.DocsURL != "" {
		name := typeInfo.Name
		if plural {
			name = typeInfo.Names
		}

		md += fmt.Sprintf(" [%s](%s).", name, typeInfo.DocsURL)
	} else {
		md += fmt.Sprintf(" %s.", typeInfo.Name)
	}

	md += "\n\n-> Ephemeral Resources are supported in HashiCorp Terraform version 1.11 and later."

	if typeInfo.IsSPNSupported {
		md += SPNSupportedResource
	} else {
		md += SPNNotSupportedResource
	}

	if typeInfo.IsPreview {
		md += PreviewEphemeralResource
	}

	return md
}

func IsPreviewMode(name string, itemIsPreview, providerPreviewMode bool) diag.Diagnostics { //revive:disable-line:flag-parameter
	var diags diag.Diagnostics

	if itemIsPreview && !providerPreviewMode {
		diags.AddError(
			common.ErrorPreviewModeHeader,
			fmt.Sprintf(common.ErrorPreviewModeDetails, name),
		)

		return diags
	}

	if itemIsPreview && providerPreviewMode {
		diags.AddWarning(
			fmt.Sprintf(common.WarningPreviewModeHeader, name),
			fmt.Sprintf(common.WarningPreviewModeDetails, name),
		)

		return diags
	}

	return nil
}

// RetryConfig holds configuration for retry operations
type RetryConfig struct {
	RetryInterval time.Duration
	Operation     string
}

// RetryOperation executes any operation with retry logic for handling "ItemDisplayNameNotAvailableYet" errors
// This will retry indefinitely until the operation succeeds or encounters a non-retryable error
func RetryOperation(ctx context.Context, config RetryConfig, operation func() error) error {
	var err error
	var errRespFabric *fabcore.ResponseError
	retryCount := 0

	for {
		err = operation()

		if err == nil {
			if retryCount > 0 {
				tflog.Debug(ctx, fmt.Sprintf("Operation succeeded after %d retries", retryCount))
			}
			return nil
		}

		// Check if context is cancelled
		if ctx.Err() != nil {
			tflog.Error(ctx, fmt.Sprintf("Context cancelled during %s operation after %d retries", config.Operation, retryCount))
			return ctx.Err()
		}

		if errors.As(err, &errRespFabric) && errRespFabric.ErrorCode == "ItemDisplayNameNotAvailableYet" {
			retryCount++
			tflog.Debug(ctx, fmt.Sprintf("Retry %d failed with ItemDisplayNameNotAvailableYet, retrying in %v...", retryCount, config.RetryInterval))

			// Use a timer to respect context cancellation during sleep
			timer := time.NewTimer(config.RetryInterval)
			select {
			case <-ctx.Done():
				timer.Stop()
				tflog.Error(ctx, fmt.Sprintf("Context cancelled during %s operation after %d retries", config.Operation, retryCount))
				return ctx.Err()
			case <-timer.C:
				continue
			}
		}

		// For any other error, don't retry
		tflog.Error(ctx, fmt.Sprintf("Non-retryable error in %s operation after %d retries: %v", config.Operation, retryCount, err))
		break
	}

	return err
}

func DefaultUpdateRetryConfig() RetryConfig {
	return RetryConfig{
		RetryInterval: 2 * time.Minute,
		Operation:     "update",
	}
}

func RetryUpdateItem(ctx context.Context, client *fabcore.ItemsClient, workspaceID, itemID string, request fabcore.UpdateItemRequest) (fabcore.ItemsClientUpdateItemResponse, error) {
	var respUpdate fabcore.ItemsClientUpdateItemResponse
	var err error

	err = RetryOperation(ctx, DefaultUpdateRetryConfig(), func() error {
		respUpdate, err = client.UpdateItem(ctx, workspaceID, itemID, request, nil)
		return err
	})

	return respUpdate, err
}

func DefaultCreateRetryConfig() RetryConfig {
	return RetryConfig{
		RetryInterval: 2 * time.Minute,
		Operation:     "create",
	}
}

func RetryCreateItem(ctx context.Context, client *fabcore.ItemsClient, workspaceID string, request fabcore.CreateItemRequest) (fabcore.ItemsClientCreateItemResponse, error) {
	var respCreate fabcore.ItemsClientCreateItemResponse
	var err error

	err = RetryOperation(ctx, DefaultCreateRetryConfig(), func() error {
		respCreate, err = client.CreateItem(ctx, workspaceID, request, nil)
		return err
	})

	return respCreate, err
}
