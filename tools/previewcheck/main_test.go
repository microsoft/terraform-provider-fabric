// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func Test_containsMarker(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		doc  string
		want bool
	}{
		{"currently in preview", "// This API is currently in preview.", true},
		{"part of a preview release", "// GetFoo is part of a Preview release.", true},
		{"mixed case marker", "// IS CURRENTLY IN PREVIEW", true},
		{"marker mid-sentence", "// note: the endpoint is part of a preview release for now", true},
		{"no marker", "// GetFoo retrieves a foo.", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := containsMarker(tt.doc); got != tt.want {
				t.Errorf("containsMarker(%q) = %v, want %v", tt.doc, got, tt.want)
			}
		})
	}
}

func Test_docCommentAbove(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		lines    []string
		declLine int
		want     string
	}{
		{
			"contiguous doc block",
			[]string{
				"// GetFoo retrieves a foo.",
				"// This API is currently in preview.",
				"func (client *FooClient) GetFoo() {}",
			},
			3,
			"// This API is currently in preview.\n// GetFoo retrieves a foo.",
		},
		{
			"blank line breaks the block",
			[]string{
				"// GetFoo retrieves a foo.",
				"",
				"func (client *FooClient) GetFoo() {}",
			},
			3,
			"",
		},
		{
			"non-comment line breaks the block",
			[]string{
				"import \"context\"",
				"func (client *FooClient) GetFoo() {}",
			},
			2,
			"",
		},
		{
			"declaration on first line does not panic",
			[]string{"func (client *FooClient) GetFoo() {}"},
			1,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := docCommentAbove(tt.lines, tt.declLine); got != tt.want {
				t.Errorf("docCommentAbove(%d) = %q, want %q", tt.declLine, got, tt.want)
			}
		})
	}
}

func Test_isItemCRUD(t *testing.T) {
	t.Parallel()

	tests := []struct {
		method string
		item   string
		want   bool
	}{
		{"GetNotebook", "Notebook", true},
		{"ListNotebooks", "Notebook", true},
		{"BeginCreateNotebook", "Notebook", true},
		{"UpdateNotebook", "Notebook", true},
		{"DeleteNotebook", "Notebook", true},
		// Item name matches but the verb is not a CRUD verb.
		{"PublishNotebook", "Notebook", false},
		// CRUD verb but the method is about a different entity.
		{"GetWorkspace", "Notebook", false},
		// Neither item nor CRUD verb.
		{"ApplyTags", "Notebook", false},
	}

	for _, tt := range tests {
		t.Run(tt.method+"/"+tt.item, func(t *testing.T) {
			t.Parallel()

			if got := isItemCRUD(tt.method, tt.item); got != tt.want {
				t.Errorf("isItemCRUD(%q, %q) = %v, want %v", tt.method, tt.item, got, tt.want)
			}
		})
	}
}

func Test_scanItemMethods(t *testing.T) {
	t.Parallel()

	lines := []string{
		"// GetNotebook gets a notebook.",
		"// This API is currently in preview.",
		"func (client *ItemsClient) GetNotebook(ctx context.Context) (Response, error) {",
		"",
		"// ListNotebooks lists notebooks.",
		"func (client *ItemsClient) ListNotebooks(ctx context.Context) (Response, error) {",
		"",
		"// ApplyLabel applies a label (not the item).",
		"func (client *ItemsClient) ApplyLabel(ctx context.Context) (Response, error) {",
		"",
		"// getInternal is unexported and must be ignored.",
		"func (client *ItemsClient) getInternal() {}",
	}

	crud, preview := scanItemMethods(lines, "Notebook")

	if crud != 2 {
		t.Errorf("crud count = %d, want 2 (GetNotebook, ListNotebooks)", crud)
	}

	want := []string{"GetNotebook"}
	if !slices.Equal(preview, want) {
		t.Errorf("preview methods = %v, want %v", preview, want)
	}
}

