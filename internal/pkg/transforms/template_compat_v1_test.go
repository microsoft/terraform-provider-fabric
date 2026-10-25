// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package transforms_test

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func TestUnit_TemplateCompatV1(t *testing.T) {
	t.Parallel()

	tests := append(testhelp.TemplateCompatV1Cases(), testhelp.TemplateCompatV1CaseConversions(t)...)
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join(t.TempDir(), "template.txt")
			require.NoError(t, os.WriteFile(path, []byte(test.Template), 0o600))

			payload, hash, diags := transforms.SourceFileToPayload(path, transforms.ProcessingModeGoTemplate, nil, nil, transforms.TokensDelimiterCurlyBraces)
			require.False(t, diags.HasError(), "%v", diags)
			require.Equal(t, base64.StdEncoding.EncodeToString([]byte(test.Want)), payload)
			require.Equal(t, fmt.Sprintf("%x", sha256.Sum256([]byte(test.Want))), hash)
		})
	}
}

func TestUnit_TemplateCompatV1Errors(t *testing.T) {
	t.Parallel()

	for name, source := range testhelp.TemplateCompatV1Errors() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join(t.TempDir(), "template.txt")
			require.NoError(t, os.WriteFile(path, []byte(source), 0o600))

			payload, hash, diags := transforms.SourceFileToPayload(path, transforms.ProcessingModeGoTemplate, nil, nil, transforms.TokensDelimiterCurlyBraces)
			require.True(t, diags.HasError(), "invalid template unexpectedly rendered: %s", source)

			summary := common.ErrorTmplParseHeader
			if strings.HasPrefix(name, "unsupported_") || name == "parse_error" {
				summary = common.ErrorFileReadHeader
			}

			require.Equal(t, summary, diags.Errors()[0].Summary())
			require.Empty(t, payload)
			require.Empty(t, hash)
		})
	}
}

func TestUnit_TemplateFunctions(t *testing.T) {
	t.Parallel()

	tests := []testhelp.TemplateCompatV1Case{
		{Name: "append_nested_list", Template: `{{ append (list "a") (list "b") | toJson }}`, Want: `["a",["b"]]`},
		{Name: "prepend_nested_list", Template: `{{ prepend (list "a") (list "b") | toJson }}`, Want: `[["b"],"a"]`},
		{Name: "new_escape_function", Template: `{{ escape "." "a.b" }}`, Want: `a\.b`},
		{Name: "new_uuid_function", Template: `{{ uuidNil }}`, Want: "00000000-0000-0000-0000-000000000000"},
		{Name: "repeat_negative", Template: `{{ repeat -1 "a" }}`, Want: ""},
		{Name: "indent_negative", Template: `{{ indent -1 "a" }}`, Want: "a"},
		{Name: "nindent_negative", Template: `{{ nindent -1 "a" }}`, Want: "\na"},
		{Name: "ellipsis_offset", Template: `{{ ellipsisBoth 99 10 "abc" }}`, Want: ""},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join(t.TempDir(), "template.txt")
			require.NoError(t, os.WriteFile(path, []byte(test.Template), 0o600))

			payload, _, diags := transforms.SourceFileToPayload(path, transforms.ProcessingModeGoTemplate, nil, nil, transforms.TokensDelimiterCurlyBraces)
			require.False(t, diags.HasError(), "%v", diags)
			require.Equal(t, base64.StdEncoding.EncodeToString([]byte(test.Want)), payload)
		})
	}

	t.Run("division_by_zero", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "template.txt")
		require.NoError(t, os.WriteFile(path, []byte(`{{ div 1 0 }}`), 0o600))

		payload, hash, diags := transforms.SourceFileToPayload(path, transforms.ProcessingModeGoTemplate, nil, nil, transforms.TokensDelimiterCurlyBraces)
		require.True(t, diags.HasError())
		require.Contains(t, diags.Errors()[0].Detail(), "cannot divide by zero")
		require.Empty(t, payload)
		require.Empty(t, hash)
	})

	t.Run("json_and_custom_delimiters", func(t *testing.T) {
		t.Parallel()

		for _, delimiter := range []string{"{{}}", "<<>>", "@{}@", "____"} {
			t.Run(delimiter, func(t *testing.T) {
				t.Parallel()

				source := `{"value":"{{ regexReplaceAll "a" .Message "x" }}"}`
				source = strings.ReplaceAll(source, "{{", delimiter[:2])
				source = strings.ReplaceAll(source, "}}", delimiter[2:])
				path := filepath.Join(t.TempDir(), "template.json")
				require.NoError(t, os.WriteFile(path, []byte(source), 0o600))

				payload, hash, diags := transforms.SourceFileToPayload(path, transforms.ProcessingModeGoTemplate, map[string]string{"Message": "abc"}, nil, delimiter)
				require.False(t, diags.HasError(), "%v", diags)

				want := `{"value":"xbc"}`
				require.Equal(t, base64.StdEncoding.EncodeToString([]byte(want)), payload)
				require.Equal(t, fmt.Sprintf("%x", sha256.Sum256([]byte(want))), hash)
			})
		}
	})
}
