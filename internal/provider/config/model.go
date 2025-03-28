// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/fabric-sdk-go/fabric"

	"github.com/microsoft/terraform-provider-fabric/internal/auth"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type ProviderData struct {
	FabricClient *fabric.Client
	Timeout      time.Duration
	Endpoint     string
	Version      string
	Preview      bool
}

type ProviderConfig struct {
	*ProviderData
	Auth *auth.Config
}

// ProviderConfigModel describes the provider data model.
type ProviderConfigModel struct {
	Timeout                        timetypes.GoDuration `tfsdk:"timeout"`
	Endpoint                       customtypes.URL      `tfsdk:"endpoint"`
	Environment                    types.String         `tfsdk:"environment"`
	AuxiliaryTenantIDs             types.Set            `tfsdk:"auxiliary_tenant_ids"`
	TenantID                       customtypes.UUID     `tfsdk:"tenant_id"`
	ClientID                       customtypes.UUID     `tfsdk:"client_id"`
	ClientIDFilePath               types.String         `tfsdk:"client_id_file_path"`
	ClientSecret                   types.String         `tfsdk:"client_secret"`
	ClientSecretFilePath           types.String         `tfsdk:"client_secret_file_path"`
	ClientCertificate              types.String         `tfsdk:"client_certificate"`
	ClientCertificateFilePath      types.String         `tfsdk:"client_certificate_file_path"`
	ClientCertificatePassword      types.String         `tfsdk:"client_certificate_password"`
	OIDCRequestToken               types.String         `tfsdk:"oidc_request_token"`
	OIDCRequestURL                 types.String         `tfsdk:"oidc_request_url"`
	OIDCToken                      types.String         `tfsdk:"oidc_token"`
	OIDCTokenFilePath              types.String         `tfsdk:"oidc_token_file_path"`
	AzureDevOpsServiceConnectionID types.String         `tfsdk:"azure_devops_service_connection_id"`
	UseOIDC                        types.Bool           `tfsdk:"use_oidc"`
	UseCLI                         types.Bool           `tfsdk:"use_cli"`
	UseDevCLI                      types.Bool           `tfsdk:"use_dev_cli"`
	UseMSI                         types.Bool           `tfsdk:"use_msi"`
	Preview                        types.Bool           `tfsdk:"preview"`
}
