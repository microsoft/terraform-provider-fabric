// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package testhelp

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

// TemplateCompatV1Case describes the provider's Sprout v1.0.3 template contract.
type TemplateCompatV1Case struct {
	Name     string `json:"name"`
	Template string `json:"template"`
	Want     string `json:"want"`
}

// TemplateCompatV1CaseConversions loads representative Sprout v1.0.3 goldens:
// capitalization boundaries, digit/acronym handling, separators, and Unicode.
// Keep one case per rule and converter, rather than a string cross-product.
func TemplateCompatV1CaseConversions(t *testing.T) []TemplateCompatV1Case {
	t.Helper()

	data, err := os.ReadFile(GetFixturesDirPath("template_compat_v1_case.json"))
	if err != nil {
		t.Fatal(err)
	}

	var cases []TemplateCompatV1Case

	err = json.Unmarshal(data, &cases)
	if err != nil {
		t.Fatal(err)
	}

	return cases
}

// TemplateCompatV1Errors exercises invalid calls and the existing registry restrictions.
func TemplateCompatV1Errors() map[string]string {
	return map[string]string{
		"get_count":                     `{{ get }}`,
		"get_type":                      `{{ get 1 2 }}`,
		"get_key_type":                  `{{ get (dict) 1 }}`,
		"set_count":                     `{{ set (dict) }}`,
		"set_type":                      `{{ set "a" "b" 1 }}`,
		"unset_count":                   `{{ unset }}`,
		"unset_type":                    `{{ unset 1 2 }}`,
		"has_key_count":                 `{{ hasKey }}`,
		"has_key_type":                  `{{ hasKey 1 2 }}`,
		"pick_count":                    `{{ pick }}`,
		"pick_type":                     `{{ pick "a" 1 }}`,
		"pick_key_type":                 `{{ pick (dict) 1 }}`,
		"omit_count":                    `{{ omit }}`,
		"omit_type":                     `{{ omit "a" 1 }}`,
		"omit_key_type":                 `{{ omit (dict) 1 }}`,
		"append_count":                  `{{ append }}`,
		"append_type":                   `{{ append 1 2 }}`,
		"prepend_count":                 `{{ prepend }}`,
		"prepend_type":                  `{{ prepend 1 2 }}`,
		"slice_count":                   `{{ slice }}`,
		"slice_type":                    `{{ slice 0 1 }}`,
		"slice_bounds":                  `{{ slice (list "a") 9 }}`,
		"without_count":                 `{{ without }}`,
		"without_type":                  `{{ without 1 2 }}`,
		"dig_count":                     `{{ dig }}`,
		"dig_type":                      `{{ dig "a" 1 }}`,
		"dig_key_type":                  `{{ dig 1 (dict) }}`,
		"dig_scalar":                    `{{ dict "a" 1 | dig "a.b" }}`,
		"dig_no_default_argument":       `{{ dict "a" "value" | dig "a" "fallback" }}`,
		"regex_find":                    `{{ regexFind "[" "abc" }}`,
		"regex_find_all":                `{{ regexFindAll "[" "abc" -1 }}`,
		"regex_match":                   `{{ regexMatch "[" "abc" }}`,
		"regex_split":                   `{{ regexSplit "[" "abc" -1 }}`,
		"regex_replace":                 `{{ regexReplaceAll "[" "abc" "x" }}`,
		"regex_literal":                 `{{ regexReplaceAllLiteral "[" "abc" "x" }}`,
		"regex_groups":                  `{{ regexFindGroups "[" "abc" }}`,
		"regex_all_groups":              `{{ regexFindAllGroups "[" -1 "abc" }}`,
		"regex_named":                   `{{ regexFindNamed "[" "abc" }}`,
		"regex_all_named":               `{{ regexFindAllNamed "[" -1 "abc" }}`,
		"must_append":                   `{{ mustAppend 1 2 }}`,
		"must_prepend":                  `{{ mustPrepend 1 2 }}`,
		"must_slice":                    `{{ mustSlice 0 1 }}`,
		"must_without":                  `{{ mustWithout 1 2 }}`,
		"must_regex_find":               `{{ mustRegexFind "[" "abc" }}`,
		"must_regex_find_all":           `{{ mustRegexFindAll "[" "abc" -1 }}`,
		"must_regex_match":              `{{ mustRegexMatch "[" "abc" }}`,
		"must_regex_split":              `{{ mustRegexSplit "[" "abc" -1 }}`,
		"must_regex_replace":            `{{ mustRegexReplaceAll "[" "abc" "x" }}`,
		"must_regex_literal":            `{{ mustRegexReplaceAllLiteral "[" "abc" "x" }}`,
		"regex_find_all_new_order":      `{{ regexFindAll "." -1 "abc" }}`,
		"regex_split_new_order":         `{{ regexSplit "," -1 "a,b" }}`,
		"must_regex_find_all_new_order": `{{ mustRegexFindAll "." -1 "abc" }}`,
		"must_regex_split_new_order":    `{{ mustRegexSplit "," -1 "a,b" }}`,
		"base64_invalid":                `{{ base64Decode "!" }}`,
		"base32_invalid":                `{{ base32Decode "!" }}`,
		"unsupported_env":               `{{ env "HOME" }}`,
		"unsupported_expand":            `{{ expandEnv "$HOME" }}`,
		"unsupported_crypto":            `{{ derivePassword 1 "long" "secret" "user" "host" }}`,
		"unsupported_file":              `{{ pathBase "/tmp/test" }}`,
		"unsupported_reflect":           `{{ kindOf "value" }}`,
		"unsupported_sprig":             `{{ upper "value" }}`,
		"unsupported_push":              `{{ push (list) "a" }}`,
		"unsupported_encoding":          `{{ b64enc "value" }}`,
		"parse_error":                   `{{ .Message! }}`,
	}
}

