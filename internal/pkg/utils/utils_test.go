// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func TestUnit_RemoveSliceByValue(t *testing.T) {
	t.Parallel()

	testCasesInt := map[string]struct {
		slice    []int
		value    int
		expected []int
	}{
		"int - remove one": {
			slice:    []int{1, 2, 3, 4, 5},
			value:    3,
			expected: []int{1, 2, 4, 5},
		},
		"int - remove none": {
			slice:    []int{1, 2, 3, 4, 5},
			value:    6,
			expected: []int{1, 2, 3, 4, 5},
		},
		"int - empty": {
			slice:    []int{},
			value:    3,
			expected: []int{},
		},
	}

	for name, testCase := range testCasesInt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := utils.RemoveSliceByValue(testCase.slice, testCase.value)
			assert.Equal(t, testCase.expected, result, "they should be equal")
		})
	}

	testCasesString := map[string]struct {
		slice    []string
		value    string
		expected []string
	}{
		"string - remove one": {
			slice:    []string{"a", "b", "c", "d", "e"},
			value:    "c",
			expected: []string{"a", "b", "d", "e"},
		},
		"string - remove none": {
			slice:    []string{"a", "b", "c", "d", "e"},
			value:    "f",
			expected: []string{"a", "b", "c", "d", "e"},
		},
		"string - empty": {
			slice:    []string{},
			value:    "c",
			expected: []string{},
		},
	}

	for name, testCase := range testCasesString {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := utils.RemoveSliceByValue(testCase.slice, testCase.value)
			assert.Equal(t, testCase.expected, result, "they should be equal")
		})
	}
}

func TestUnit_RemoveSlicesByValues(t *testing.T) {
	t.Parallel()

	testCasesInt := map[string]struct {
		slice    []int
		values   []int
		expected []int
	}{
		"int - empty": {
			slice:    []int{},
			values:   []int{},
			expected: []int{},
		},
		"int - empty values": {
			slice:    []int{1, 2, 3, 4, 5},
			values:   []int{},
			expected: []int{1, 2, 3, 4, 5},
		},
		"int - empty slice": {
			slice:    []int{},
			values:   []int{1, 2, 3, 4, 5},
			expected: []int{},
		},
		"int - remove one": {
			slice:    []int{1, 2, 3, 4, 5},
			values:   []int{3},
			expected: []int{1, 2, 4, 5},
		},
		"int - remove multiple": {
			slice:    []int{1, 2, 3, 4, 5},
			values:   []int{2, 4},
			expected: []int{1, 3, 5},
		},
		"int - remove all": {
			slice:    []int{1, 2, 3, 4, 5},
			values:   []int{1, 2, 3, 4, 5},
			expected: []int{},
		},
		"int - remove none": {
			slice:    []int{1, 2, 3, 4, 5},
			values:   []int{6, 7},
			expected: []int{1, 2, 3, 4, 5},
		},
	}

	for name, testCase := range testCasesInt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := utils.RemoveSlicesByValues(testCase.slice, testCase.values)
			assert.Equal(t, testCase.expected, result, "they should be equal")
		})
	}

	testCasesString := map[string]struct {
		slice    []string
		values   []string
		expected []string
	}{
		"string - empty": {
			slice:    []string{},
			values:   []string{},
			expected: []string{},
		},
		"string - empty values": {
			slice:    []string{"a", "b", "c", "d", "e"},
			values:   []string{},
			expected: []string{"a", "b", "c", "d", "e"},
		},
		"string - empty slice": {
			slice:    []string{},
			values:   []string{"a", "b", "c", "d", "e"},
			expected: []string{},
		},
		"string - remove one": {
			slice:    []string{"a", "b", "c", "d", "e"},
			values:   []string{"c"},
			expected: []string{"a", "b", "d", "e"},
		},
		"string - remove multiple": {
			slice:    []string{"a", "b", "c", "d", "e"},
			values:   []string{"b", "d"},
			expected: []string{"a", "c", "e"},
		},
		"string - remove all": {
			slice:    []string{"a", "b", "c", "d", "e"},
			values:   []string{"a", "b", "c", "d", "e"},
			expected: []string{},
		},
		"string - remove none": {
			slice:    []string{"a", "b", "c", "d", "e"},
			values:   []string{"f", "g"},
			expected: []string{"a", "b", "c", "d", "e"},
		},
	}

	for name, testCase := range testCasesString {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := utils.RemoveSlicesByValues(testCase.slice, testCase.values)
			assert.Equal(t, testCase.expected, result, "they should be equal")
		})
	}
}

func TestUnit_ConvertEnumsToStringSlices(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		values   []int
		expected []string
		sorting  bool
	}{
		"unsorted": {
			values:   []int{1, 5, 4, 3, 2},
			expected: []string{"1", "5", "4", "3", "2"},
			sorting:  false,
		},
		"sorted": {
			values:   []int{1, 5, 4, 3, 2},
			expected: []string{"1", "2", "3", "4", "5"},
			sorting:  true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := utils.ConvertEnumsToStringSlices(testCase.values, testCase.sorting)
			assert.Equal(t, testCase.expected, result, "they should be equal")
		})
	}
}

func TestUnit_ConvertStringSlicesToString(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input     []string
		expected  string
		backticks bool
		sorting   bool
		separator string
	}{
		"empty": {
			input:    []string{},
			expected: "",
		},
		"single": {
			input:     []string{"Hello"},
			expected:  "Hello",
			backticks: false,
			separator: ", ",
		},
		"multiple": {
			input:     []string{"Hello", "World", "Fabric"},
			expected:  "Hello, World, Fabric",
			backticks: false,
			sorting:   false,
			separator: ", ",
		},
		"multiple with backticks": {
			input:     []string{"Hello", "World", "Fabric"},
			expected:  "`Hello`, `World`, `Fabric`",
			backticks: true,
			sorting:   false,
			separator: ", ",
		},
		"multiple with sorting": {
			input:     []string{"Hello", "World", "Fabric"},
			expected:  "Fabric, Hello, World",
			backticks: false,
			sorting:   true,
			separator: ", ",
		},
		"multiple with backticks and sorting": {
			input:     []string{"Hello", "World", "Fabric"},
			expected:  "`Fabric`, `Hello`, `World`",
			backticks: true,
			sorting:   true,
			separator: ", ",
		},
		"multiple with custom separator": {
			input:     []string{"Hello", "World", "Fabric"},
			expected:  "Hello|World|Fabric",
			backticks: false,
			sorting:   false,
			separator: "|",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := utils.ConvertStringSlicesToString(testCase.input, testCase.backticks, testCase.sorting, testCase.separator)
			assert.Equal(t, testCase.expected, result, "they should be equal")
		})
	}
}

func TestUnit_SortMapStringByKeys(t *testing.T) {
	t.Parallel()

	input := map[string]string{
		"b": "value2",
		"a": "value1",
		"c": "value3",
	}

	expected := map[string]string{
		"a": "value1",
		"b": "value2",
		"c": "value3",
	}

	testCases := map[string]struct {
		input    map[string]string
		expected map[string]string
	}{
		"unsorted": {
			input:    input,
			expected: expected,
		},
		"sorted": {
			input:    expected,
			expected: expected,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := utils.SortMapStringByKeys(testCase.input)
			assert.Equal(t, testCase.expected, result, "they should be equal")
		})
	}
}

func TestUnit_Sha256(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input    string
		expected string
	}{
		"empty": {
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		"hello world": {
			input:    "Hello, World!",
			expected: "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := utils.Sha256(testCase.input)
			assert.Equal(t, testCase.expected, result, "they should be equal")
		})
	}
}
