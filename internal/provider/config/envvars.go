// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package config

func GetEnvVarsTimeout() []string {
	return []string{"FABRIC_TIMEOUT"}
}

func GetEnvVarsEndpoint() []string {
	return []string{"FABRIC_ENDPOINT"}
}

func GetEnvVarsEnvironment() []string {
	return []string{"FABRIC_ENVIRONMENT", "ARM_ENVIRONMENT"}
}

func GetEnvVarsTenantID() []string {
	return []string{"FABRIC_TENANT_ID", "ARM_TENANT_ID"}
}

func GetEnvVarsAuxiliaryTenantIDs() []string {
	return []string{"FABRIC_AUXILIARY_TENANT_IDS", "ARM_AUXILIARY_TENANT_IDS"}
}

func GetEnvVarsClientID() []string {
	return []string{"FABRIC_CLIENT_ID", "ARM_CLIENT_ID"}
}

func GetEnvVarsClientIDFilePath() []string {
	return []string{"FABRIC_CLIENT_ID_FILE_PATH", "ARM_CLIENT_ID_FILE_PATH"}
}

func GetEnvVarsClientSecret() []string {
	return []string{"FABRIC_CLIENT_SECRET", "ARM_CLIENT_SECRET"}
}

func GetEnvVarsClientSecretFilePath() []string {
	return []string{"FABRIC_CLIENT_SECRET_FILE_PATH", "ARM_CLIENT_SECRET_FILE_PATH"}
}

func GetEnvVarsClientCertificate() []string {
	return []string{"FABRIC_CLIENT_CERTIFICATE", "ARM_CLIENT_CERTIFICATE"}
}

func GetEnvVarsClientCertificateFilePath() []string {
	return []string{"FABRIC_CLIENT_CERTIFICATE_FILE_PATH", "ARM_CLIENT_CERTIFICATE_FILE_PATH", "ARM_CLIENT_CERTIFICATE_PATH"}
}

func GetEnvVarsClientCertificatePassword() []string {
	return []string{"FABRIC_CLIENT_CERTIFICATE_PASSWORD", "ARM_CLIENT_CERTIFICATE_PASSWORD"}
}

func GetEnvVarsOIDCRequestURL() []string {
	return []string{"FABRIC_OIDC_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_URL", "ARM_OIDC_REQUEST_URL"}
}

func GetEnvVarsOIDCRequestToken() []string {
	return []string{"FABRIC_OIDC_REQUEST_TOKEN", "ACTIONS_ID_TOKEN_REQUEST_TOKEN", "SYSTEM_ACCESSTOKEN", "ARM_OIDC_REQUEST_TOKEN"}
}

func GetEnvVarsOIDCToken() []string {
	return []string{"FABRIC_OIDC_TOKEN", "ARM_OIDC_TOKEN"}
}

func GetEnvVarsOIDCTokenFilePath() []string {
	return []string{"FABRIC_OIDC_TOKEN_FILE_PATH", "ARM_OIDC_TOKEN_FILE_PATH"}
}

func GetEnvVarsAzureDevOpsServiceConnectionID() []string {
	return []string{"FABRIC_AZURE_DEVOPS_SERVICE_CONNECTION_ID"}
}

func GetEnvVarsToken() []string {
	return []string{"FABRIC_TOKEN"}
}

func GetEnvVarsTokenFilePath() []string {
	return []string{"FABRIC_TOKEN_FILE_PATH"}
}

func GetEnvVarsUseOIDC() []string {
	return []string{"FABRIC_USE_OIDC", "ARM_USE_OIDC"}
}

func GetEnvVarsUseMSI() []string {
	return []string{"FABRIC_USE_MSI", "ARM_USE_MSI"}
}

func GetEnvVarsUseDevCLI() []string {
	return []string{"FABRIC_USE_DEV_CLI", "ARM_USE_DEV_CLI"}
}

func GetEnvVarsUseCLI() []string {
	return []string{"FABRIC_USE_CLI", "ARM_USE_CLI"}
}
