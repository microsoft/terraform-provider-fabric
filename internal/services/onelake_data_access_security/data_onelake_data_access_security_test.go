// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelake_data_access_security_test

import (
	"testing"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_OneLakeDataAccessSecurityDataSource(t *testing.T) {
}
