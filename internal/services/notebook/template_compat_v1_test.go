// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package notebook_test

import (
	"crypto/sha256"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

func TestUnit_NotebookResource_TemplateCompatV1(t *testing.T) {
	server := fakes.NewFakeServer()
	workspaceID := testhelp.RandomUUID()
	steps := notebookTemplateCompatV1Steps(t, workspaceID)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, server.ServerFactory, nil, steps))
}

func TestAcc_NotebookResource_TemplateCompatV1(t *testing.T) {
	require.True(t, testhelp.IsWellKnownDataAvailable(), "Acceptance tests require FABRIC_TESTACC_WELLKNOWN or the well-known fixture file")

	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)
	steps := notebookTemplateCompatV1Steps(t, workspaceID)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, steps))
}

func notebookTemplateCompatV1Steps(t *testing.T, workspaceID string) []resource.TestStep {
	t.Helper()

	directory := t.TempDir()
	root, err := os.OpenRoot(directory)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, root.Close())
	})

	displayName := testhelp.RandomName()
	scaffold, err := os.ReadFile(testhelp.GetFixturesDirPath("notebook", "template_compat_v1.py.tmpl"))
	require.NoError(t, err)

	notebookScaffold := strings.ReplaceAll(string(scaffold), "\r\n", "\n")
	const bodyMarker = "# TEMPLATE_COMPAT_V1_BODY\n"
	require.Equal(t, 1, strings.Count(notebookScaffold, bodyMarker))

	var source, expected strings.Builder

	tests := append(testhelp.TemplateCompatV1Cases(), testhelp.TemplateCompatV1CaseConversions(t)...)
	for _, test := range tests {
		fmt.Fprintf(&source, "# %s: %s\n", test.Name, test.Template)
		fmt.Fprintf(&expected, "# %s: %s\n", test.Name, test.Want)
	}

	_, _ = source.WriteString("print(\"{{ .Message }}\")\n")
	templateSource := strings.Replace(notebookScaffold, bodyMarker, source.String(), 1)

	legacyPath := filepath.Join(directory, "legacy.py.tmpl")

	require.NoError(t, root.WriteFile("legacy.py.tmpl", []byte(templateSource), 0o600))

	customPath := filepath.Join(directory, "custom.py.tmpl")
	customSource := strings.ReplaceAll(strings.ReplaceAll(templateSource, "{{", "<<"), "}}", ">>")
	require.NoError(t, root.WriteFile("custom.py.tmpl", []byte(customSource), 0o600))

	config := func(path, message, delimiter string) string {
		return at.CompileConfig(testResourceItemHeader, map[string]any{
			"workspace_id": workspaceID,
			"display_name": displayName,
			"format":       "py",
			"definition": map[string]any{
				`"notebook-content.py"`: map[string]any{
					"source":           filepath.ToSlash(path),
					"tokens":           map[string]any{"Message": message},
					"tokens_delimiter": delimiter,
				},
			},
		})
	}
	check := func(message string) resource.TestCheckFunc {
		body := expected.String() + fmt.Sprintf("print(%q)\n", message)
		content := strings.Replace(notebookScaffold, bodyMarker, body, 1)
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))

		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
			resource.TestCheckResourceAttr(testResourceItemFQN, "definition.notebook-content.py.source_content_sha256", hash),
		)
	}

	failures := testhelp.TemplateCompatV1Errors()

	steps := make([]resource.TestStep, 0, len(failures)+4)
	for _, name := range slices.Sorted(maps.Keys(failures)) {
		path := filepath.Join(directory, name+".py.tmpl")
		require.NoError(t, root.WriteFile(name+".py.tmpl", []byte(failures[name]), 0o600))

		summary := common.ErrorTmplParseHeader
		if strings.HasPrefix(name, "unsupported_") || name == "parse_error" {
			summary = common.ErrorFileReadHeader
		}

		steps = append(steps, resource.TestStep{
			ResourceName: testResourceItemFQN,
			Config:       config(path, "before", "{{}}"),
			PlanOnly:     true,
			ExpectError:  regexp.MustCompile(regexp.QuoteMeta(summary)),
		})
	}

	steps = append(steps,
		resource.TestStep{
			ResourceName: testResourceItemFQN,
			Config:       config(legacyPath, "before", "{{}}"),
			Check:        check("before"),
		},
		resource.TestStep{
			ResourceName: testResourceItemFQN,
			Config:       config(legacyPath, "before", "{{}}"),
			PlanOnly:     true,
		},
		resource.TestStep{
			ResourceName: testResourceItemFQN,
			Config:       config(legacyPath, "after", "{{}}"),
			Check:        check("after"),
		},
		resource.TestStep{
			ResourceName: testResourceItemFQN,
			Config:       config(customPath, "custom", "<<>>"),
			Check:        check("custom"),
		},
	)

	return steps
}
