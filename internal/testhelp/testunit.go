// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package testhelp

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/provider"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var NewTestUnitCase = func(t *testing.T, testResource *string, fakeServer *fabfake.ServerFactory, preCheck func(*testing.T), steps []resource.TestStep) resource.TestCase {
	t.Helper()

	if preCheck == nil {
		return newTestUnitCase(t, testResource, fakeServer, nil, TestUnitPreCheck, steps)
	}

	return newTestUnitCase(t, testResource, fakeServer, nil, preCheck, steps)
}

var NewTestUnitCaseWithState = func(t *testing.T, testResource *string, fakeServer *fabfake.ServerFactory, testState *TestState, preCheck func(*testing.T), steps []resource.TestStep) resource.TestCase {
	t.Helper()

	if preCheck == nil {
		return newTestUnitCase(t, testResource, fakeServer, testState, TestUnitPreCheck, steps)
	}

	return newTestUnitCase(t, testResource, fakeServer, testState, preCheck, steps)
}

func TestUnitPreCheck(t *testing.T) {
	t.Helper()
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func TestUnitPreCheckNoEnvs(t *testing.T) {
	t.Helper()

	for _, env := range os.Environ() {
		envPair := strings.SplitN(env, "=", 2) // Split into key and value
		if (strings.HasPrefix(envPair[0], "FABRIC_") && !strings.HasPrefix(envPair[0], "FABRIC_TESTACC_")) || strings.HasPrefix(envPair[0], "ARM_") {
			os.Unsetenv(envPair[0]) //revive:disable-line:unhandled-error
		}
	}
}

func newTestUnitCase(t *testing.T, testResource *string, fakeServer *fabfake.ServerFactory, testState *TestState, preCheck func(*testing.T), steps []resource.TestStep) resource.TestCase {
	t.Helper()

	fabricClientOptions := &fabric.ClientOptions{}
	fabricClientOptions.Transport = fabfake.NewServerFactoryTransport(fakeServer)

	testCase := resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { preCheck(t) },
		CheckDestroy: func(s *terraform.State) error {
			if testResource != nil {
				_, ok := s.RootModule().Resources[*testResource]
				if !ok {
					return errors.New(*testResource + ` - resource still exists`)
				}
			}

			return nil
		},
		ProtoV6ProviderFactories: GetTestUnitProtoV6ProviderFactories(fabricClientOptions, testState),
		Steps:                    steps,
	}

	// ephemeral specific configurations
	if testResource != nil && strings.HasPrefix(*testResource, "ephemeral") {
		testCase.TerraformVersionChecks = []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		}
		testCase.CheckDestroy = func(_ *terraform.State) error {
			return nil // No need to check destroy for ephemeral resources
		}
	}

	// writeOnly specific configurations
	if strings.Contains(strings.ToLower(t.Name()), "writeonly") {
		testCase.TerraformVersionChecks = []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		}
	}

	return testCase
}

// GetTestUnitProtoV6ProviderFactories are used to instantiate a provider during
// unit testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
func GetTestUnitProtoV6ProviderFactories(fabricClientOpts *fabric.ClientOptions, testState *TestState) map[string]func() (tfprotov6.ProviderServer, error) {
	prov := provider.New("testUnit")
	prov.ConfigureCreateClient(func(ctx context.Context, cfg *pconfig.ProviderConfig) (*fabric.Client, error) {
		client, err := fabric.NewClient(&azfake.TokenCredential{}, &cfg.Endpoint, fabricClientOpts)
		if err != nil {
			tflog.Error(ctx, "Failed to initialize Microsoft Fabric test client", map[string]any{"error": err})

			return nil, err
		}

		return client, nil
	})
	prov.ConfigureAffirmProviderConfig(func(cfg *pconfig.ProviderConfig) {
		if testState != nil {
			testState.Config = cfg
		}
	})

	return map[string]func() (tfprotov6.ProviderServer, error){
		"fabric": providerserver.NewProtocol6WithError(prov),
		"echo":   echoprovider.NewProviderServer(),
	}
}

type TestState struct {
	Config *pconfig.ProviderConfig
}

func NewTestState() *TestState {
	return &TestState{}
}
