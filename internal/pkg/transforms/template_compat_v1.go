// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package transforms

import (
	"fmt"
	"reflect"
	"slices"
	"text/template"

	sproutmaps "github.com/go-sprout/sprout/registry/maps"
	sproutregex "github.com/go-sprout/sprout/registry/regex"
	sproutslices "github.com/go-sprout/sprout/registry/slices"
	sproutstrings "github.com/go-sprout/sprout/registry/strings"
)

// applyTemplateCompatV1 preserves selected Sprout v1.0.3 call forms and results,
// not full Sprig compatibility. Remove this call and the compatibility files in provider 2.0.
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

	funcs["get"] = templateCompatV1MapFunction2("get", mapRegistry.Get)
	funcs["hasKey"] = templateCompatV1MapFunction2("hasKey", mapRegistry.HasKey)
	funcs["unset"] = templateCompatV1MapFunction2("unset", mapRegistry.Unset)
	funcs["set"] = templateCompatV1MapFunction3("set", mapRegistry.Set)
	// These variadic methods already validate their arguments after reordering.
	funcs["pick"] = func(args ...any) (map[string]any, error) {
		mapLastArgs := templateCompatV1MapArgs(args)

		return mapRegistry.Pick(mapLastArgs...)
	}
	funcs["omit"] = func(args ...any) (map[string]any, error) {
		mapLastArgs := templateCompatV1MapArgs(args)

		return mapRegistry.Omit(mapLastArgs...)
	}
	funcs["append"] = templateCompatV1ListFunction(sliceRegistry.Append)
	funcs["prepend"] = templateCompatV1ListFunction(sliceRegistry.Prepend)
	funcs["slice"] = func(args ...any) (any, error) {
		listLastArgs := templateCompatV1ListArgs(args)

		result, err := sliceRegistry.Slice(listLastArgs...)
		if err == nil && result == nil && len(args) == 1 {
			list := reflect.ValueOf(args[0])
			if (list.Kind() == reflect.Slice || list.Kind() == reflect.Array) && list.Len() == 0 {
				// Preserve the original empty list without bypassing Sprout's validation.
				return args[0], nil
			}
		}

		return result, err
	}
	funcs["without"] = func(args ...any) ([]any, error) {
		listLastArgs := templateCompatV1ListArgs(args)

		return sliceRegistry.Without(listLastArgs...)
	}

	// Preserve the legacy pattern, input, limit order in provider 1.x.
	funcs["regexFindAll"] = func(pattern, value string, limit int) ([]string, error) {
		return regexRegistry.RegexFindAll(pattern, limit, value)
	}
	funcs["regexSplit"] = func(pattern, value string, limit int) ([]string, error) {
		return regexRegistry.RegexSplit(pattern, limit, value)
	}
	// Replacement arguments are both strings, so their two orders are ambiguous.
	funcs["regexReplaceAll"] = func(pattern, value, replacement string) (string, error) {
		return regexRegistry.RegexReplaceAll(pattern, replacement, value)
	}
	funcs["regexReplaceAllLiteral"] = func(pattern, value, replacement string) (string, error) {
		return regexRegistry.RegexReplaceAllLiteral(pattern, replacement, value)
	}

	funcs["substr"] = func(start, end int, value string) string {
		// In v1.0.3 a resolved end index of zero selects an empty substring.
		if end == 0 || end == -len(value) {
			return ""
		}

		return stringRegistry.Substring(start, end, value)
	}

	funcs["toCamelCase"] = templateCompatV1CamelCase
	funcs["toPascalCase"] = templateCompatV1PascalCase

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
	remainingArgs := slices.Clone(args[1:])

	return append(remainingArgs, args[0])
}

// templateCompatV1MapFunction2 wraps fixed-arity map functions with two operands:
// They share the typed (key, dict) signature; T accommodates their different
// results. The returned function validates template arguments before delegation.
// The template engine enforces the fixed argument count.
func templateCompatV1MapFunction2[T any](name string, call func(string, map[string]any) (T, error)) func(any, any) (T, error) {
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

// templateCompatV1MapFunction3 wraps map functions with a key, value, and map.
// It supports legacy map-first calls and delegates normalized arguments.
func templateCompatV1MapFunction3[T any](name string, call func(string, any, map[string]any) (T, error)) func(any, any, any) (T, error) {
	return func(keyArg, value, mapArg any) (T, error) {
		if _, mapFirst := keyArg.(map[string]any); mapFirst {
			keyArg, value, mapArg = value, mapArg, keyArg
		}

		key, keyOK := keyArg.(string)

		dict, dictOK := mapArg.(map[string]any)
		if !keyOK || !dictOK || dict == nil {
			var zero T

			return zero, fmt.Errorf("%s requires a string key and a non-nil map[string]any", name)
		}

		return call(key, value, dict)
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
