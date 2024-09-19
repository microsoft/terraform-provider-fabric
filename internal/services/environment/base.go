// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package environment

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName           = "Environment"
	ItemTFName         = "environment"
	ItemsName          = "Environments"
	ItemsTFName        = "environments"
	ItemType           = fabcore.ItemTypeEnvironment
	ItemDocsSPNSupport = common.DocsSPNSupported
	ItemDocsURL        = "https://learn.microsoft.com/fabric/data-engineering/create-and-use-environment"
)