// TemplateCompatV1Cases is shared by renderer and Notebook lifecycle tests.
func TemplateCompatV1Cases() []TemplateCompatV1Case {
	cases := []TemplateCompatV1Case{
		{"get_missing", `{{ get (dict) "missing" }}`, ""},
		{"set_map_value", `{{ $d := dict }}{{ $_ := set $d "a" (dict "b" "value") }}{{ $d | dig "a.b" }}`, "value"},
		{"has_key_missing", `{{ dict | hasKey "a" }}`, "false"},
		{"append_accumulation", `{{ $l := list }}{{ range list "a" "b" }}{{ $l = append $l . }}{{ end }}{{ $l | join "," }}`, "a,b"},
		{"append_no_mutation", `{{ $l := list "a" }}{{ $_ := append $l "b" }}{{ $l | join "," }}`, "a"},
		{"slice_one_index", `{{ slice (list "a" "b") 1 | join "," }}`, "b"},
		{"slice_no_index", `{{ slice (list "a" "b") | join "," }}`, "a,b"},
		{"dig_nested", `{{ dict "a" (dict "b" "value") | dig "a.b" }}`, "value"},
		{"dig_separate_keys", `{{ dict "a" (dict "b" "value") | dig "a" "b" }}`, "value"},
		{"dig_empty_segment", `{{ dict "a" (dict "" (dict "b" "value")) | dig "a..b" }}`, "value"},
		{"dig_leading_dot", `{{ dict "" (dict "a" "value") | dig ".a" }}`, "value"},
		{"dig_trailing_dot", `{{ dict "a" (dict "" "value") | dig "a." }}`, "value"},
		{"dig_literal_backslash", `{{ dict "a\\b" "value" | dig "a\\b" }}`, "value"},
		{"dig_backslash_before_dot", `{{ dict "a\\" (dict "b" "value") | dig "a\\.b" }}`, "value"},
		{"dig_missing", `{{ dict | dig "a" | default "fallback" }}`, "fallback"},
		{"regex_find", `{{ regexFind "[a-z]+" "ab12cd" }}`, "ab"},
		{"regex_find_all", `{{ regexFindAll "[a-z]+" "ab12cd" -1 | join "," }}`, "ab,cd"},
		{"regex_find_all_limit", `{{ regexFindAll "[a-z]+" "ab12cd" 1 | join "," }}`, "ab"},
		{"regex_find_all_zero", `{{ regexFindAll "[a-z]+" "ab12cd" 0 | len }}`, "0"},
		{"regex_split", `{{ regexSplit "," "a,b,c" -1 | join ":" }}`, "a:b:c"},
		{"regex_replace", `{{ regexReplaceAll "a(x*)" "-ab-axxb-" "${1}W" }}`, "-Wb-xxWb-"},
		{"regex_replace_literal", `{{ regexReplaceAllLiteral "a(x*)" "-ab-axxb-" "${1}W" }}`, "-${1}Wb-${1}Wb-"},
		{"regex_match", `{{ regexMatch "^ab" "abcd" }}`, "true"},
		{"regex_no_match", `{{ regexFind "[0-9]+" "abc" }}`, ""},
		{"regex_groups", `{{ regexFindGroups "(a)(b)" "ab" | join "," }}`, "ab,a,b"},
		{"regex_all_groups", `{{ regexFindAllGroups "(a)" -1 "aba" | toJson }}`, `[["a","a"],["a","a"]]`},
		{"regex_named", `{{ regexFindNamed "(?P<letter>a)" "ab" | toJson }}`, `{"letter":"a"}`},
		{"regex_all_named", `{{ regexFindAllNamed "(?P<letter>a)" -1 "aba" | toJson }}`, `[{"letter":"a"},{"letter":"a"}]`},
		{"regex_quote", `{{ regexQuoteMeta "a.b" }}`, `a\.b`},
		{"must_regex_find_alias", `{{ mustRegexFind "[a-z]+" "ab12cd" }}`, "ab"},
		{"must_regex_match_alias", `{{ mustRegexMatch "^ab" "abcd" }}`, "true"},
		{"must_regex_alias", `{{ mustRegexFindAll "[a-z]+" "ab12cd" -1 | join "," }}`, "ab,cd"},
		{"must_regex_split_alias", `{{ mustRegexSplit "," "a,b" -1 | join ":" }}`, "a:b"},
		{"must_regex_replace_alias", `{{ mustRegexReplaceAll "a" "abc" "x" }}`, "xbc"},
		{"must_regex_literal_alias", `{{ mustRegexReplaceAllLiteral "a" "abc" "$1" }}`, "$1bc"},
		{"substr_zero", `{{ substr 0 0 "hello" }}`, ""},
		{"substr_empty_range", `{{ substr 2 0 "hello" }}`, ""},
		{"substr_negative_end_zero", `{{ substr 0 -5 "hello" }}`, ""},
		{"substr_bounds", `{{ substr -1 99 "hello" }}`, "o"},
		{"ellipsis_unicode_byte_offset", `{{ ellipsisBoth 2 10 "\u00e9abc" }}`, "abc"},
		{"ellipsis_unicode_past_rune_count", `{{ ellipsisBoth 4 10 "\u00e9\u00e9a" }}`, "a"},
		{"ellipsis_short_width", `{{ ellipsisBoth 2 3 "\u00e9abc" }}`, "\u00e9abc"},
		{"ellipsis_short_both_width", `{{ ellipsisBoth 2 6 "\u00e9abc" }}`, "\u00e9abc"},
		{"ellipsis_truncated", `{{ ellipsisBoth 2 7 "abcdefghij" }}`, "...c..."},
		{"camel_case", `{{ toCamelCase "firstName" }}`, "firstname"},
		{"pascal_case", `{{ toPascalCase "hELLO world" }}`, "HELLOWorld"},
		{"snake_case", `{{ toSnakeCase "2FOo" }}`, "_2f_oo"},
		{"kebab_case", `{{ toKebabCase "2FOo" }}`, "-2f-oo"},
		{"dot_case", `{{ toDotCase "2FOo" }}`, ".2f.oo"},
		{"path_case", `{{ toPathCase "2FOo" }}`, "/2f/oo"},
		{"constant_case", `{{ toConstantCase "2FOo" }}`, "_2F_OO"},
		{"title_case", `{{ toTitleCase "HELLO it's" }}`, "Hello It's"},
		{"base64_decode", `{{ base64Decode "aGVsbG8=" }}`, "hello"},
		{"base32_decode", `{{ base32Decode "NBSWY3DP" }}`, "hello"},
		{"merge_zero", `{{ mergeOverwrite (dict "a" 1) (dict "a" 0) | toJson }}`, `{"a":0}`},
		{"date_timezone", `{{ date "15:04 -0700" (toDate "2006-01-02T15:04:05Z07:00" "2026-01-02T03:04:05+02:00") }}`, "03:04 +0200"},
		{"random_zero", `{{ randAlpha 0 }}`, ""},
		{"numeric_alias/addf", `{{ addf 1.5 2 }}`, "3.5"},
		{"numeric_alias/add1f", `{{ add1f 1.5 }}`, "2.5"},
		{"numeric_alias/subf", `{{ subf 3.5 2 }}`, "1.5"},
	}

	return append(templateCompatV1SyntaxCases(), cases...)
}

