// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package toolutil_test

import (
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
