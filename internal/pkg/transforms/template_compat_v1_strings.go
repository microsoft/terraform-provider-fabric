// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package transforms

import (
	"strings"
	"unicode"
)

// templateCompatV1Case keeps v1.0.3's capitalization and word-boundary rules,
// including leading digit separators and the byte-based lookahead.
// Do not replace this with the current Sprout converters: their output changes
// definition content and hashes. Captured v1.0.3 expectations live in
// template_compat_v1_case.json and template_compat_v1_edges.json.
func templateCompatV1Case(name, value string) string {
	separator := map[string]rune{
		"toSnakeCase": '_', "toKebabCase": '-', "toDotCase": '.',
		"toPathCase": '/', "toConstantCase": '_',
	}[name]
	capitalizing := separator == 0
	upperNext := name == "toPascalCase"
	seenLetter := false
	var previous, lookahead rune
	var output strings.Builder

	for offset, letter := range value {
		// Byte lookahead and retaining it at EOF are intentional legacy behavior,
		// even though iteration over value itself yields Unicode runes.
		if offset+1 < len(value) {
			lookahead = rune(value[offset+1])
		}

		if strings.ContainsRune(" -_", letter) {
			if separator != 0 && previous != separator {
				_, _ = output.WriteRune(separator)
			}

			if seenLetter {
				upperNext = true
			}

			previous = separator

			continue
		}

		if separator != 0 && previous != separator {
			wordBoundary := unicode.IsUpper(letter) && output.Len() > 0 &&
				(unicode.IsLower(previous) || unicode.IsUpper(previous) && unicode.IsLower(lookahead))

			digitBoundary := unicode.IsDigit(letter) && !unicode.IsDigit(previous)
			if wordBoundary || digitBoundary {
				_, _ = output.WriteRune(separator)
			}
		}

		converted := letter

		switch {
		case capitalizing && upperNext:
			converted = unicode.ToUpper(letter)
			upperNext = false
		case name == "toConstantCase":
			converted = unicode.ToUpper(letter)
		case name != "toPascalCase":
			converted = unicode.ToLower(letter)
		default:
		}

		_, _ = output.WriteRune(converted)

		previous = letter
		seenLetter = seenLetter || unicode.IsLetter(letter)
	}

	return output.String()
}
