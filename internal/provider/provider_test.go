// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/auth"
	"github.com/microsoft/terraform-provider-fabric/internal/common"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testHelperEntity = fakes.NewRandomWorkspaceInfo()
	testHelperHCL    = fmt.Sprintf(`
		data "fabric_workspace" "test" {
			id = "%s"
		}`, *testHelperEntity.ID)
)

func TestUnit_Provider_Configurations(t *testing.T) {
	testState := testhelp.NewTestState()

	fakes.FakeServer.Upsert(testHelperEntity)

	tenantID := testhelp.RandomUUID()

	resource.ParallelTest(t, testhelp.NewTestUnitCaseWithState(t, nil, fakes.FakeServer.ServerFactory, testState, testhelp.TestUnitPreCheckNoEnvs, []resource.TestStep{
		{
			// lintignore:AT004
			Config: fmt.Sprintf(`
				provider "fabric" {
					tenant_id = "%s"
				}`+testHelperHCL, tenantID),
			Check: resource.ComposeAggregateTestCheckFunc(
				testhelp.CheckAuthConfig("TenantID", tenantID, testState),
			),
		},
		// explicitly setting service principal with secret
		{
			// lintignore:AT004
			Config: fmt.Sprintf(`
				provider "fabric" {
					client_id = "00000000-0000-0000-0000-000000000000"
					tenant_id = "00000000-0000-0000-0000-000000000000"
					client_secret = "%s"
				}`+testHelperHCL, testhelp.RandomUUID()),
			Check: resource.ComposeAggregateTestCheckFunc(
				testhelp.CheckAuthMethod(auth.ServicePrincipalSecretAuth, testState),
			),
		},
	}))
}

func TestUnit_Provider_AuthAttributes(t *testing.T) {
	testState := testhelp.NewTestState()

	fakes.FakeServer.Upsert(testHelperEntity)

	resource.ParallelTest(t, testhelp.NewTestUnitCaseWithState(t, nil, fakes.FakeServer.ServerFactory, testState, testhelp.TestUnitPreCheckNoEnvs, []resource.TestStep{
		// Not allowed multiple auth methods
		{
			// lintignore:AT004
			Config: `
				provider "fabric" {
					use_oidc = true
					use_msi = true
				}` + testHelperHCL,
			ExpectError: regexp.MustCompile(common.ErrorAttComboInvalid),
		},
		// Inconsistent client_id and client_id_file_path
		{
			// lintignore:AT004
			Config: `
				provider "fabric" {
					client_id = "00000000-0000-0000-0000-000000000000"
					client_id_file_path = "fixtures/client_id_1.txt"
				}` + testHelperHCL,
			ExpectError: regexp.MustCompile(common.ErrorInvalidValue),
		},
	}))
}