func Test_categorize(t *testing.T) {
	t.Parallel()

	results := []result{
		{service: "no-decl", hasDeclaration: false},
		{service: "undetermined", hasDeclaration: true, sdkCalls: 0},
		{service: "under-marked", hasDeclaration: true, sdkCalls: 2, declared: false, previewAPIs: []string{"core.GetX"}},
		{service: "review", hasDeclaration: true, sdkCalls: 2, declared: true, previewAPIs: nil},
		{service: "correct-preview", hasDeclaration: true, sdkCalls: 2, declared: true, previewAPIs: []string{"core.GetX"}},
		{service: "correct-ga", hasDeclaration: true, sdkCalls: 2, declared: false, previewAPIs: nil},
	}

	b := categorize(results, nil)

	if len(b.undeclared) != 1 || b.undeclared[0].service != "no-decl" {
		t.Errorf("undeclared = %+v, want [no-decl]", b.undeclared)
	}

	if b.undetermined != 1 {
		t.Errorf("undetermined = %d, want 1", b.undetermined)
	}

	if len(b.underMarked) != 1 || b.underMarked[0].service != "under-marked" {
		t.Errorf("underMarked = %+v, want [under-marked]", b.underMarked)
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
		{service: "under-marked", hasDeclaration: true, sdkCalls: 2, declared: false, previewAPIs: []string{"core.GetX"}},
		{service: "review", hasDeclaration: true, sdkCalls: 2, declared: true, previewAPIs: nil},
		{service: "correct-ga", hasDeclaration: true, sdkCalls: 2, declared: false, previewAPIs: nil},
	}

	exclusions := map[string]string{
		"under-marked": "confirmed GA, stale SDK marker",
		"correct-ga":   "stale exclusion, service is already consistent",
	}

	b := categorize(results, exclusions)

	// under-marked is failing + excluded -> moves to excluded, not underMarked.
	if len(b.underMarked) != 0 {
		t.Errorf("underMarked = %+v, want empty (excluded)", b.underMarked)
	}

	if len(b.excluded) != 1 || b.excluded[0].res.service != "under-marked" {
		t.Errorf("excluded = %+v, want [under-marked]", b.excluded)
	}

	if b.excluded[0].reason != "confirmed GA, stale SDK marker" {
		t.Errorf("excluded reason = %q, want the configured reason", b.excluded[0].reason)
	}

	// review is not excluded -> still reported.
	if len(b.review) != 1 || b.review[0].service != "review" {
		t.Errorf("review = %+v, want [review]", b.review)
	}

	// correct-ga is consistent, so its exclusion is stale.
	if !slices.Contains(b.staleExcl, "correct-ga") {
		t.Errorf("staleExcl = %+v, want to contain correct-ga", b.staleExcl)
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
		{"under-marked", result{hasDeclaration: true, sdkCalls: 1, declared: false, previewAPIs: []string{"core.X"}}, catUnderMarked},
		{"review", result{hasDeclaration: true, sdkCalls: 1, declared: true}, catReview},
		{"consistent-preview", result{hasDeclaration: true, sdkCalls: 1, declared: true, previewAPIs: []string{"core.X"}}, catConsistent},
		{"consistent-ga", result{hasDeclaration: true, sdkCalls: 1, declared: false}, catConsistent},
	}

	for _, tt := range tests {
		if got := classify(tt.res); got != tt.want {
			t.Errorf("%s: classify = %d, want %d", tt.name, got, tt.want)
		}
	}
}

func Test_loadExclusions(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		path := writeExclusions(t, "exclusions:\n  - service: foo\n    reason: GA, stale marker\n")

		m, err := loadExclusions(path)
		if err != nil {
			t.Fatalf("loadExclusions() error = %v", err)
		}

		if m["foo"] != "GA, stale marker" {
			t.Errorf("m[foo] = %q, want the configured reason", m["foo"])
		}
	})

	t.Run("missing service", func(t *testing.T) {
		t.Parallel()

		path := writeExclusions(t, "exclusions:\n  - reason: no service\n")

		_, err := loadExclusions(path)
		if err == nil {
			t.Error("loadExclusions() error = nil, want error for missing service")
		}
	})

	t.Run("missing reason", func(t *testing.T) {
		t.Parallel()

		path := writeExclusions(t, "exclusions:\n  - service: foo\n")

		_, err := loadExclusions(path)
		if err == nil {
			t.Error("loadExclusions() error = nil, want error for missing reason")
		}
	})

	t.Run("malformed yaml", func(t *testing.T) {
		t.Parallel()

		path := writeExclusions(t, "exclusions: [::::\n")

		_, err := loadExclusions(path)
		if err == nil {
			t.Error("loadExclusions() error = nil, want parse error")
		}
	})

	t.Run("nonexistent path", func(t *testing.T) {
		t.Parallel()

		_, err := loadExclusions(filepath.Join(t.TempDir(), "nope.yaml"))
		if err == nil {
			t.Error("loadExclusions() error = nil, want read error")
		}
	})
}

func writeExclusions(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "exclusions.yaml")

	err := os.WriteFile(path, []byte(content), 0o600)
	if err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	return path
}

func Test_result_expected_determinable(t *testing.T) {
	t.Parallel()

	if (result{}).expected() {
		t.Error("empty result should not be expected preview")
	}

	if !(result{previewAPIs: []string{"core.GetX"}}).expected() {
		t.Error("result with previewAPIs should be expected preview")
	}

	if (result{sdkCalls: 0}).determinable() {
		t.Error("result with no SDK calls should not be determinable")
	}

	if !(result{sdkCalls: 1}).determinable() {
		t.Error("result with SDK calls should be determinable")
	}
}

func Test_fabricRoot(t *testing.T) {
	t.Parallel()

	coreFile := filepath.Join("home", "user", "sdk", "fabric", "core", "zz_generated_models.go")
	want := filepath.Join("home", "user", "sdk", "fabric")

	if got := fabricRoot(coreFile); got != want {
		t.Errorf("fabricRoot(%q) = %q, want %q", coreFile, got, want)
	}

	notCore := filepath.Join("home", "user", "sdk", "fabric", "admin", "client.go")
	if got := fabricRoot(notCore); got != "" {
		t.Errorf("fabricRoot(%q) = %q, want empty (parent dir is not core)", notCore, got)
	}
}
