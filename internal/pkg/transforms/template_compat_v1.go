// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package transforms

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"unicode/utf8"

	sproutmaps "github.com/go-sprout/sprout/registry/maps"
	sproutregex "github.com/go-sprout/sprout/registry/regex"
	sproutslices "github.com/go-sprout/sprout/registry/slices"
	sproutstrings "github.com/go-sprout/sprout/registry/strings"
)

// applyTemplateCompatV1 preserves the Sprout v1.0.3 contract, not Sprig's
// different behavior. Remove this call and the compatibility files in provider 2.0.
func applyTemplateCompatV1(funcs template.FuncMap) {
	maps := sproutmaps.NewRegistry()
	lists := sproutslices.NewRegistry()
	regex := sproutregex.NewRegistry()
	text := sproutstrings.NewRegistry()

	funcs["get"] = templateCompatV1MapFunction("get", maps.Get)
	funcs["hasKey"] = templateCompatV1MapFunction("hasKey", maps.HasKey)
	funcs["unset"] = templateCompatV1MapFunction("unset", maps.Unset)
	funcs["set"] = func(args ...any) (map[string]any, error) {
		args = templateCompatV1MapArgs(args)
		if len(args) != 3 {
			return nil, errors.New("set requires exactly three arguments")
		}

		key, keyOK := args[0].(string)

		dict, dictOK := args[2].(map[string]any)
		if !keyOK || !dictOK || dict == nil {
			return nil, errors.New("set requires a string key and a non-nil map[string]any")
		}

		return maps.Set(key, args[1], dict)
	}
	funcs["pick"] = func(args ...any) (map[string]any, error) {
		return maps.Pick(templateCompatV1MapArgs(args)...)
	}
	funcs["omit"] = func(args ...any) (map[string]any, error) {
		return maps.Omit(templateCompatV1MapArgs(args)...)
	}
	funcs["dig"] = templateCompatV1Dig

	funcs["append"] = templateCompatV1ListFunction("append", lists.Append)
	funcs["prepend"] = templateCompatV1ListFunction("prepend", lists.Prepend)
	funcs["slice"] = func(args ...any) (any, error) {
		if len(args) == 1 {
			return args[0], nil
		}

		return lists.Slice(templateCompatV1ListArgs(args)...)
	}
	funcs["without"] = func(args ...any) ([]any, error) {
		if len(args) < 2 {
			return nil, errors.New("without requires at least two arguments")
		}

		return lists.Without(templateCompatV1ListArgs(args)...)
	}

	// These all-string signatures are ambiguous: keep only the old order in 1.x.
	funcs["regexFindAll"] = func(pattern, value string, limit int) ([]string, error) {
		return regex.RegexFindAll(pattern, limit, value)
	}
	funcs["regexSplit"] = func(pattern, value string, limit int) ([]string, error) {
		return regex.RegexSplit(pattern, limit, value)
	}
	funcs["regexReplaceAll"] = func(pattern, value, replacement string) (string, error) {
		return regex.RegexReplaceAll(pattern, replacement, value)
	}
	funcs["regexReplaceAllLiteral"] = func(pattern, value, replacement string) (string, error) {
		return regex.RegexReplaceAllLiteral(pattern, replacement, value)
	}

	funcs["substr"] = func(start, end int, value string) string {
		if end == 0 || end == -len(value) {
			return ""
		}

		return text.Substring(start, end, value)
	}

	funcs["ellipsisBoth"] = func(offset, maxWidth int, value string) string {
		if maxWidth < 4 || offset > 0 && maxWidth < 7 {
			return value
		}

		// Preserve byte offsets in the old non-truncating branch without restoring
		// its out-of-range panic. The truncating branch already used rune offsets.
		runeCount := utf8.RuneCountInString(value)
		if offset >= 0 && offset <= len(value) && (runeCount <= maxWidth || runeCount <= offset) {
			return value[offset:]
		}

		return text.EllipsisBoth(offset, maxWidth, value)
	}
	for _, name := range []string{"toCamelCase", "toPascalCase", "toSnakeCase", "toKebabCase", "toDotCase", "toPathCase", "toConstantCase"} {
		funcs[name] = func(value string) string {
			return templateCompatV1Case(name, value)
		}
	}

	// Build resolves aliases before these overrides; regex also omits the old must* aliases.
	for alias, name := range map[string]string{
		"mustAppend":                 "append",
		"mustPrepend":                "prepend",
		"mustSlice":                  "slice",
		"mustWithout":                "without",
		"mustRegexFind":              "regexFind",
		"mustRegexFindAll":           "regexFindAll",
		"mustRegexMatch":             "regexMatch",
		"mustRegexSplit":             "regexSplit",
		"mustRegexReplaceAll":        "regexReplaceAll",
		"mustRegexReplaceAllLiteral": "regexReplaceAllLiteral",
	} {
		funcs[alias] = funcs[name]
	}
}

func templateCompatV1MapArgs(args []any) []any {
	if len(args) > 1 {
		if _, ok := args[0].(map[string]any); ok {
			return append(append([]any{}, args[1:]...), args[0])
		}
	}

	return args
}

func templateCompatV1MapFunction[T any](name string, call func(string, map[string]any) (T, error)) func(...any) (T, error) {
	return func(args ...any) (T, error) {
		var zero T

		args = templateCompatV1MapArgs(args)
		if len(args) != 2 {
			return zero, fmt.Errorf("%s requires exactly two arguments", name)
		}

		key, keyOK := args[0].(string)

		dict, dictOK := args[1].(map[string]any)
		if !keyOK || !dictOK {
			return zero, fmt.Errorf("%s requires a string key and a map[string]any", name)
		}

		return call(key, dict)
	}
}

func templateCompatV1ListArgs(args []any) []any {
	if len(args) > 1 {
		switch reflect.ValueOf(args[0]).Kind() {
		case reflect.Array, reflect.Slice, reflect.Invalid:
			return append(append([]any{}, args[1:]...), args[0])
		default:
		}
	}

	return args
}

func templateCompatV1ListFunction(name string, call func(any, any) ([]any, error)) func(...any) ([]any, error) {
	return func(args ...any) ([]any, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("%s requires exactly two arguments", name)
		}

		// Reorder once rather than recursively: two lists used to overflow the stack.
		args = templateCompatV1ListArgs(args)

		return call(args[0], args[1])
	}
}

func templateCompatV1Dig(args ...any) (any, error) {
	if len(args) < 2 {
		return nil, errors.New("dig requires at least two arguments: a sequence of keys and a dictionary")
	}

	dict, ok := args[len(args)-1].(map[string]any)
	if !ok {
		return nil, errors.New("last argument must be a map[string]any")
	}

	var keys []string

	for i, arg := range args[:len(args)-1] {
		key, ok := arg.(string)
		if !ok {
			return nil, fmt.Errorf("cannot parse keys: all keys must be strings, got %T at position %d", arg, i)
		}

		keys = append(keys, strings.Split(key, ".")...)
	}

	for i, key := range keys {
		value, exists := dict[key]
		if !exists || i == len(keys)-1 {
			return value, nil
		}

		dict, ok = value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("value at key %q is not a nested dictionary but %T", key, value)
		}
	}

	return nil, errors.New("unexpected termination of key traversal")
}