func TestUnit_Provider_AuthOIDC(t *testing.T) {
	testState := testhelp.NewTestState()

	fakes.FakeServer.Upsert(testHelperEntity)

	resource.Test(t, testhelp.NewTestUnitCaseWithState(t, nil, fakes.FakeServer.ServerFactory, testState, testhelp.TestUnitPreCheckNoEnvs, []resource.TestStep{
		// Missing tenant_id
		{
			// lintignore:AT004
			Config: `
				provider "fabric" {
					use_oidc = true
				}` + testHelperHCL,
			ExpectError: regexp.MustCompile(common.ErrorInvalidConfig),
		},
		// Missing client_id
		{
			// lintignore:AT004
			Config: `
				provider "fabric" {
					use_oidc = true
					tenant_id = "00000000-0000-0000-0000-000000000000"
				}` + testHelperHCL,
			ExpectError: regexp.MustCompile(common.ErrorInvalidConfig),
		},
		// Missing OIDC details
		{
			PreConfig: func() {
				envKeys := pconfig.GetEnvVarsOIDCRequestURL()
				envKeys = append(envKeys, pconfig.GetEnvVarsOIDCRequestToken()...)
				envKeys = append(envKeys, pconfig.GetEnvVarsOIDCToken()...)
				envKeys = append(envKeys, pconfig.GetEnvVarsOIDCTokenFilePath()...)

				for _, envKey := range envKeys {
					os.Unsetenv(envKey) //revive:disable-line:unhandled-error
				}
			},
			// lintignore:AT004
			Config: `
				provider "fabric" {
					use_oidc = true
					tenant_id = "00000000-0000-0000-0000-000000000000"
					client_id = "00000000-0000-0000-0000-000000000000"
				}` + testHelperHCL,
			ExpectError: regexp.MustCompile(common.ErrorInvalidConfig),
		},
		// OIDC auth OK
		{
			PreConfig: func() {
				t.Setenv(pconfig.GetEnvVarsOIDCRequestURL()[0], "https://localhost")
				t.Setenv(pconfig.GetEnvVarsOIDCRequestToken()[0], testhelp.RandomUUID())
			},
			// lintignore:AT004
			Config: `
				provider "fabric" {
					use_oidc = true
					tenant_id = "00000000-0000-0000-0000-000000000000"
					client_id = "00000000-0000-0000-0000-000000000000"
				}` + testHelperHCL,
			Check: resource.ComposeAggregateTestCheckFunc(
				testhelp.CheckAuthMethod(auth.ServicePrincipalOIDCAuth, testState),
			),
		},
	}))
}

func TestUnit_Provider_AuthAzDevOpsWI(t *testing.T) {
	testState := testhelp.NewTestState()

	fakes.FakeServer.Upsert(testHelperEntity)

	resource.Test(t, testhelp.NewTestUnitCaseWithState(t, nil, fakes.FakeServer.ServerFactory, testState, testhelp.TestUnitPreCheckNoEnvs, []resource.TestStep{
		// Missing OIDC details
		{
			PreConfig: func() {
				envKeys := pconfig.GetEnvVarsOIDCRequestURL()
				envKeys = append(envKeys, pconfig.GetEnvVarsOIDCRequestToken()...)
				envKeys = append(envKeys, pconfig.GetEnvVarsOIDCToken()...)
				envKeys = append(envKeys, pconfig.GetEnvVarsOIDCTokenFilePath()...)

				for _, envKey := range envKeys {
					os.Unsetenv(envKey) //revive:disable-line:unhandled-error
				}
			},
			// lintignore:AT004
			Config: `
				provider "fabric" {
					use_oidc = true
					tenant_id = "00000000-0000-0000-0000-000000000000"
					client_id = "00000000-0000-0000-0000-000000000000"
					azure_devops_service_connection_id = "00000000-0000-0000-0000-000000000000"
				}` + testHelperHCL,
			ExpectError: regexp.MustCompile(common.ErrorInvalidConfig),
		},
		// OIDC auth OK
		{
			PreConfig: func() {
				t.Setenv("SYSTEM_OIDCREQUESTURI", "https://localhost")
				t.Setenv(pconfig.GetEnvVarsOIDCRequestToken()[0], testhelp.RandomUUID())
			},
			// lintignore:AT004
			Config: `
				provider "fabric" {
					use_oidc = true
					tenant_id = "00000000-0000-0000-0000-000000000000"
					client_id = "00000000-0000-0000-0000-000000000000"
					azure_devops_service_connection_id = "00000000-0000-0000-0000-000000000000"
				}` + testHelperHCL,
			Check: resource.ComposeAggregateTestCheckFunc(
				testhelp.CheckAuthMethod(auth.AzureDevOpsWorkloadIdentityFederationAuth, testState),
			),
		},
	}))
}

