// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/assert"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func TestUnit_IsPreviewMode_ItemIsPreview_ProviderPreviewModeDisabled(t *testing.T) {
	name := testhelp.RandomName()
	itemIsPreview := true
	providerPreviewMode := false

	diags := fabricitem.IsPreviewMode(name, itemIsPreview, providerPreviewMode)

	assert.Len(t, diags, 1)
	assert.Equal(t, diag.SeverityError, diags[0].Severity())
	assert.Equal(t, common.ErrorPreviewModeHeader, diags[0].Summary())
	assert.Equal(t, fmt.Sprintf(common.ErrorPreviewModeDetails, name), diags[0].Detail())
}

func TestUnit_IsPreviewMode_ItemIsPreview_ProviderPreviewModeEnabled(t *testing.T) {
	name := testhelp.RandomName()
	itemIsPreview := true
	providerPreviewMode := true

	diags := fabricitem.IsPreviewMode(name, itemIsPreview, providerPreviewMode)

	assert.Len(t, diags, 1)
	assert.Equal(t, diag.SeverityWarning, diags[0].Severity())
	assert.Equal(t, fmt.Sprintf(common.WarningPreviewModeHeader, name), diags[0].Summary())
	assert.Equal(t, fmt.Sprintf(common.WarningPreviewModeDetails, name), diags[0].Detail())
}

func TestUnit_IsPreviewMode_ItemIsNotPreview(t *testing.T) {
	name := testhelp.RandomName()
	itemIsPreview := false
	providerPreviewMode := false

	diags := fabricitem.IsPreviewMode(name, itemIsPreview, providerPreviewMode)

	assert.Empty(t, diags)
}
