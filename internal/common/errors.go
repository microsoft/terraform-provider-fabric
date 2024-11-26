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
	ErrorCreateHeader                 string = "create operation"
	ErrorCreateDetails                string = "Could not create resource"
	ErrorReadHeader                   string = "read operation"
	ErrorReadDetails                  string = "Could not read resource"
	ErrorUpdateHeader                 string = "update operation"
	ErrorUpdateDetails                string = "Could not update resource"
	ErrorDeleteHeader                 string = "delete operation"
	ErrorDeleteDetails                string = "Could not delete resource"
	ErrorListHeader                   string = "list operation"
	ErrorListDetails                  string = "Could not list resource"
	ErrorImportHeader                 string = "import operation"
	ErrorImportDetails                string = "Could not import resource"
	ErrorImportIdentifierHeader       string = "Invalid import identifier"
	ErrorImportIdentifierDetails      string = "Expected identifier must be in the format: %s"
	ErrorInvalidURL                   string = "must be a valid URL."
	ErrorFabricClientType             string = "Expected *fabric.Client, got: %T. Please report this issue to the provider developers."
	ErrorGenericUnexpected            string = "Unexpected error occurred"
	ErrorBase64DecodeHeader           string = "base64 decode operation"
	ErrorBase64EncodeHeader           string = "base64 encode operation"
	ErrorBase64GzipEncodeHeader       string = "base64 gzip encode operation"
	ErrorJSONNormalizeHeader          string = "json normalize operation"
	ErrorFileReadHeader               string = "file read operation"
	ErrorTmplParseHeader              string = "template parse operation"
)

// Warnings.
const (
	WarningItemDefinitionUpdateHeader  = "Item definition update"
	WarningItemDefinitionUpdateDetails = "%s definition update will overwrite the existing definition."
)
