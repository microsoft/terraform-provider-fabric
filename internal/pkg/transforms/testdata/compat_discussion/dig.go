// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

// Package compatdiscussion contains unregistered compatibility alternatives.
package compatdiscussion

import (
	"errors"
	"fmt"
	"strings"
)

// LegacyDig is a discussion-only implementation of Sprout v1.0.3 dig.
func LegacyDig(args ...any) (any, error) {
	if len(args) < 2 {
		return nil, errors.New("dig requires at least two arguments: a sequence of keys and a dictionary")
	}

	dict, ok := args[len(args)-1].(map[string]any)
	if !ok {
		return nil, errors.New("last argument must be a map[string]any")
	}

	var keys []string

	for i, arg := range args[:len(args)-1] {
		key, keyOK := arg.(string)
		if !keyOK {
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
