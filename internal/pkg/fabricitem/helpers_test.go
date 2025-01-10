// Copyright (c) Microsoft Corporation
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

func TestUnit_GetDataSourcePreviewNote_PreviewEnabled(t *testing.T) {
	md := "This is a test data-source."
	expected := md + fabricitem.PreviewDataSource
	result := fabricitem.GetDataSourcePreviewNote(md, true)

	assert.Equal(t, expected, result)
}

func TestUnit_GetDataSourcePreviewNote_PreviewDisabled(t *testing.T) {
	md := "This is a test data-source."
	expected := md
	result := fabricitem.GetDataSourcePreviewNote(md, false)

	assert.Equal(t, expected, result)
}

func TestUnit_GetResourcePreviewNote_PreviewEnabled(t *testing.T) {
	md := "This is a test resource."
	expected := md + fabricitem.PreviewResource
	result := fabricitem.GetResourcePreviewNote(md, true)

	assert.Equal(t, expected, result)
}

func TestUnit_GetResourcePreviewNote_PreviewDisabled(t *testing.T) {
	md := "This is a test resource."
	expected := md
	result := fabricitem.GetResourcePreviewNote(md, false)

	assert.Equal(t, expected, result)
}

func TestUnit_IsPreviewMode_ItemIsPreview_ProviderPreviewModeDisabled(t *testing.T) {
	name := testhelp.RandomName()
	itemIsPreview := true
	providerPreviewMode := false

	diags := fabricitem.IsPreviewMode(name, itemIsPreview, providerPreviewMode)

	assert.Len(t, diags, 1)
	assert.Equal(t, diags[0].Severity(), diag.SeverityError)
	assert.Equal(t, diags[0].Summary(), common.ErrorPreviewModeHeader)
	assert.Equal(t, diags[0].Detail(), fmt.Sprintf(common.ErrorPreviewModeDetails, name))
}

func TestUnit_IsPreviewMode_ItemIsPreview_ProviderPreviewModeEnabled(t *testing.T) {
	name := testhelp.RandomName()
	itemIsPreview := true
	providerPreviewMode := true

	diags := fabricitem.IsPreviewMode(name, itemIsPreview, providerPreviewMode)

	assert.Len(t, diags, 1)
	assert.Equal(t, diags[0].Severity(), diag.SeverityWarning)
	assert.Equal(t, diags[0].Summary(), common.WarningPreviewModeHeader)
	assert.Equal(t, diags[0].Detail(), fmt.Sprintf(common.WarningPreviewModeDetails, name))
}

func TestUnit_IsPreviewMode_ItemIsNotPreview(t *testing.T) {
	name := testhelp.RandomName()
	itemIsPreview := false
	providerPreviewMode := false

	diags := fabricitem.IsPreviewMode(name, itemIsPreview, providerPreviewMode)

	assert.Len(t, diags, 0)
}
