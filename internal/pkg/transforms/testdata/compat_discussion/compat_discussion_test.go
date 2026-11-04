// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package compatdiscussion_test

import (
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"

	compatdiscussion "github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms/testdata/compat_discussion"
)

func TestUnit_DiscussionLegacyDig(t *testing.T) {
	t.Parallel()

	dict := map[string]any{
		"a": map[string]any{
			"":  map[string]any{"b": "empty"},
			"b": "nested",
		},
		`a\b`: "backslash",
		"":    map[string]any{"a": "leading"},
	}

	for _, test := range []struct {
		name string
		args []any
		want any
	}{
		{"nested", []any{"a.b", dict}, "nested"},
		{"separate_keys", []any{"a", "b", dict}, "nested"},
		{"empty_segment", []any{"a..b", dict}, "empty"},
		{"leading_dot", []any{".a", dict}, "leading"},
		{"backslash", []any{`a\b`, dict}, "backslash"},
		{"missing", []any{"missing", dict}, nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			value, err := compatdiscussion.LegacyDig(test.args...)
			require.NoError(t, err)
			require.Equal(t, test.want, value)
		})
	}

	for _, args := range [][]any{
		nil,
		{"a", 1},
		{1, dict},
		{"a.b.c", dict},
	} {
		_, err := compatdiscussion.LegacyDig(args...)
		require.Error(t, err)
	}
}

func TestUnit_DiscussionLegacySeparatedCase(t *testing.T) {
	t.Parallel()

	for _, style := range []struct {
		name      string
		separator rune
		uppercase bool
		convert   func(rune) rune
	}{
		{"snake", '_', false, unicode.ToLower},
		{"kebab", '-', false, unicode.ToLower},
		{"dot", '.', false, unicode.ToLower},
		{"path", '/', false, unicode.ToLower},
		{"constant", '_', true, unicode.ToUpper},
	} {
		t.Run(style.name, func(t *testing.T) {
			t.Parallel()

			for _, test := range []struct {
				input string
				want  string
			}{
				{"0", string(style.separator) + "0"},
				{"firstName", "first" + string(style.separator) + "name"},
				{"HTMLParser", "html" + string(style.separator) + "parser"},
				{"foo2Bar", "foo" + string(style.separator) + "2bar"},
				{"a.b/c", "a.b/c"},
				{"αΒγ", "α" + string(style.separator) + "βγ"},
			} {
				want := test.want
				if style.uppercase {
					want = strings.ToUpper(want)
				}

				require.Equal(t, want, compatdiscussion.LegacySeparatedCase(test.input, style.separator, style.convert))
			}
		})
	}
}
