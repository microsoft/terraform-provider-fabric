// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package digitaltwinbuilder_test

import (
	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")
