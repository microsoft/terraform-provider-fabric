// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

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

func IsPreviewModeEnabled(name string, itemIsPreview, providerPreviewMode bool) diag.Diagnostics { //revive:disable-line:flag-parameter
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
			common.WarningPreviewModeHeader,
			fmt.Sprintf(common.WarningPreviewModeDetails, name),
		)

		return diags
	}

	return nil
}
