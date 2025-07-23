// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

type dataSourceConnectionModel[ConnectionDetails dsConnectionDetailsModel | rsConnectionDetailsModel, CredentialDetails dsCredentialDetailsModel | rsCredentialDetailsModel] struct {
	baseConnectionModel[ConnectionDetails, CredentialDetails]

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type dsConnectionDetailsModel struct {
	Path types.String `tfsdk:"path"`
	Type types.String `tfsdk:"type"`
}

func setDSConnectionDetails(from fabcore.ListConnectionDetails, to *dsConnectionDetailsModel) {
	to.Path = types.StringPointerValue(from.Path)
	to.Type = types.StringPointerValue(from.Type)
}

type dsCredentialDetailsModel struct {
	ConnectionEncryption types.String `tfsdk:"connection_encryption"`
	CredentialType       types.String `tfsdk:"credential_type"`
	SingleSignOnType     types.String `tfsdk:"single_sign_on_type"`
	SkipTestConnection   types.Bool   `tfsdk:"skip_test_connection"`
}

func setDSCredentialDetails(from fabcore.ListCredentialDetails, to *dsCredentialDetailsModel) {
	to.ConnectionEncryption = types.StringPointerValue((*string)(from.ConnectionEncryption))
	to.CredentialType = types.StringPointerValue((*string)(from.CredentialType))
	to.SingleSignOnType = types.StringPointerValue((*string)(from.SingleSignOnType))
	to.SkipTestConnection = types.BoolPointerValue(from.SkipTestConnection)
}
