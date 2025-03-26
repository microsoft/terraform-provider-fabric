// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sparkcustompool

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Spark Custom Pool",
	Type:           "spark_custom_pool",
	DocsURL:        "https://learn.microsoft.com/fabric/data-engineering/create-custom-spark-pools",
	IsPreview:      true,
	IsSPNSupported: true,
}
