// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package mlmodel

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const FabricItemType = fabcore.ItemTypeMLModel

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "ML Model",
	Type:           "ml_model",
	Names:          "ML Models",
	Types:          "ml_models",
	DocsURL:        "https://learn.microsoft.com/fabric/data-science/machine-learning-model",
	IsPreview:      true,
	IsSPNSupported: false,
}
