// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	superstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DefinitionFormat struct {
	Type  string
	API   string
	Paths []string
}

func GetDefinitionFormats(values []DefinitionFormat) []string {
	results := make([]string, len(values))

	for i, value := range values {
		results[i] = value.Type
	}

	return results
}

func GetDefinitionFormatsPaths(values []DefinitionFormat) []string {
	var results []string

	for _, value := range values {
		results = append(results, value.Paths...)
	}

	return results
}

func GetDefinitionFormatPaths(values []DefinitionFormat, format string) []string {
	for _, value := range values {
		if value.Type == format {
			return value.Paths
		}
	}

	return nil
}

func GetDefinitionFormatAPI(values []DefinitionFormat, format string) string {
	for _, value := range values {
		if value.Type == format {
			return value.API
		}
	}

	return ""
}

func DefinitionPathKeysValidator(values []DefinitionFormat) []validator.String {
	results := make([]validator.String, 0, len(values))

	for _, value := range values {
		paths := []superstringvalidator.OneOfWithDescriptionIfAttributeIsOneOfValues{}

		for _, p := range value.Paths {
			paths = append(paths, superstringvalidator.OneOfWithDescriptionIfAttributeIsOneOfValues{
				Value:       p,
				Description: p,
			})
		}

		stringValidator := superstringvalidator.OneOfWithDescriptionIfAttributeIsOneOf(
			path.MatchRoot("format"),
			[]attr.Value{types.StringValue(value.Type)},
			paths...,
		)

		results = append(results, stringValidator)
	}

	return results
}
