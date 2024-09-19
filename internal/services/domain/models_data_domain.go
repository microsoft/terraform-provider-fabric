// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
)

type dataSourceDomainModel struct {
	baseDomainModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
