// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehouse

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	fabwarehouse "github.com/microsoft/fabric-sdk-go/fabric/warehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getResourceWarehousePropertiesAttributes() map[string]schema.Attribute {
	result := map[string]schema.Attribute{
		"collation_type": schema.StringAttribute{
			MarkdownDescription: "The collation type of the warehouse. Possible values:" + utils.ConvertStringSlicesToString(fabwarehouse.PossibleCollationTypeValues(), true, true) + ".",
			Computed:            true,
		},
		"connection_string": schema.StringAttribute{
			MarkdownDescription: "The SQL connection string connected to the workspace containing this warehouse.",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "The date and time the warehouse was created.",
			Computed:            true,
			CustomType:          timetypes.RFC3339Type{},
		},
		"last_updated_time": schema.StringAttribute{
			MarkdownDescription: "The date and time the warehouse was last updated.",
			Computed:            true,
			CustomType:          timetypes.RFC3339Type{},
		},
	}

	return result
}

func getResourceWarehouseConfigurationAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"collation_type": schema.StringAttribute{
			MarkdownDescription: "The collation type of the warehouse. Accepted values: " + utils.ConvertStringSlicesToString(fabwarehouse.PossibleCollationTypeValues(), true, true) + ".",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabwarehouse.PossibleCollationTypeValues(), false)...),
			},
		},
	}
}
