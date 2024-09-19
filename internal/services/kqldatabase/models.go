// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package kqldatabase

import (
	"context"
	"fmt"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	timeoutsd "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsr "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabkqldatabase "github.com/microsoft/fabric-sdk-go/fabric/kqldatabase"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceKQLDatabaseModel struct {
	baseKQLDatabasePropertiesModel
	Timeouts timeoutsd.Value `tfsdk:"timeouts"`
}

type resourceKQLDatabaseModel struct {
	baseKQLDatabasePropertiesModel
	Configuration supertypes.SingleNestedObjectValueOf[kqlDatabaseConfigurationModel] `tfsdk:"configuration"`
	Timeouts      timeoutsr.Value                                                     `tfsdk:"timeouts"`
}

type baseKQLDatabaseModel struct {
	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
	Description types.String     `tfsdk:"description"`
}

func (to *baseKQLDatabaseModel) set(from fabkqldatabase.KQLDatabase) {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
}

type baseKQLDatabasePropertiesModel struct {
	baseKQLDatabaseModel
	Properties supertypes.SingleNestedObjectValueOf[kqlDatabasePropertiesModel] `tfsdk:"properties"`
}

func (to *baseKQLDatabasePropertiesModel) setProperties(ctx context.Context, from fabkqldatabase.KQLDatabase) diag.Diagnostics {
	properties := supertypes.NewSingleNestedObjectValueOfNull[kqlDatabasePropertiesModel](ctx)

	if from.Properties != nil {
		propertiesModel := &kqlDatabasePropertiesModel{}
		propertiesModel.set(from.Properties)

		if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
			return diags
		}
	}

	to.Properties = properties

	return nil
}

type kqlDatabasePropertiesModel struct {
	DatabaseType        types.String     `tfsdk:"database_type"`
	EventhouseID        customtypes.UUID `tfsdk:"eventhouse_id"`
	IngestionServiceURI customtypes.URL  `tfsdk:"ingestion_service_uri"`
	QueryServiceURI     customtypes.URL  `tfsdk:"query_service_uri"`
	// OneLakeStandardStoragePeriod types.String     `tfsdk:"onelake_standard_storage_period"`
	// OneLakeCachingPeriod         types.String     `tfsdk:"onelake_caching_period"`
}

func (to *kqlDatabasePropertiesModel) set(from *fabkqldatabase.Properties) {
	to.DatabaseType = types.StringPointerValue((*string)(from.DatabaseType))
	to.EventhouseID = customtypes.NewUUIDPointerValue(from.ParentEventhouseItemID)
	to.IngestionServiceURI = customtypes.NewURLPointerValue(from.IngestionServiceURI)
	to.QueryServiceURI = customtypes.NewURLPointerValue(from.QueryServiceURI)
	// to.OneLakeStandardStoragePeriod = types.StringPointerValue(from.OneLakeStandardStoragePeriod)
	// to.OneLakeCachingPeriod = types.StringPointerValue(from.OneLakeCachingPeriod)
}

type kqlDatabaseConfigurationModel struct {
	DatabaseType       types.String     `tfsdk:"database_type"`
	EventhouseID       customtypes.UUID `tfsdk:"eventhouse_id"`
	InvitationToken    types.String     `tfsdk:"invitation_token"`
	SourceClusterURI   customtypes.URL  `tfsdk:"source_cluster_uri"`
	SourceDatabaseName types.String     `tfsdk:"source_database_name"`
}

type requestUpdateKQLDatabase struct {
	fabkqldatabase.UpdateKQLDatabaseRequest
}

func (to *requestUpdateKQLDatabase) set(from resourceKQLDatabaseModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
}

type requestCreateKQLDatabase struct {
	fabkqldatabase.CreateKQLDatabaseRequest
}

func (to *requestCreateKQLDatabase) set(ctx context.Context, from resourceKQLDatabaseModel) diag.Diagnostics {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()

	configuration, diags := from.Configuration.Get(ctx)
	if diags.HasError() {
		return diags
	}

	kqlDatabaseType := (fabkqldatabase.Type)(configuration.DatabaseType.ValueString())

	switch kqlDatabaseType {
	case fabkqldatabase.TypeReadWrite:
		to.CreationPayload = &fabkqldatabase.ReadWriteDatabaseCreationPayload{
			DatabaseType:           &kqlDatabaseType,
			ParentEventhouseItemID: configuration.EventhouseID.ValueStringPointer(),
		}
	case fabkqldatabase.TypeShortcut:
		creationPayload := fabkqldatabase.ShortcutDatabaseCreationPayload{}
		creationPayload.DatabaseType = &kqlDatabaseType
		creationPayload.ParentEventhouseItemID = configuration.EventhouseID.ValueStringPointer()

		if !configuration.InvitationToken.IsNull() && !configuration.InvitationToken.IsUnknown() {
			creationPayload.InvitationToken = configuration.InvitationToken.ValueStringPointer()

			to.CreationPayload = &creationPayload

			return nil
		}

		if !configuration.SourceClusterURI.IsNull() && !configuration.SourceClusterURI.IsUnknown() {
			creationPayload.SourceClusterURI = configuration.SourceClusterURI.ValueStringPointer()
		}

		creationPayload.SourceDatabaseName = configuration.SourceDatabaseName.ValueStringPointer()

		to.CreationPayload = &creationPayload

		return nil
	default:
		diags.AddError(
			"Unsupported KQL database type",
			fmt.Sprintf("The KQL database type '%s' is not supported.", string(kqlDatabaseType)),
		)

		return diags
	}

	return nil
}
