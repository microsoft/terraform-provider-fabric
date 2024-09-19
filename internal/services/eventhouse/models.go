// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventhouse

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	timeoutsd "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsr "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabeventhouse "github.com/microsoft/fabric-sdk-go/fabric/eventhouse"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceEventhouseModel struct {
	baseEventhousePropertiesModel
	Timeouts timeoutsd.Value `tfsdk:"timeouts"`
}

type resourceEventhouseModel struct {
	baseEventhousePropertiesModel
	Timeouts timeoutsr.Value `tfsdk:"timeouts"`
}

type baseEventhouseModel struct {
	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
	Description types.String     `tfsdk:"description"`
}

func (to *baseEventhouseModel) set(from fabeventhouse.Eventhouse) {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
}

type baseEventhousePropertiesModel struct {
	baseEventhouseModel
	Properties supertypes.SingleNestedObjectValueOf[eventhousePropertiesModel] `tfsdk:"properties"`
}

func (to *baseEventhousePropertiesModel) setProperties(ctx context.Context, from fabeventhouse.Eventhouse) diag.Diagnostics {
	properties := supertypes.NewSingleNestedObjectValueOfNull[eventhousePropertiesModel](ctx)

	if from.Properties != nil {
		propertiesModel := &eventhousePropertiesModel{}
		propertiesModel.set(ctx, from.Properties)

		if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
			return diags
		}
	}

	to.Properties = properties

	return nil
}

type eventhousePropertiesModel struct {
	IngestionServiceURI types.String                   `tfsdk:"ingestion_service_uri"`
	QueryServiceURI     types.String                   `tfsdk:"query_service_uri"`
	DatabaseIDs         supertypes.ListValueOf[string] `tfsdk:"database_ids"`
}

func (to *eventhousePropertiesModel) set(ctx context.Context, from *fabeventhouse.Properties) {
	to.IngestionServiceURI = types.StringPointerValue(from.IngestionServiceURI)
	to.QueryServiceURI = types.StringPointerValue(from.QueryServiceURI)
	to.DatabaseIDs = supertypes.NewListValueOfSlice(ctx, from.DatabasesItemIDs)
}

type requestUpdateEventhouse struct {
	fabeventhouse.UpdateEventhouseRequest
}

func (to *requestUpdateEventhouse) set(from resourceEventhouseModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
}

type requestCreateEventhouse struct {
	fabeventhouse.CreateEventhouseRequest
}

func (to *requestCreateEventhouse) set(from resourceEventhouseModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
}
