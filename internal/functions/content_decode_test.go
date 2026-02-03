// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package functions_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testFunctionContentDecodeHeader = testhelp.FunctionHeader("fabric", "content_decode")

func TestUnit_ContentDecodeFunction(t *testing.T) {
	// Valid Base64/gzip content, valid json
	const testFunctionContentDecodeFixture1 = "H4sIAAAAAAAACqtWKi4pUrJSKkvMKU01VNJRyswrUbIy1FFKys/PUbIqKSpN1VHKT8pSsqpGVmkEU2kEU5mWmFOcWlsLACn4/TdQAAAA" // `{"str":"value1","int":1,"bool":true,"obj":{"str":"value2","int":2,"bool":false}}`

	// Valid Base64/gzip content, invalid JSON
	const testFunctionContentDecodeFixture2 = "H4sIAAAAAAAACvPJL0rNVcgsKC7NVUjJz8kvAgAy+4dOEQAAAA==" // `Lorem ipsum dolor`

	// Valid Base64/gzip content, valid json array
	const testFunctionContentDecodeFixture3 = "H4sIAAAAAAAACos21DHSUSrJKEpNVYoFAIFAs94NAAAA" // `[1,2,"three"]`

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// -- Happy path tests --
		{
			Config: fmt.Sprintf(`
				output "test" {
					value = %s("%s")
				}
			`, testFunctionContentDecodeHeader, testFunctionContentDecodeFixture1),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue("test", knownvalue.ObjectExact(map[string]knownvalue.Check{
					"str":  knownvalue.StringExact("value1"),
					"int":  knownvalue.Int64Exact(1),
					"bool": knownvalue.Bool(true),
					"obj": knownvalue.ObjectExact(map[string]knownvalue.Check{
						"str":  knownvalue.StringExact("value2"),
						"int":  knownvalue.Int64Exact(2),
						"bool": knownvalue.Bool(false),
					}),
				})),
			},
		},
		// Test with empty JSONPath (should default to whole document)
		{
			Config: fmt.Sprintf(`
				output "test" {
					value = %s("%s", "")
				}
			`, testFunctionContentDecodeHeader, testFunctionContentDecodeFixture1),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue("test", knownvalue.ObjectExact(map[string]knownvalue.Check{
					"str":  knownvalue.StringExact("value1"),
					"int":  knownvalue.Int64Exact(1),
					"bool": knownvalue.Bool(true),
					"obj": knownvalue.ObjectExact(map[string]knownvalue.Check{
						"str":  knownvalue.StringExact("value2"),
						"int":  knownvalue.Int64Exact(2),
						"bool": knownvalue.Bool(false),
					}),
				})),
			},
		},
		{
			Config: fmt.Sprintf(`
				output "test" {
					value = %s("%s", ".str")
				}
			`, testFunctionContentDecodeHeader, testFunctionContentDecodeFixture1),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact("value1")),
			},
		},
		{
			Config: fmt.Sprintf(`
				output "test" {
					value = %s("%s").str
				}
			`, testFunctionContentDecodeHeader, testFunctionContentDecodeFixture1),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact("value1")),
			},
		},
		{
			Config: fmt.Sprintf(`
				output "test" {
					value = %s("%s", ".obj").str
				}
			`, testFunctionContentDecodeHeader, testFunctionContentDecodeFixture1),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact("value2")),
			},
		},
		{
			Config: fmt.Sprintf(`
				output "test" {
					value = %s("%s")
				}
			`, testFunctionContentDecodeHeader, testFunctionContentDecodeFixture2),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact("Lorem ipsum dolor")),
			},
		},
		{
			Config: fmt.Sprintf(`
				output "test" {
					value = %s("%s", ".test")
				}
			`, testFunctionContentDecodeHeader, testFunctionContentDecodeFixture2),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact("Lorem ipsum dolor")),
			},
		},
		// Test with JSON array
		{
			Config: fmt.Sprintf(`
				output "test" {
					value = %s("%s")
				}
					`, testFunctionContentDecodeHeader, testFunctionContentDecodeFixture3),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue("test", knownvalue.ListExact([]knownvalue.Check{
					knownvalue.Int64Exact(1),
					knownvalue.Int64Exact(2),
					knownvalue.StringExact("three"),
				})),
			},
		},
		// Test with JSON array and JSONPath to access specific element
		{
			Config: fmt.Sprintf(`
				output "test" {
					value = %s("%s", "$[2]")
				}
			`, testFunctionContentDecodeHeader, testFunctionContentDecodeFixture3),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact("three")),
			},
		},
		// -- Error path tests --
		{
			Config: fmt.Sprintf(`
				output "test" {
					value = %s("")
				}
			`, testFunctionContentDecodeHeader),
			// Terraform: Error in function call
			// OpenTofu: Invalid function argument
			ExpectError: regexp.MustCompile(`Error in function call|Invalid function argument`),
		},
		{
			Config: fmt.Sprintf(`
				output "test" {
					value = %s()
				}
			`, testFunctionContentDecodeHeader),
			ExpectError: regexp.MustCompile("Not enough function arguments"),
		},
		// Invalid JSONPath expression test
		{
			Config: fmt.Sprintf(`
				output "test" {
					value = %s("%s", "[[[")
				}
			`, testFunctionContentDecodeHeader, testFunctionContentDecodeFixture1),
			ExpectError: regexp.MustCompile(`Error in function call|Invalid function argument`), // "Failed to parse JSONPath expression"
		},
		// Test with non-existent JSONPath but valid expression
		{
			Config: fmt.Sprintf(`
				output "test" {
					value = %s("%s", ".non_existent_field")
				}
			`, testFunctionContentDecodeHeader, testFunctionContentDecodeFixture1),
			ExpectError: regexp.MustCompile(`Error in function call|Invalid function argument`), // "JSONPath expression did not match any elements"
		},
	}))
}
