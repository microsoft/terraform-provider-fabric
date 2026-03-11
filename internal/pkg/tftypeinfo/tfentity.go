// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package tftypeinfo

import (
	"fmt"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

type TFTypeInfo struct {
	ProviderTypeName string
	Name             string
	Type             string
	Names            string
	Types            string
	DocsURL          string
	IsPreview        bool
	IsSPNSupported   bool
}

func (t TFTypeInfo) FullTypeName(plural bool) string { //revive:disable-line:flag-parameter
	typeName := t.Type

	if plural {
		typeName = t.Types
	}

	providerTypeName := common.ProviderTypeName

	if t.ProviderTypeName != "" {
		providerTypeName = t.ProviderTypeName
	}

	return fmt.Sprintf("%s_%s", providerTypeName, typeName)
}
