// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package transforms_test

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// Golden results were captured with the provider's Sprout v1.0.3 registries.
// These cases exclude process-crashing inputs, which belong in safety tests.
func TestUnit_TemplateCompatV1Edges(t *testing.T) {
	t.Parallel()

	var cases []struct {
		Bucket   string `json:"bucket"`
		Name     string `json:"name"`
		Template string `json:"template"`
		Want     string `json:"want"`
		Error    bool   `json:"error"`
	}

	data, err := os.ReadFile(testhelp.GetFixturesDirPath("template_compat_v1_edges.json"))
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &cases))
	require.NotEmpty(t, cases)

	for _, test := range cases {
		t.Run(test.Bucket+"/"+test.Name, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join(t.TempDir(), "template.txt")
			require.NoError(t, os.WriteFile(path, []byte(test.Template), 0o600))

			payload, hash, diags := transforms.SourceFileToPayload(path, transforms.ProcessingModeGoTemplate, nil, nil, transforms.TokensDelimiterCurlyBraces)
			if test.Error {
				require.True(t, diags.HasError())
				require.Equal(t, common.ErrorTmplParseHeader, diags.Errors()[0].Summary())
				require.Empty(t, payload)
				require.Empty(t, hash)

				return
			}

			require.False(t, diags.HasError(), "%v", diags)
			require.Equal(t, base64.StdEncoding.EncodeToString([]byte(test.Want)), payload)
			require.Equal(t, fmt.Sprintf("%x", sha256.Sum256([]byte(test.Want))), hash)
		})
	}
}