func TestUnit_Provider_AuthMSI(t *testing.T) {
	testState := testhelp.NewTestState()

	fakes.FakeServer.Upsert(testHelperEntity)

	resource.ParallelTest(t, testhelp.NewTestUnitCaseWithState(t, nil, fakes.FakeServer.ServerFactory, testState, testhelp.TestUnitPreCheckNoEnvs, []resource.TestStep{
		// Missing tenant_id
		{
			// lintignore:AT004
			Config: `
				provider "fabric" {
					use_msi = true
				}` + testHelperHCL,
			ExpectError: regexp.MustCompile(common.ErrorInvalidConfig),
		},
		// MSI User auth OK
		{
			// lintignore:AT004
			Config: `
				provider "fabric" {
					use_msi = true
					tenant_id = "00000000-0000-0000-0000-000000000000"
					client_id = "00000000-0000-0000-0000-000000000000"
				}` + testHelperHCL,
			Check: resource.ComposeAggregateTestCheckFunc(
				testhelp.CheckAuthMethod(auth.ManagedServiceIdentityUserAuth, testState),
			),
		},
		// MSI System auth OK
		{
			// lintignore:AT004
			Config: `
				provider "fabric" {
					use_msi = true
					tenant_id = "00000000-0000-0000-0000-000000000000"
				}` + testHelperHCL,
			Check: resource.ComposeAggregateTestCheckFunc(
				testhelp.CheckAuthMethod(auth.ManagedServiceIdentitySystemAuth, testState),
			),
		},
	}))
}

func TestUnit_Provider_AuthCLI(t *testing.T) {
	testState := testhelp.NewTestState()

	fakes.FakeServer.Upsert(testHelperEntity)

	resource.ParallelTest(t, testhelp.NewTestUnitCaseWithState(t, nil, fakes.FakeServer.ServerFactory, testState, testhelp.TestUnitPreCheckNoEnvs, []resource.TestStep{
		// if auth is not explicitly set, use cli should be true
		{
			// lintignore:AT004
			Config: `
				provider "fabric" {}
			` + testHelperHCL,
			Check: resource.ComposeAggregateTestCheckFunc(
				testhelp.CheckAuthMethod(auth.AzureCLIAuth, testState),
			),
		},
		// explicitly setting use cli to true
		{
			// lintignore:AT004
			Config: `
				provider "fabric" {
					use_cli = true
				}` + testHelperHCL,
			Check: resource.ComposeAggregateTestCheckFunc(
				testhelp.CheckAuthMethod(auth.AzureCLIAuth, testState),
			),
		},
		// explicitly setting use cli to true and use_msi to false
		{
			// lintignore:AT004
			Config: `
				provider "fabric" {
					use_msi = false
					use_cli = true
				}` + testHelperHCL,
			Check: resource.ComposeAggregateTestCheckFunc(
				testhelp.CheckAuthMethod(auth.AzureCLIAuth, testState),
			),
		},
	}))
}

func TestUnit_Provider_AuthDevCLI(t *testing.T) {
	testState := testhelp.NewTestState()

	fakes.FakeServer.Upsert(testHelperEntity)

	resource.ParallelTest(t, testhelp.NewTestUnitCaseWithState(t, nil, fakes.FakeServer.ServerFactory, testState, testhelp.TestUnitPreCheckNoEnvs, []resource.TestStep{
		// explicitly setting use dev cli to true
		{
			// lintignore:AT004
			Config: `
				provider "fabric" {
					use_dev_cli = true
				}` + testHelperHCL,
			Check: resource.ComposeAggregateTestCheckFunc(
				testhelp.CheckAuthMethod(auth.AzureDevCLIAuth, testState),
			),
		},
		// explicitly setting use dev cli to true and use_msi to false
		{
			// lintignore:AT004
			Config: `
				provider "fabric" {
					use_msi = false
					use_dev_cli = true
				}` + testHelperHCL,
			Check: resource.ComposeAggregateTestCheckFunc(
				testhelp.CheckAuthMethod(auth.AzureDevCLIAuth, testState),
			),
		},
	}))
}

