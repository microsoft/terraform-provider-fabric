// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package transforms

import (
	"strings"
	"unicode"
)

type templateCompatV1CaseStyle uint8

const (
	templateCompatV1CamelStyle templateCompatV1CaseStyle = iota
	templateCompatV1PascalStyle
)

// Camel and Pascal case preserve v1.0.3 capitalization to keep generated names
// and hashes stable; the other case converters use native Sprout behavior.
// Captured expectations live in template_compat_v1_{camel,pascal}.json.
func templateCompatV1CamelCase(value string) string {
	return templateCompatV1CapitalizedCase(value, templateCompatV1CamelStyle)
}

func templateCompatV1PascalCase(value string) string {
	return templateCompatV1CapitalizedCase(value, templateCompatV1PascalStyle)
}

func templateCompatV1CapitalizedCase(value string, style templateCompatV1CaseStyle) string {
	upperNext := style == templateCompatV1PascalStyle
	seenLetter := false
	var output strings.Builder

	for _, letter := range value {
		if strings.ContainsRune(" -_", letter) {
			// Legacy conversion starts a new word only after a letter was seen.
			if seenLetter {
				upperNext = true
			}

			continue
		}

		converted := letter
		if upperNext {
			converted = unicode.ToUpper(letter)
			upperNext = false
		} else if style == templateCompatV1CamelStyle {
			converted = unicode.ToLower(letter)
		}

		_, _ = output.WriteRune(converted)

		seenLetter = seenLetter || unicode.IsLetter(letter)
	}

	return output.String()
}
