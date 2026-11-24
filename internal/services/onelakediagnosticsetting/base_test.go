// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakediagnosticsetting_test

import (
	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/services/onelakediagnosticsetting"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var itemTypeInfo = onelakediagnosticsetting.ItemTypeInfo

var (
	testResourceItemFQN, testResourceItemHeader     = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")
	testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")
)
