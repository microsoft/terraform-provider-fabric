// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package kqldatabase

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabkqldatabase "github.com/microsoft/fabric-sdk-go/fabric/kqldatabase"
	superstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getResourceKQLDatabasePropertiesAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"database_type": schema.StringAttribute{
			MarkdownDescription: "The type of the database. Possible values:" + utils.ConvertStringSlicesToString(fabkqldatabase.PossibleKqlDatabaseTypeValues(), true, true) + ".",
			Computed:            true,
		},
		"eventhouse_id": schema.StringAttribute{
			MarkdownDescription: "Parent Eventhouse ID.",
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"ingestion_service_uri": schema.StringAttribute{
			MarkdownDescription: "Ingestion service URI.",
			Computed:            true,
			CustomType:          customtypes.URLType{},
		},
		"query_service_uri": schema.StringAttribute{
			MarkdownDescription: "Query service URI.",
			Computed:            true,
			CustomType:          customtypes.URLType{},
		},
	}
}

func getResourceKQLDatabaseConfigurationAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"database_type": schema.StringAttribute{
			MarkdownDescription: "The type of the KQL database. Accepted values: " + utils.ConvertStringSlicesToString(fabkqldatabase.PossibleKqlDatabaseTypeValues(), true, true) + ".\n\n" +
				"`" + string(fabkqldatabase.TypeReadWrite) + "` Allows read and write operations on the database.\n\n" +
				"`" + string(fabkqldatabase.TypeShortcut) + "` A shortcut is an embedded reference allowing read only operations on a source database. The source can be in the same or different tenants, either in an Azure Data Explorer cluster or a Fabric Eventhouse.",
			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabkqldatabase.PossibleKqlDatabaseTypeValues(), false)...),
			},
		},
		"eventhouse_id": schema.StringAttribute{
			MarkdownDescription: "Parent Eventhouse ID.",
			Required:            true,
			CustomType:          customtypes.UUIDType{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"invitation_token": schema.StringAttribute{
			MarkdownDescription: "Invitation token to follow the source database. Only allowed when `database_type` is `" + string(fabkqldatabase.TypeShortcut) + "`.",
			Optional:            true,
			Sensitive:           true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.ConflictsWith(
					path.MatchRelative().AtParent().AtName("source_cluster_uri"),
					path.MatchRelative().AtParent().AtName("source_database_name"),
				),
				superstringvalidator.NullIfAttributeIsOneOf(
					path.MatchRelative().AtParent().AtName("database_type"),
					[]attr.Value{types.StringValue(string(fabkqldatabase.TypeReadWrite))},
				),
			},
		},
		"source_cluster_uri": schema.StringAttribute{
			MarkdownDescription: "The URI of the source Eventhouse or Azure Data Explorer cluster. Only allowed when `database_type` is `" + string(fabkqldatabase.TypeShortcut) + "`.",
			Optional:            true,
			CustomType:          customtypes.URLType{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("invitation_token")),
				stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("source_database_name")),
				superstringvalidator.NullIfAttributeIsOneOf(
					path.MatchRelative().AtParent().AtName("database_type"),
					[]attr.Value{types.StringValue(string(fabkqldatabase.TypeReadWrite))},
				),
			},
		},
		"source_database_name": schema.StringAttribute{
			MarkdownDescription: "The name of the database to follow in the source Eventhouse or Azure Data Explorer cluster. Only allowed when `database_type` is `" + string(fabkqldatabase.TypeShortcut) + "`.",
			Optional:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("invitation_token")),
				superstringvalidator.NullIfAttributeIsOneOf(
					path.MatchRelative().AtParent().AtName("database_type"),
					[]attr.Value{types.StringValue(string(fabkqldatabase.TypeReadWrite))},
				),
			},
		},
	}
}
