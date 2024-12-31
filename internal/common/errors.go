// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package common

// Errors.
const (
	ErrorWorkspaceNotSupportedHeader  string = "Workspace type not supported"
	ErrorWorkspaceNotSupportedDetails string = "Cannot interact with '%s' workspace type"
	ErrorInvalidConfig                string = "Invalid Configuration"
	ErrorInvalidValue                 string = "Invalid Value"
	ErrorAttComboInvalid              string = "Invalid Attribute Combination"
	ErrorAttConfigMissing             string = "Missing Attribute Configuration"
	ErrorAttValueMatch                string = "Invalid Attribute Value Match"
	ErrorDataSourceConfigType         string = "Unexpected Data Source Configure Type"
	ErrorResourceConfigType           string = "Unexpected Resource Configure Type"
	ErrorModelConversion              string = "Data Model Conversion Error"
	ErrorCreateHeader                 string = "Create operation"
	ErrorCreateDetails                string = "Could not create resource"
	ErrorReadHeader                   string = "Read operation"
	ErrorReadDetails                  string = "Could not read resource"
	ErrorUpdateHeader                 string = "Update operation"
	ErrorUpdateDetails                string = "Could not update resource"
	ErrorDeleteHeader                 string = "Delete operation"
	ErrorDeleteDetails                string = "Could not delete resource"
	ErrorListHeader                   string = "List operation"
	ErrorListDetails                  string = "Could not list resource"
	ErrorImportHeader                 string = "Import operation"
	ErrorImportDetails                string = "Could not import resource"
	ErrorImportIdentifierHeader       string = "Invalid import identifier"
	ErrorImportIdentifierDetails      string = "Expected identifier must be in the format: %s"
	ErrorInvalidURL                   string = "must be a valid URL."
	ErrorFabricClientType             string = "Expected *fabric.Client, got: %T. Please report this issue to the provider developers."
	ErrorGenericUnexpected            string = "Unexpected error occurred"
	ErrorBase64DecodeHeader           string = "Base64 decode operation"
	ErrorBase64EncodeHeader           string = "Base64 encode operation"
	ErrorBase64GzipEncodeHeader       string = "Base64 Gzip encode operation"
	ErrorJSONNormalizeHeader          string = "JSON normalize operation"
	ErrorFileReadHeader               string = "File read operation"
	ErrorTmplParseHeader              string = "template parse operation"
)

// Warnings.
const (
	WarningItemDefinitionUpdateHeader  = "Fabric Item definition update"
	WarningItemDefinitionUpdateDetails = "%s definition update operation will overwrite the existing definition on the Fabric side."
)
