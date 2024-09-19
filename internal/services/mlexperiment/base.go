// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mlexperiment

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName           = "ML Experiment"
	ItemTFName         = "ml_experiment"
	ItemsName          = "ML Experiments"
	ItemsTFName        = "ml_experiments"
	ItemType           = fabcore.ItemTypeMLExperiment
	ItemDocsSPNSupport = common.DocsSPNNotSupported
	ItemDocsURL        = "https://learn.microsoft.com/fabric/data-science/machine-learning-experiment"
)
