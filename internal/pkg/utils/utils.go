// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"
)

// RemoveSliceByValue removes the first occurrence of the specified value from the slice.
func RemoveSliceByValue[T comparable](slice []T, value T) []T {
	index := -1

	for i, v := range slice {
		if v == value {
			index = i

			break
		}
	}

	if index == -1 {
		return slice // Return the original slice if the value is not found
	}

	// Make a copy of the slice
	sliceCopy := append([]T(nil), slice...)

	return slices.Delete(sliceCopy, index, index+1)
}

func RemoveSlicesByValues[T comparable](slice, value []T) []T {
	for _, v := range value {
		slice = RemoveSliceByValue(slice, v)
	}

	return slice
}

func ConvertEnumsToStringSlices[T any](values []T, sorting bool) []string { //revive:disable-line:flag-parameter
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = fmt.Sprintf("%v", v)
	}

	if sorting {
		slices.Sort(result)
	}

	return result
}

func ConvertStringSlicesToString[T any](values []T, backticks, sorting bool, separator ...string) string { //revive:disable-line:flag-parameter
	result := ConvertEnumsToStringSlices(values, sorting)

	if backticks {
		// Add backticks to each string
		for i, value := range result {
			result[i] = fmt.Sprintf("`%s`", value)
		}
	}

	var sepValue string
	if len(separator) == 0 {
		sepValue = ", " // default
	} else {
		sepValue = separator[0]
	}

	return strings.Join(result, sepValue)
}

// SortMapStringByKeys sorts a map[string]string by keys.
func SortMapStringByKeys[T any](m map[string]T) map[string]T {
	sortedKeys := make([]string, 0, len(m))
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}

	slices.Sort(sortedKeys)

	sortedMap := make(map[string]T)
	for _, k := range sortedKeys {
		sortedMap[k] = m[k]
	}

	return sortedMap
}

func Sha256(content string) string {
	hash := sha256.Sum256([]byte(content))

	return hex.EncodeToString(hash[:])
}
