// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package mlexperiment

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const FabricItemType = fabcore.ItemTypeMLExperiment

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "ML Experiment",
	Type:           "ml_experiment",
	Names:          "ML Experiments",
	Types:          "ml_experiments",
	DocsURL:        "https://learn.microsoft.com/fabric/data-science/machine-learning-experiment",
	IsPreview:      true,
	IsSPNSupported: false,
}
