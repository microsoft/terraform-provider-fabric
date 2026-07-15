// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"slices"
	"testing"
)

func Test_parseSPNSupport(t *testing.T) {
	t.Parallel()

	const heading = "// MICROSOFT ENTRA SUPPORTED IDENTITIES This API supports the Microsoft identities listed in this section.\n"

	tests := []struct {
		name string
		doc  string
		want spnSupport
	}{
		{
			"service principal yes",
			heading + "// | Identity | Support | |-|-| | User | Yes | | Service principal [/entra/foo] and Managed identities [/entra/bar] | Yes |",
			spnYes,
		},
		{
			"service principal no",
			heading + "// | Identity | Support | |-|-| | User | Yes | | Service principal [/entra/foo] | No |",
			spnNo,
		},
		{
			"conditional support",
			heading + "// | Identity | Support | |-|-| | User | Yes | | Service principal [/entra/foo] and Managed identities [/entra/bar] | When the item type in the call is supported. |",
			spnConditional,
		},
		{
			"table wrapped across comment lines",
			heading +
				"// | Identity | Support | |-|-| | User | Yes | | Service principal [/entra/foo]\n" +
				"// and Managed identities\n" +
				"// [/entra/bar] | No |",
			spnNo,
		},
		{
			"no identity table",
			"// GetFoo retrieves a foo. It does not document identities.",
			spnUnknown,
		},
		{
			"heading present but no service principal row",
			heading + "// | Identity | Support | |-|-| | User | Yes |",
			spnUnknown,
		},
		{
			"service principal in free text before table",
			"// PERMISSIONS The caller must authenticate using a service principal.\n" +
				heading +
				"// | Identity | Support | |-|-| | User | Yes | | Service principal [/entra/foo] | Yes |",
			spnYes,
		},
		{
			"empty",
			"",
			spnUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := parseSPNSupport(tt.doc); got != tt.want {
				t.Errorf("parseSPNSupport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_categorize(t *testing.T) {
	t.Parallel()

	results := []result{
		{service: "no-decl", hasDeclaration: false},
		{service: "undetermined", hasDeclaration: true, sdkCalls: 0},
		{service: "over-marked", hasDeclaration: true, sdkCalls: 2, declared: true, nonSPNAPIs: []string{"core.GetX"}},
		{service: "review", hasDeclaration: true, sdkCalls: 2, declared: false, nonSPNAPIs: nil},
		{service: "correct-supported", hasDeclaration: true, sdkCalls: 2, declared: true, nonSPNAPIs: nil},
		{service: "correct-unsupported", hasDeclaration: true, sdkCalls: 2, declared: false, nonSPNAPIs: []string{"core.GetX"}},
	}

	b := categorize(results, nil)

	if len(b.undeclared) != 1 || b.undeclared[0].service != "no-decl" {
		t.Errorf("undeclared = %+v, want [no-decl]", b.undeclared)
	}

	if b.undetermined != 1 {
		t.Errorf("undetermined = %d, want 1", b.undetermined)
	}

	if len(b.overMarked) != 1 || b.overMarked[0].service != "over-marked" {
		t.Errorf("overMarked = %+v, want [over-marked]", b.overMarked)
	}

	if len(b.review) != 1 || b.review[0].service != "review" {
		t.Errorf("review = %+v, want [review]", b.review)
	}

	if len(b.excluded) != 0 || len(b.staleExcl) != 0 {
		t.Errorf("excluded/stale should be empty without exclusions, got %+v / %+v", b.excluded, b.staleExcl)
	}
}

func Test_categorize_exclusions(t *testing.T) {
	t.Parallel()

	results := []result{
		{service: "over-marked", hasDeclaration: true, sdkCalls: 2, declared: true, nonSPNAPIs: []string{"core.GetX"}},
		{service: "review", hasDeclaration: true, sdkCalls: 2, declared: false, nonSPNAPIs: nil},
		{service: "correct-supported", hasDeclaration: true, sdkCalls: 2, declared: true, nonSPNAPIs: nil},
	}

	exclusions := map[string]string{
		"over-marked":       "SPN-supported, stale SDK annotation",
		"correct-supported": "stale exclusion, service is already consistent",
	}

	b := categorize(results, exclusions)

	if len(b.overMarked) != 0 {
		t.Errorf("overMarked = %+v, want empty (excluded)", b.overMarked)
	}

	if len(b.excluded) != 1 || b.excluded[0].res.service != "over-marked" {
		t.Errorf("excluded = %+v, want [over-marked]", b.excluded)
	}

	if b.excluded[0].reason != "SPN-supported, stale SDK annotation" {
		t.Errorf("excluded reason = %q, want the configured reason", b.excluded[0].reason)
	}

	if len(b.review) != 1 || b.review[0].service != "review" {
		t.Errorf("review = %+v, want [review]", b.review)
	}

	if !slices.Contains(b.staleExcl, "correct-supported") {
		t.Errorf("staleExcl = %+v, want to contain correct-supported", b.staleExcl)
	}
}

func Test_classify(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		res  result
		want category
	}{
		{"undeclared", result{hasDeclaration: false}, catUndeclared},
		{"undetermined", result{hasDeclaration: true, sdkCalls: 0}, catUndetermined},
		{"over-marked", result{hasDeclaration: true, sdkCalls: 1, declared: true, nonSPNAPIs: []string{"core.X"}}, catOverMarked},
		{"review", result{hasDeclaration: true, sdkCalls: 1, declared: false}, catReview},
		{"consistent-supported", result{hasDeclaration: true, sdkCalls: 1, declared: true}, catConsistent},
		{"consistent-unsupported", result{hasDeclaration: true, sdkCalls: 1, declared: false, nonSPNAPIs: []string{"core.X"}}, catConsistent},
	}

	for _, tt := range tests {
		if got := classify(tt.res); got != tt.want {
			t.Errorf("%s: classify = %d, want %d", tt.name, got, tt.want)
		}
	}
}

func Test_result_expected_determinable(t *testing.T) {
	t.Parallel()

	if !(result{}).expected() {
		t.Error("result with no non-SPN APIs should be expected SPN-supported")
	}

	if (result{nonSPNAPIs: []string{"core.GetX"}}).expected() {
		t.Error("result with a non-SPN API should not be expected SPN-supported")
	}

	if (result{sdkCalls: 0}).determinable() {
		t.Error("result with no SDK calls should not be determinable")
	}

	if !(result{sdkCalls: 1}).determinable() {
		t.Error("result with SDK calls should be determinable")
	}
}