// templateCompatV1SyntaxCases checks each unambiguous function in all three
// supported forms. Aliases use the same matrix as their canonical function.
func templateCompatV1SyntaxCases() []TemplateCompatV1Case {
	tests := []struct {
		names    []string
		legacy   string
		direct   string
		pipeline string
		want     string
	}{
		{
			names:    []string{"get"},
			legacy:   `{{ %s (dict "a" "value") "a" }}`,
			direct:   `{{ %s "a" (dict "a" "value") }}`,
			pipeline: `{{ dict "a" "value" | %s "a" }}`,
			want:     "value",
		},
		{
			names:    []string{"set"},
			legacy:   `{{ $d := dict }}{{ $_ := %s $d "a" "value" }}{{ get $d "a" }}`,
			direct:   `{{ $d := dict }}{{ $_ := %s "a" "value" $d }}{{ get $d "a" }}`,
			pipeline: `{{ $d := dict }}{{ $_ := $d | %s "a" "value" }}{{ get $d "a" }}`,
			want:     "value",
		},
		{
			names:    []string{"unset"},
			legacy:   `{{ $d := dict "a" 1 }}{{ $_ := %s $d "a" }}{{ len $d }}`,
			direct:   `{{ $d := dict "a" 1 }}{{ $_ := %s "a" $d }}{{ len $d }}`,
			pipeline: `{{ $d := dict "a" 1 }}{{ $_ := $d | %s "a" }}{{ len $d }}`,
			want:     "0",
		},
		{
			names:    []string{"hasKey"},
			legacy:   `{{ %s (dict "a" 0) "a" }}`,
			direct:   `{{ %s "a" (dict "a" 0) }}`,
			pipeline: `{{ dict "a" 0 | %s "a" }}`,
			want:     "true",
		},
		{
			names:    []string{"pick"},
			legacy:   `{{ %s (dict "a" 1 "b" 2 "c" 3) "a" "b" | toJson }}`,
			direct:   `{{ %s "a" "b" (dict "a" 1 "b" 2 "c" 3) | toJson }}`,
			pipeline: `{{ dict "a" 1 "b" 2 "c" 3 | %s "a" "b" | toJson }}`,
			want:     `{"a":1,"b":2}`,
		},
		{
			names:    []string{"omit"},
			legacy:   `{{ %s (dict "a" 1 "b" 2 "c" 3) "a" "b" | toJson }}`,
			direct:   `{{ %s "a" "b" (dict "a" 1 "b" 2 "c" 3) | toJson }}`,
			pipeline: `{{ dict "a" 1 "b" 2 "c" 3 | %s "a" "b" | toJson }}`,
			want:     `{"c":3}`,
		},
		{
			names:    []string{"append", "mustAppend"},
			legacy:   `{{ %s (list "a") "b" | join "," }}`,
			direct:   `{{ %s "b" (list "a") | join "," }}`,
			pipeline: `{{ list "a" | %s "b" | join "," }}`,
			want:     "a,b",
		},
		{
			names:    []string{"prepend", "mustPrepend"},
			legacy:   `{{ %s (list "a") "b" | join "," }}`,
			direct:   `{{ %s "b" (list "a") | join "," }}`,
			pipeline: `{{ list "a" | %s "b" | join "," }}`,
			want:     "b,a",
		},
		{
			names:    []string{"slice", "mustSlice"},
			legacy:   `{{ %s (list "a" "b" "c") 1 3 | join "," }}`,
			direct:   `{{ %s 1 3 (list "a" "b" "c") | join "," }}`,
			pipeline: `{{ list "a" "b" "c" | %s 1 3 | join "," }}`,
			want:     "b,c",
		},
		{
			names:    []string{"without", "mustWithout"},
			legacy:   `{{ %s (list "a" "b" "c") "a" "c" | join "," }}`,
			direct:   `{{ %s "a" "c" (list "a" "b" "c") | join "," }}`,
			pipeline: `{{ list "a" "b" "c" | %s "a" "c" | join "," }}`,
			want:     "b",
		},
	}

	var cases []TemplateCompatV1Case

	for _, test := range tests {
		for _, name := range test.names {
			for _, form := range []struct{ name, call string }{
				{"legacy", test.legacy},
				{"direct", test.direct},
				{"pipeline", test.pipeline},
			} {
				cases = append(cases, TemplateCompatV1Case{
					Name:     "syntax/" + name + "/" + form.name,
					Template: fmt.Sprintf(form.call, name),
					Want:     test.want,
				})
			}
		}
	}

	return cases
}
