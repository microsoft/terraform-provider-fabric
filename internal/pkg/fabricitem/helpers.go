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
	MaxRetries    int
	RetryInterval time.Duration
	Operation     string
}

// DefaultUpdateRetryConfig returns the default retry configuration for update operations
func DefaultUpdateRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:    7,
		RetryInterval: time.Minute,
		Operation:     "update",
	}
}

// RetryUpdateOperation executes an update operation with retry logic for handling "ItemDisplayNameNotAvailableYet" errors
func RetryUpdateOperation(ctx context.Context, config RetryConfig, operation func() error) error {
	var err error
	var errRespFabric *fabcore.ResponseError

	for i := 0; i < config.MaxRetries; i++ {
		err = operation()

		if err == nil {
			return nil
		}

		if errors.As(err, &errRespFabric) && errRespFabric.ErrorCode == "ItemDisplayNameNotAvailableYet" {
			if i < config.MaxRetries-1 {
				tflog.Debug(ctx, fmt.Sprintf("Retry %d/%d failed with ItemDisplayNameNotAvailableYet, retrying in %v...", i+1, config.MaxRetries, config.RetryInterval))
				time.Sleep(config.RetryInterval)
				continue
			} else {
				tflog.Error(ctx, fmt.Sprintf("Max retries (%d) reached for %s operation. Operation failed.", config.MaxRetries, config.Operation))
			}
		}

		break
	}

	return err
}

// RetryUpdateItem executes UpdateItem with retry logic and returns the response
func RetryUpdateItem(ctx context.Context, client *fabcore.ItemsClient, workspaceID, itemID string, request fabcore.UpdateItemRequest) (fabcore.ItemsClientUpdateItemResponse, error) {
	var respUpdate fabcore.ItemsClientUpdateItemResponse
	var err error

	err = RetryUpdateOperation(ctx, DefaultUpdateRetryConfig(), func() error {
		respUpdate, err = client.UpdateItem(ctx, workspaceID, itemID, request, nil)
		return err
	})

	return respUpdate, err
}
