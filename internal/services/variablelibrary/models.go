// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package variablelibrary

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabvariablelibrary "github.com/microsoft/fabric-sdk-go/fabric/variablelibrary"
)

type variableLibraryPropertiesModel struct {
	ActiveValueSetName types.String `tfsdk:"active_value_set_name"`
}

func (to *variableLibraryPropertiesModel) set(from fabvariablelibrary.Properties) {
	to.ActiveValueSetName = types.StringPointerValue((*string)(from.ActiveValueSetName))
}
