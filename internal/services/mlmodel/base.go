// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mlmodel

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName           = "ML Model"
	ItemTFName         = "ml_model"
	ItemsName          = "ML Models"
	ItemsTFName        = "ml_models"
	ItemType           = fabcore.ItemTypeMLModel
	ItemDocsSPNSupport = common.DocsSPNNotSupported
	ItemDocsURL        = "https://learn.microsoft.com/fabric/data-science/machine-learning-model"
)
