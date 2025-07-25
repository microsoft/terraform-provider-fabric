// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package common

const (
	ErrorWorkspaceNotSupportedHeader  = "Workspace type not supported"
	ErrorWorkspaceNotSupportedDetails = "Cannot interact with '%s' workspace type"
	ErrorInvalidConfig                = "Invalid Configuration"
	ErrorInvalidValue                 = "Invalid Value"
	ErrorAttComboInvalid              = "Invalid Attribute Combination"
	ErrorAttConfigMissing             = "Missing Attribute Configuration"
	ErrorAttValueMatch                = "Invalid Attribute Value Match"
	ErrorDataSourceConfigType         = "Unexpected Data Source Configure Type"
	ErrorResourceConfigType           = "Unexpected Resource Configure Type"
	ErrorEphemeralResourceConfigType  = "Unexpected Ephemeral Resource Configure Type"
	ErrorModelConversion              = "Data Model Conversion Error"
	ErrorCreateHeader                 = "Create operation"
	ErrorCreateDetails                = "Could not create resource"
	ErrorReadHeader                   = "Read operation"
	ErrorReadDetails                  = "Could not read resource"
	ErrorUpdateHeader                 = "Update operation"
	ErrorUpdateDetails                = "Could not update resource"
	ErrorDeleteHeader                 = "Delete operation"
	ErrorDeleteDetails                = "Could not delete resource"
	ErrorListHeader                   = "List operation"
	ErrorListDetails                  = "Could not list resource"
	ErrorImportHeader                 = "Import operation"
	ErrorImportDetails                = "Could not import resource"
	ErrorImportIdentifierHeader       = "Invalid import identifier"
	ErrorImportIdentifierDetails      = "Expected identifier must be in the format: %s"
	ErrorOpenHeader                   = "Open operation"
	ErrorOpenDetails                  = "Could not open resource"
	ErrorInvalidURL                   = "must be a valid URL."
	ErrorFabricClientType             = "Expected *fabric.Client, got: %T. Please report this issue to the provider developers."
	ErrorGenericUnexpected            = "Unexpected error occurred"
	ErrorBase64DecodeHeader           = "Base64 decode operation"
	ErrorBase64EncodeHeader           = "Base64 encode operation"
	ErrorBase64GzipEncodeHeader       = "Base64 Gzip encode operation"
	ErrorJSONNormalizeHeader          = "JSON normalize operation"
	ErrorFileReadHeader               = "File read operation"
	ErrorTmplParseHeader              = "Template parse operation"
	ErrorPreviewModeHeader            = "Preview mode not enabled"
	ErrorPreviewModeDetails           = "'%s' is not available without explicitly opt-in to the preview mode on the provider level configuration."
)
