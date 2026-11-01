// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package compatdiscussion

import (
	"strings"
	"unicode"
)

// LegacySeparatedCase is a discussion-only implementation for the five
// separated case converters. Callers supply their separator and case converter.
func LegacySeparatedCase(value string, separator rune, convert func(rune) rune) string {
	var previous, lookahead rune
	var output strings.Builder

	for offset, letter := range value {
		// v1.0.3 used byte lookahead and retained the previous value at EOF.
		if offset+1 < len(value) {
			lookahead = rune(value[offset+1])
		}

		if strings.ContainsRune(" -_", letter) {
			if previous != separator {
				_, _ = output.WriteRune(separator)
			}

			previous = separator

			continue
		}

		if previous != separator {
			wordBoundary := unicode.IsUpper(letter) && output.Len() > 0 &&
				(unicode.IsLower(previous) || unicode.IsUpper(previous) && unicode.IsLower(lookahead))

			digitBoundary := unicode.IsDigit(letter) && !unicode.IsDigit(previous)
			if wordBoundary || digitBoundary {
				_, _ = output.WriteRune(separator)
			}
		}

		_, _ = output.WriteRune(convert(letter))
		previous = letter
	}

	return output.String()
}
