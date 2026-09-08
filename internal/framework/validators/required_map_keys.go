// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package validators

import (
	"context"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Map = RequiredMapKeysValidator{}

// RequiredMapKeysValidator validates that a map contains all configured keys.
type RequiredMapKeysValidator struct {
	keys []string
}

// RequiredMapKeys returns a validator that requires the given map keys.
func RequiredMapKeys(keys ...string) RequiredMapKeysValidator {
	return RequiredMapKeysValidator{keys: slices.Sorted(slices.Values(keys))}
}

func (v RequiredMapKeysValidator) Description(_ context.Context) string {
	return "map must contain required keys: " + strings.Join(v.keys, ", ")
}

func (v RequiredMapKeysValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v RequiredMapKeysValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elements := req.ConfigValue.Elements()
	missing := make([]string, 0, len(v.keys))

	for _, key := range v.keys {
		if _, ok := elements[key]; !ok {
			missing = append(missing, key)
		}
	}

	if len(missing) == 0 {
		return
	}

	configured := slices.Sorted(maps.Keys(elements))

	resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
		req.Path,
		v.Description(ctx),
		strings.Join(configured, ", "),
	))
}
