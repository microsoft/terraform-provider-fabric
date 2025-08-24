// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package testhelp

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/microsoft/terraform-provider-fabric/internal/provider"
)

var NewTestAccCase = func(t *testing.T, testResource *string, preCheck func(*testing.T), steps []resource.TestStep) resource.TestCase {
	t.Helper()

	if preCheck == nil {
		return newTestAccCase(t, testResource, TestAccPreCheck, steps)
	}

	return newTestAccCase(t, testResource, preCheck, steps)
}

// lintignore:AT003
func TestAccPreCheck(t *testing.T) {
	t.Helper()
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

// lintignore:AT003
func TestAccPreCheckNoEnvs(t *testing.T) {
	t.Helper()

	for _, env := range os.Environ() {
		envPair := strings.SplitN(env, "=", 2) // Split into key and value
		if strings.HasPrefix(envPair[0], "FABRIC_") && !strings.HasPrefix(envPair[0], "FABRIC_TESTACC_") {
			os.Unsetenv(envPair[0]) //revive:disable-line:unhandled-error
		}
	}
}

func newTestAccCase(t *testing.T, testResource *string, preCheck func(*testing.T), steps []resource.TestStep) resource.TestCase {
	t.Helper()

	testCase := resource.TestCase{
		IsUnitTest: false,
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
		ProtoV6ProviderFactories: getTestAccProtoV6ProviderFactories(),
		ExternalProviders: map[string]resource.ExternalProvider{
			"azurerm": {
				Source: "hashicorp/azurerm",
			},
		},
		Steps: steps,
	}

	// ephemeral specific configurations
	if testResource != nil && strings.HasPrefix(*testResource, "ephemeral") {
		testCase.TerraformVersionChecks = []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		}
		testCase.CheckDestroy = func(_ *terraform.State) error {
			return nil // No need to check destroy for ephemeral resources
		}
	}

	// writeOnly specific configurations
	if testResource != nil && strings.Contains(*testResource, "WriteOnly") {
		testCase.TerraformVersionChecks = []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		}
	}

	return testCase
}

// getTestAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
func getTestAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"fabric": providerserver.NewProtocol6WithError(provider.New("testAcc")),
		"echo":   echoprovider.NewProviderServer(),
	}
}

// lintignore:AT003
func TestAccWorkspaceResource(t *testing.T, capacityID string) (resourceHCL, resourceFQN string) { //nolint:nonamedreturns
	t.Helper()

	resourceHCL = at.CompileConfig(
		at.ResourceHeader(TypeName("fabric", "workspace"), "test"),
		map[string]any{
			"display_name": RandomName(),
			"description":  "testacc",
			"capacity_id":  capacityID,
		},
	)

	resourceFQN = ResourceFQN("fabric", "workspace", "test")

	return resourceHCL, resourceFQN
}

// ShouldSkipTest checks if a test should be skipped based on FABRIC_TESTACC_SKIP_NO_SPN environment variable.
func ShouldSkipTest(t *testing.T) bool {
	t.Helper()

	if skip, ok := os.LookupEnv("FABRIC_TESTACC_SKIP_NO_SPN"); ok && skip != "" {
		skipBool, err := strconv.ParseBool(skip)
		if err != nil {
			return false
		}

		return skipBool
	}

	return false
}

func GetFixturesDirPath(fixtureDir ...string) string {
	_, filename, _, _ := runtime.Caller(0) //nolint:dogsled
	testHelpDir := filepath.Dir(filename)

	var tempPath []string

	tempPath = append(tempPath, testHelpDir)
	tempPath = append(tempPath, "fixtures")
	tempPath = append(tempPath, fixtureDir...)

	return filepath.ToSlash(filepath.Join(tempPath...))
}
