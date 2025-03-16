// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

func GetDataSourceSPNSupportNote(md string, spn bool) string { //revive:disable-line:flag-parameter
	if spn {
		return md + SPNSupportedDataSource
	}

	return md + SPNNotSupportedDataSource
}

func GetResourceSPNSupportNote(md string, spn bool) string { //revive:disable-line:flag-parameter
	if spn {
		return md + SPNSupportedResource
	}

	return md + SPNNotSupportedResource
}

func GetDataSourcePreviewNote(md string, preview bool) string { //revive:disable-line:flag-parameter
	if preview {
		return md + PreviewDataSource
	}

	return md
}

func GetResourcePreviewNote(md string, preview bool) string { //revive:disable-line:flag-parameter
	if preview {
		return md + PreviewResource
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
