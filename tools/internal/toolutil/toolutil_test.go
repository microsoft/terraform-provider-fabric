// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package toolutil_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/microsoft/terraform-provider-fabric/tools/internal/toolutil"
)

func TestParseClientMethod(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		line     string
		wantRecv string
		wantMeth string
	}{
		{"pointer receiver", "func (client *WorkspacesClient) GetWorkspace(ctx context.Context) {", "WorkspacesClient", "GetWorkspace"},
		{"value receiver", "func (c ItemsClient) ListItems() {", "ItemsClient", "ListItems"},
		{"alternate receiver name", "func (d *PublishedClient) ListTags(id string) {", "PublishedClient", "ListTags"},
		{"no space before paren", "func (c *T)Do(x int) {", "T", "Do"},
		// Non-method / rejected shapes.
		{"unexported method", "func (c *T) do() {", "", ""},
		{"plain function", "func Helper() {", "", ""},
		{"not a func", "type WorkspacesClient struct {", "", ""},
		{"empty receiver", "func () Foo() {", "", ""},
		{"no method params", "func (c *T) Field", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotRecv, gotMeth := toolutil.ParseClientMethod(tt.line)
			if gotRecv != tt.wantRecv || gotMeth != tt.wantMeth {
				t.Errorf("ParseClientMethod(%q) = (%q, %q), want (%q, %q)",
					tt.line, gotRecv, gotMeth, tt.wantRecv, tt.wantMeth)
			}
		})
	}
}

func TestDocCommentAbove(t *testing.T) {
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
			"// GetFoo retrieves a foo.\n// This API is currently in preview.",
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

			if got := toolutil.DocCommentAbove(tt.lines, tt.declLine); got != tt.want {
				t.Errorf("DocCommentAbove(%d) = %q, want %q", tt.declLine, got, tt.want)
			}
		})
	}
}

func TestIsItemCRUD(t *testing.T) {
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

			if got := toolutil.IsItemCRUD(tt.method, tt.item); got != tt.want {
				t.Errorf("IsItemCRUD(%q, %q) = %v, want %v", tt.method, tt.item, got, tt.want)
			}
		})
	}
}

func TestFabricRoot(t *testing.T) {
	t.Parallel()

	coreFile := filepath.Join("home", "user", "sdk", "fabric", "core", "zz_generated_models.go")
	want := filepath.Join("home", "user", "sdk", "fabric")

	if got := toolutil.FabricRoot(coreFile); got != want {
		t.Errorf("FabricRoot(%q) = %q, want %q", coreFile, got, want)
	}

	notCore := filepath.Join("home", "user", "sdk", "fabric", "admin", "client.go")
	if got := toolutil.FabricRoot(notCore); got != "" {
		t.Errorf("FabricRoot(%q) = %q, want empty (parent dir is not core)", notCore, got)
	}
}

func TestScanItemCRUD(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	content := "package notebook\n\n" +
		"// GetNotebook gets a notebook.\n" +
		"// This API is currently in preview.\n" +
		"func (client *ItemsClient) GetNotebook(ctx context.Context) (Response, error) {}\n\n" +
		"// ListNotebooks lists notebooks.\n" +
		"func (client *ItemsClient) ListNotebooks(ctx context.Context) (Response, error) {}\n\n" +
		"// ApplyLabel applies a label (not the item).\n" +
		"func (client *ItemsClient) ApplyLabel(ctx context.Context) (Response, error) {}\n\n" +
		"// getInternal is unexported and must be ignored.\n" +
		"func (client *ItemsClient) getInternal() {}\n"

	err := os.WriteFile(filepath.Join(dir, "items_client.go"), []byte(content), 0o600)
	if err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	methods := toolutil.NewDocCache().ScanItemCRUD(dir, "Notebook")

	names := make([]string, 0, len(methods))
	for _, m := range methods {
		names = append(names, m.Name)
	}

	wantNames := []string{"GetNotebook", "ListNotebooks"}
	if !slices.Equal(names, wantNames) {
		t.Fatalf("CRUD methods = %v, want %v", names, wantNames)
	}

	if got := methods[0].Doc; got != "// GetNotebook gets a notebook.\n// This API is currently in preview." {
		t.Errorf("GetNotebook doc = %q, want the preview marker block", got)
	}

	if got := methods[1].Doc; got != "// ListNotebooks lists notebooks." {
		t.Errorf("ListNotebooks doc = %q, want its own doc line", got)
	}
}

func TestLoadExclusions(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		path := writeExclusions(t, "exclusions:\n  - service: foo\n    reason: GA, stale marker\n")

		m, err := toolutil.LoadExclusions(path)
		if err != nil {
			t.Fatalf("LoadExclusions() error = %v", err)
		}

		if m["foo"] != "GA, stale marker" {
			t.Errorf("m[foo] = %q, want the configured reason", m["foo"])
		}
	})

	t.Run("missing service", func(t *testing.T) {
		t.Parallel()

		path := writeExclusions(t, "exclusions:\n  - reason: no service\n")

		_, err := toolutil.LoadExclusions(path)
		if err == nil {
			t.Error("LoadExclusions() error = nil, want error for missing service")
		}
	})

	t.Run("missing reason", func(t *testing.T) {
		t.Parallel()

		path := writeExclusions(t, "exclusions:\n  - service: foo\n")

		_, err := toolutil.LoadExclusions(path)
		if err == nil {
			t.Error("LoadExclusions() error = nil, want error for missing reason")
		}
	})

	t.Run("malformed yaml", func(t *testing.T) {
		t.Parallel()

		path := writeExclusions(t, "exclusions: [::::\n")

		_, err := toolutil.LoadExclusions(path)
		if err == nil {
			t.Error("LoadExclusions() error = nil, want parse error")
		}
	})

	t.Run("nonexistent path", func(t *testing.T) {
		t.Parallel()

		_, err := toolutil.LoadExclusions(filepath.Join(t.TempDir(), "nope.yaml"))
		if err == nil {
			t.Error("LoadExclusions() error = nil, want read error")
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
