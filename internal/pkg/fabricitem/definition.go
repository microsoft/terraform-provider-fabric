// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	superstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

type DefinitionFormat struct {
	Type  string
	API   string
	Paths []string
}

func getDefinitionFormats(values []DefinitionFormat) []string {
	results := make([]string, len(values))

	for i, value := range values {
		results[i] = value.Type
	}

	return results
}

func getDefinitionFormatsPaths(values []DefinitionFormat) map[string][]string {
	results := make(map[string][]string)

	for _, v := range values {
		results[v.Type] = v.Paths
	}

	return results
}

func getDefinitionFormatsPathsDocs(values []DefinitionFormat) string {
	elements := getDefinitionFormatsPaths(values)

	var results string

	i := 0

	for k, v := range elements {
		results += "**" + k + "** format: "
		results += utils.ConvertStringSlicesToString(v, true, true)

		if i != len(elements)-1 {
			results += " "
		}

		i++
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

func getDefinitionFormatAPI(values []DefinitionFormat, format string) string {
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
