// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetListStringValues(value types.Set, envVarKeys, defaultValue []string) types.Set {
	if value.IsUnknown() || value.IsNull() {
		values := []attr.Value{}

		for _, v := range getEnvList(envVarKeys) {
			values = append(values, types.StringValue(strings.TrimSpace(v)))
		}

		if len(values) == 0 {
			for _, v := range defaultValue {
				values = append(values, types.StringValue(strings.TrimSpace(v)))
			}
		}

		return types.SetValueMust(types.StringType, values)
	}

	return value
}

func getEnvList(envVarKeys []string) []string {
	if v, ok := getMultiEnvVar(envVarKeys); ok {
		return strings.Split(v, ";")
	}

	return nil
}

func GetStringValue(value types.String, envVarKeys []string, defaultValue string) types.String {
	if value.IsUnknown() || value.IsNull() {
		return types.StringValue(strings.TrimSpace(getEnvString(envVarKeys, defaultValue)))
	}

	return value
}

func getEnvString(envVarKeys []string, defaultValue string) string {
	if v, ok := getMultiEnvVar(envVarKeys); ok {
		return v
	}

	return defaultValue
}

func GetBoolValue(value types.Bool, envVarKeys []string, defaultValue bool) types.Bool {
	if value.IsUnknown() || value.IsNull() {
		return types.BoolValue(getEnvBool(envVarKeys, defaultValue))
	}

	return value
}

func getEnvBool(envVarKeys []string, defaultValue bool) bool {
	if v, ok := getMultiEnvVar(envVarKeys); ok {
		strBool, err := strconv.ParseBool(v)
		if err != nil {
			return defaultValue
		}

		return strBool
	}

	return defaultValue
}

// getMultiEnvVar returns the value of the first environment variable that is set.
func getMultiEnvVar(envVarNames []string) (string, bool) {
	for _, envVarName := range envVarNames {
		if value, ok := os.LookupEnv(envVarName); ok {
			return value, ok
		}
	}

	return "", false
}