func TestUnit_Provider_AuthSecret(t *testing.T) {
	testState := testhelp.NewTestState()

	fakes.FakeServer.Upsert(testHelperEntity)

	resource.ParallelTest(t, testhelp.NewTestUnitCaseWithState(t, nil, fakes.FakeServer.ServerFactory, testState, testhelp.TestUnitPreCheckNoEnvs, []resource.TestStep{
		// Missing client_id
		{
			// lintignore:AT004
			Config: fmt.Sprintf(`
				provider "fabric" {
					tenant_id = "00000000-0000-0000-0000-000000000000"
					client_secret = "%s"
				}`+testHelperHCL, testhelp.RandomUUID()),
			ExpectError: regexp.MustCompile(common.ErrorInvalidConfig),
		},
		// Missing tenant_id
		{
			// lintignore:AT004
			Config: fmt.Sprintf(`
				provider "fabric" {
					client_id = "00000000-0000-0000-0000-000000000000"
					client_secret = "%s"
				}`+testHelperHCL, testhelp.RandomUUID()),
			ExpectError: regexp.MustCompile(common.ErrorInvalidConfig),
		},
		// Secret auth OK
		{
			// lintignore:AT004
			Config: fmt.Sprintf(`
				provider "fabric" {
					tenant_id = "00000000-0000-0000-0000-000000000000"
					client_id = "00000000-0000-0000-0000-000000000000"
					client_secret = "%s"
				}`+testHelperHCL, testhelp.RandomUUID()),
			Check: resource.ComposeAggregateTestCheckFunc(
				testhelp.CheckAuthMethod(auth.ServicePrincipalSecretAuth, testState),
			),
		},
	}))
}

func TestUnit_Provider_AuthCertificate(t *testing.T) {
	testState := testhelp.NewTestState()

	fakes.FakeServer.Upsert(testHelperEntity)

	certPass := testhelp.RandomUUID()
	certWithPass := testhelp.RandomP12CertB64(certPass)
	certNoPass := testhelp.RandomP12CertB64("")

	resource.ParallelTest(t, testhelp.NewTestUnitCaseWithState(t, nil, fakes.FakeServer.ServerFactory, testState, testhelp.TestUnitPreCheckNoEnvs, []resource.TestStep{
		// Missing client_id
		{
			// lintignore:AT004
			Config: fmt.Sprintf(`
				provider "fabric" {
					tenant_id = "00000000-0000-0000-0000-000000000000"
					client_certificate = "%s"
				}`+testHelperHCL, certNoPass),
			ExpectError: regexp.MustCompile(common.ErrorInvalidConfig),
		},
		// Missing tenant_id
		{
			// lintignore:AT004
			Config: fmt.Sprintf(`
				provider "fabric" {
					client_id = "00000000-0000-0000-0000-000000000000"
					client_certificate = "%s"
				}`+testHelperHCL, certNoPass),
			ExpectError: regexp.MustCompile(common.ErrorInvalidConfig),
		},
		// Certificate no pass - auth OK
		{
			// lintignore:AT004
			Config: fmt.Sprintf(`
				provider "fabric" {
					tenant_id = "00000000-0000-0000-0000-000000000000"
					client_id = "00000000-0000-0000-0000-000000000000"
					client_certificate = "%s"
				}`+testHelperHCL, certNoPass),
			Check: resource.ComposeAggregateTestCheckFunc(
				testhelp.CheckAuthMethod(auth.ServicePrincipalCertificateAuth, testState),
			),
		},
		// Certificate with pass - auth OK
		{
			// lintignore:AT004
			Config: fmt.Sprintf(`
				provider "fabric" {
					tenant_id = "00000000-0000-0000-0000-000000000000"
					client_id = "00000000-0000-0000-0000-000000000000"
					client_certificate = "%s"
					client_certificate_password = "%s"
				}`+testHelperHCL, certWithPass, certPass),
			Check: resource.ComposeAggregateTestCheckFunc(
				testhelp.CheckAuthMethod(auth.ServicePrincipalCertificateAuth, testState),
			),
		},
	}))
}
