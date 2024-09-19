// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/microsoft/terraform-provider-fabric/internal/auth"
)

func GetValueOrFileValue(attValue, attFile string, value, file types.String) (string, error) {
	valueResult := value.ValueString()

	if path := file.ValueString(); path != "" {
		fileRaw, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("reading '%s' from file %q: %w", attFile, path, err)
		}

		fileResult := strings.TrimSpace(string(fileRaw))
		if valueResult != "" && valueResult != fileResult {
			return "", fmt.Errorf("mismatch between supplied '%s' and supplied '%s' file content - please either remove one or ensure they match", attValue, attFile)
		}

		valueResult = fileResult
	}

	return valueResult, nil
}

func GetCertOrFileCert(attValue, attFile string, value, file types.String) (string, error) {
	valueResult := strings.TrimSpace(value.ValueString())

	if path := file.ValueString(); path != "" {
		b64, err := auth.ConvertFileToBase64(path)
		if err != nil {
			return "", fmt.Errorf("reading '%s' from file %q: %w", attFile, path, err)
		}

		fileResult := strings.TrimSpace(b64)
		if valueResult != "" && valueResult != fileResult {
			return "", fmt.Errorf("mismatch between supplied '%s' and supplied '%s' file content - please either remove one or ensure they match", attValue, attFile)
		}

		valueResult = fileResult
	}

	return valueResult, nil
}
