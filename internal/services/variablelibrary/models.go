// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package variablelibrary

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabvariablelibrary "github.com/microsoft/fabric-sdk-go/fabric/variablelibrary"
)

type variablelibraryPropertiesModel struct {
	ActiveValueSetName types.String `tfsdk:"active_value_set_name"`
}

func (to *variablelibraryPropertiesModel) set(ctx context.Context, from *fabvariablelibrary.Properties) diag.Diagnostics {
	if from == nil {
		return nil
	}

	to.ActiveValueSetName = types.StringPointerValue(from.ActiveValueSetName)

	return nil
}
