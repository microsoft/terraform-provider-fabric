// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

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
