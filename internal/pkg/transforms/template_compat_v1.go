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
//
// Funcs is the name-to-function lookup built by the provider's Sprout handler.
// Replacing its entries installs wrappers; template arguments arrive only when
// a template executes. The registry instances below provide implementation
// methods, not additional registrations or access to unsupported function groups.
func applyTemplateCompatV1(funcs template.FuncMap) {
	mapRegistry := sproutmaps.NewRegistry()
	sliceRegistry := sproutslices.NewRegistry()
	regexRegistry := sproutregex.NewRegistry()
	stringRegistry := sproutstrings.NewRegistry()

	funcs["get"] = templateCompatV1MapFunction("get", mapRegistry.Get)
	funcs["hasKey"] = templateCompatV1MapFunction("hasKey", mapRegistry.HasKey)
	funcs["unset"] = templateCompatV1MapFunction("unset", mapRegistry.Unset)
	funcs["set"] = func(keyArg, value, mapArg any) (map[string]any, error) {
		// Historical calls use (dict, key, value); normalize to (key, value, dict).
		if _, mapFirst := keyArg.(map[string]any); mapFirst {
			keyArg, value, mapArg = value, mapArg, keyArg
		}

		key, keyOK := keyArg.(string)

		dict, dictOK := mapArg.(map[string]any)
		if !keyOK || !dictOK || dict == nil {
			return nil, errors.New("set requires a string key and a non-nil map[string]any")
		}

		return mapRegistry.Set(key, value, dict)
	}
	// These variadic methods already validate their arguments after reordering.
	funcs["pick"] = func(args ...any) (map[string]any, error) {
		mapLastArgs := templateCompatV1MapArgs(args)

		return mapRegistry.Pick(mapLastArgs...)
	}
	funcs["omit"] = func(args ...any) (map[string]any, error) {
		mapLastArgs := templateCompatV1MapArgs(args)

		return mapRegistry.Omit(mapLastArgs...)
	}
	funcs["dig"] = templateCompatV1Dig

	funcs["append"] = templateCompatV1ListFunction(sliceRegistry.Append)
	funcs["prepend"] = templateCompatV1ListFunction(sliceRegistry.Prepend)
	funcs["slice"] = func(args ...any) (any, error) {
		if len(args) == 1 {
			return args[0], nil
		}

		listLastArgs := templateCompatV1ListArgs(args)

		return sliceRegistry.Slice(listLastArgs...)
	}
	funcs["without"] = func(args ...any) ([]any, error) {
		if len(args) < 2 {
			return nil, errors.New("without requires at least two arguments")
		}

		listLastArgs := templateCompatV1ListArgs(args)

		return sliceRegistry.Without(listLastArgs...)
	}

	// These all-string signatures are ambiguous: keep only the old order in 1.x.
	funcs["regexFindAll"] = func(pattern, value string, limit int) ([]string, error) {
		return regexRegistry.RegexFindAll(pattern, limit, value)
	}
	funcs["regexSplit"] = func(pattern, value string, limit int) ([]string, error) {
		return regexRegistry.RegexSplit(pattern, limit, value)
	}
	funcs["regexReplaceAll"] = func(pattern, value, replacement string) (string, error) {
		return regexRegistry.RegexReplaceAll(pattern, replacement, value)
	}
	funcs["regexReplaceAllLiteral"] = func(pattern, value, replacement string) (string, error) {
		return regexRegistry.RegexReplaceAllLiteral(pattern, replacement, value)
	}

	funcs["substr"] = func(start, end int, value string) string {
		if end == 0 || end == -len(value) {
			return ""
		}

		return stringRegistry.Substring(start, end, value)
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

		return stringRegistry.EllipsisBoth(offset, maxWidth, value)
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

// templateCompatV1MapArgs normalizes variadic map-first calls to map-last order:
// [dict, key1, key2] becomes [key1, key2, dict]. Map-last calls are unchanged.
// Argument validation belongs to the caller or the delegated registry method.
func templateCompatV1MapArgs(args []any) []any {
	if len(args) > 1 {
		if _, ok := args[0].(map[string]any); ok {
			return templateCompatV1MoveFirstArgLast(args)
		}
	}

	return args
}

// templateCompatV1MoveFirstArgLast copies the arguments so normalization never
// overwrites the caller's slice. Callers supply at least two arguments.
func templateCompatV1MoveFirstArgLast(args []any) []any {
	reordered := make([]any, len(args))
	copy(reordered, args[1:])
	reordered[len(reordered)-1] = args[0]

	return reordered
}

// templateCompatV1MapFunction creates a wrapper for get, hasKey, and unset.
// They share the typed (key, dict) signature; T accommodates their different
// results. The returned function validates template arguments before delegation.
// The template engine enforces the fixed argument count.
func templateCompatV1MapFunction[T any](name string, call func(string, map[string]any) (T, error)) func(any, any) (T, error) {
	return func(keyArg, mapArg any) (T, error) {
		if _, mapFirst := keyArg.(map[string]any); mapFirst {
			keyArg, mapArg = mapArg, keyArg
		}

		key, keyOK := keyArg.(string)

		dict, dictOK := mapArg.(map[string]any)
		if !keyOK || !dictOK {
			var zero T

			return zero, fmt.Errorf("%s requires a string key and a map[string]any", name)
		}

		return call(key, dict)
	}
}

// templateCompatV1ListArgs preserves v1.0.3's first-argument type detection,
// including nil. A list-like first argument moves to the end exactly once;
// already list-last calls with a scalar first argument are unchanged.
func templateCompatV1ListArgs(args []any) []any {
	if len(args) > 1 {
		switch reflect.ValueOf(args[0]).Kind() {
		case reflect.Array, reflect.Slice, reflect.Invalid:
			return templateCompatV1MoveFirstArgLast(args)
		default:
		}
	}

	return args
}

// templateCompatV1ListFunction wraps append and prepend. Both take (value, list)
// in new Sprout, but provider 1.x also accepts the historical (list, value).
func templateCompatV1ListFunction(call func(any, any) ([]any, error)) func(any, any) ([]any, error) {
	return func(value, list any) ([]any, error) {
		// Reorder once rather than recursively: two lists used to overflow the stack.
		switch reflect.ValueOf(value).Kind() {
		case reflect.Array, reflect.Slice, reflect.Invalid:
			value, list = list, value
		default:
		}

		return call(value, list)
	}
}

// templateCompatV1Dig retains dotted-path splitting, empty segments, and literal
// backslashes. Unlike Sprig's dig, it has no positional default-value argument.
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
