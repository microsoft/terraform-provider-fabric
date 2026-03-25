// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sqldatabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabsqldatabase "github.com/microsoft/fabric-sdk-go/fabric/sqldatabase"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	superint32validator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/int32validator"
	superobjectvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/objectvalidator"
	superstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getResourceSQLDatabasePropertiesAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"connection_string": schema.StringAttribute{
			MarkdownDescription: "The connection string of the database.",
			Computed:            true,
		},
		"database_name": schema.StringAttribute{
			MarkdownDescription: "The database name.",
			Computed:            true,
		},
		"server_fqdn": schema.StringAttribute{
			MarkdownDescription: "The server fully qualified domain name (FQDN).",
			Computed:            true,
		},
		"backup_retention_days": schema.Int32Attribute{
			MarkdownDescription: "The backup retention period in days.",
			Computed:            true,
		},
		"collation": schema.StringAttribute{
			MarkdownDescription: "The collation of the SQL database.",
			Computed:            true,
		},
		"earliest_restore_point": schema.StringAttribute{
			MarkdownDescription: "The earliest restore point of the database in UTC, using the YYYY-MM-DDTHH:mm:ssZ format.",
			Computed:            true,
			CustomType:          timetypes.RFC3339Type{},
		},
		"latest_restore_point": schema.StringAttribute{
			MarkdownDescription: "The latest restore point of the database in UTC, using the YYYY-MM-DDTHH:mm:ssZ format.",
			Computed:            true,
			CustomType:          timetypes.RFC3339Type{},
		},
	}
}

func getResourceSQLDatabaseConfigurationAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"creation_mode": schema.StringAttribute{
			MarkdownDescription: "The creation mode of the SQL database.",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabsqldatabase.PossibleCreationModeValues(), false)...),
			},
		},
		"backup_retention_days": schema.Int32Attribute{
			MarkdownDescription: "Set the backup retention period in days. The minimum is 1 days. The maximum is 35 days.",
			Optional:            true,
			PlanModifiers: []planmodifier.Int32{
				int32planmodifier.RequiresReplace(),
			},
			Validators: []validator.Int32{
				superint32validator.NullIfAttributeIsOneOf(
					path.MatchRelative().AtParent().AtName("creation_mode"),
					[]attr.Value{types.StringValue(string(fabsqldatabase.CreationModeRestore))},
				),
				int32validator.Between(1, 35),
			},
		},
		"collation": schema.StringAttribute{
			MarkdownDescription: "Set the collation of the SQL database.",
			Optional:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				superstringvalidator.NullIfAttributeIsOneOf(
					path.MatchRelative().AtParent().AtName("creation_mode"),
					[]attr.Value{types.StringValue(string(fabsqldatabase.CreationModeRestore))},
				),
			},
		},
		"restore_point_in_time": schema.StringAttribute{
			MarkdownDescription: "Set the time to restore the source database in UTC, using the YYYY-MM-DDTHH:mm:ssZ format.",
			Optional:            true,
			CustomType:          timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				superstringvalidator.NullIfAttributeIsOneOf(
					path.MatchRelative().AtParent().AtName("creation_mode"),
					[]attr.Value{types.StringValue(string(fabsqldatabase.CreationModeNew))},
				),
				superstringvalidator.RequireIfAttributeIsOneOf(
					path.MatchRelative().AtParent().AtName("creation_mode"),
					[]attr.Value{types.StringValue(string(fabsqldatabase.CreationModeRestore))},
				),
			},
		},
		"source_database_reference": schema.SingleNestedAttribute{
			MarkdownDescription: "Set the reference for the source database to be restored from.",
			Optional:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[sourceDatabaseReferenceModel](ctx),
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.RequiresReplace(),
			},
			Validators: []validator.Object{
				superobjectvalidator.NullIfAttributeIsOneOf(
					path.MatchRelative().AtParent().AtName("creation_mode"),
					[]attr.Value{types.StringValue(string(fabsqldatabase.CreationModeNew))},
				),
				superobjectvalidator.RequireIfAttributeIsOneOf(
					path.MatchRelative().AtParent().AtName("creation_mode"),
					[]attr.Value{types.StringValue(string(fabsqldatabase.CreationModeRestore))},
				),
			},
			Attributes: map[string]schema.Attribute{
				"item_id": schema.StringAttribute{
					MarkdownDescription: "The ID of the item.",
					Optional:            true,
					CustomType:          customtypes.UUIDType{},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						superstringvalidator.RequireIfAttributeIsOneOf(
							path.MatchRelative().AtParent().AtName("reference_type"),
							[]attr.Value{types.StringValue(string(fabsqldatabase.ItemReferenceTypeByID))},
						),
						superstringvalidator.NullIfAttributeIsOneOf(
							path.MatchRelative().AtParent().AtName("reference_type"),
							[]attr.Value{types.StringValue(string(fabsqldatabase.ItemReferenceTypeByVariable))},
						),
					},
				},
				"reference_type": schema.StringAttribute{
					MarkdownDescription: "The item reference type.",
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabsqldatabase.PossibleItemReferenceTypeValues(), false)...),
					},
				},
				"workspace_id": schema.StringAttribute{
					MarkdownDescription: "The workspace ID of the item.",
					Optional:            true,
					CustomType:          customtypes.UUIDType{},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						superstringvalidator.RequireIfAttributeIsOneOf(
							path.MatchRelative().AtParent().AtName("reference_type"),
							[]attr.Value{types.StringValue(string(fabsqldatabase.ItemReferenceTypeByID))},
						),
						superstringvalidator.NullIfAttributeIsOneOf(
							path.MatchRelative().AtParent().AtName("reference_type"),
							[]attr.Value{types.StringValue(string(fabsqldatabase.ItemReferenceTypeByVariable))},
						),
					},
				},
				"variable_reference": schema.StringAttribute{
					MarkdownDescription: "The variable reference. Required when `reference_type` is `ByVariable`.",
					Optional:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						superstringvalidator.RequireIfAttributeIsOneOf(
							path.MatchRelative().AtParent().AtName("reference_type"),
							[]attr.Value{types.StringValue(string(fabsqldatabase.ItemReferenceTypeByVariable))},
						),
						superstringvalidator.NullIfAttributeIsOneOf(
							path.MatchRelative().AtParent().AtName("reference_type"),
							[]attr.Value{types.StringValue(string(fabsqldatabase.ItemReferenceTypeByID))},
						),
					},
				},
			},
		},
	}
}
