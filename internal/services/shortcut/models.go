// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package shortcut

import (
	"context"
	"fmt"
	"strings"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseShortcutModel struct {
	ID          types.String                                      `tfsdk:"id"`
	Name        types.String                                      `tfsdk:"name"`
	Path        types.String                                      `tfsdk:"path"`
	Target      supertypes.SingleNestedObjectValueOf[targetModel] `tfsdk:"target"`
	WorkspaceID customtypes.UUID                                  `tfsdk:"workspace_id"`
	ItemID      customtypes.UUID                                  `tfsdk:"item_id"`
}

type targetModel struct {
	Type               types.String                                                `tfsdk:"type"`
	Onelake            supertypes.SingleNestedObjectValueOf[oneLakeModel]          `tfsdk:"onelake"`
	AdlsGen2           supertypes.SingleNestedObjectValueOf[targetDataSourceModel] `tfsdk:"adls_gen2"`
	AmazonS3           supertypes.SingleNestedObjectValueOf[targetDataSourceModel] `tfsdk:"amazon_s3"`
	Dataverse          supertypes.SingleNestedObjectValueOf[dataverse]             `tfsdk:"dataverse"`
	GoogleCloudStorage supertypes.SingleNestedObjectValueOf[targetDataSourceModel] `tfsdk:"google_cloud_storage"`
	S3Compatible       supertypes.SingleNestedObjectValueOf[s3Compatible]          `tfsdk:"s3_compatible"`
	AzureBlobStorage   supertypes.SingleNestedObjectValueOf[targetDataSourceModel] `tfsdk:"azure_blob_storage"`
	ExternalDataShare  supertypes.SingleNestedObjectValueOf[externalDataShare]     `tfsdk:"external_data_share"`
}

type targetDataSourceModel struct {
	ConnectionID customtypes.UUID `tfsdk:"connection_id"`
	Location     customtypes.URL  `tfsdk:"location"`
	Subpath      types.String     `tfsdk:"subpath"`
}

type oneLakeModel struct {
	ItemID      customtypes.UUID `tfsdk:"item_id"`
	Path        types.String     `tfsdk:"path"`
	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
}

type dataverse struct {
	ConnectionID      customtypes.UUID `tfsdk:"connection_id"`
	DeltaLakeFolder   types.String     `tfsdk:"deltalake_folder"`
	EnvironmentDomain customtypes.URL  `tfsdk:"environment_domain"`
	TableName         types.String     `tfsdk:"table_name"`
}

type s3Compatible struct {
	ConnectionID customtypes.UUID `tfsdk:"connection_id"`
	Location     customtypes.URL  `tfsdk:"location"`
	Subpath      types.String     `tfsdk:"subpath"`
	Bucket       types.String     `tfsdk:"bucket"`
}
type externalDataShare struct {
	ConnectionID customtypes.UUID `tfsdk:"connection_id"`
}

func (to *baseShortcutModel) set(ctx context.Context, workspaceID, itemID string, from fabcore.Shortcut) diag.Diagnostics {
	to.Name = types.StringPointerValue(from.Name)
	to.Path = types.StringPointerValue(from.Path)
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.ItemID = customtypes.NewUUIDValue(itemID)

	shortcutComputedID := fmt.Sprintf("%s%s%s%s", workspaceID, itemID, strings.TrimPrefix(to.Path.ValueString(), "/"), to.Name.ValueString())

	to.ID = types.StringValue(shortcutComputedID)

	to.Target = supertypes.NewSingleNestedObjectValueOfNull[targetModel](ctx)
	target := supertypes.NewSingleNestedObjectValueOfNull[targetModel](ctx)

	if from.Target != nil {
		targetModel := &targetModel{}
		targetModel.set(ctx, *from.Target)

		if diags := target.Set(ctx, targetModel); diags.HasError() {
			return diags
		}
	}

	to.Target = target

	return nil
}

func (to *targetModel) set(ctx context.Context, from fabcore.Target) {
	to.Type = types.StringPointerValue((*string)(from.Type))
	to.Onelake = supertypes.NewSingleNestedObjectValueOfNull[oneLakeModel](ctx)
	to.AdlsGen2 = supertypes.NewSingleNestedObjectValueOfNull[targetDataSourceModel](ctx)
	to.AmazonS3 = supertypes.NewSingleNestedObjectValueOfNull[targetDataSourceModel](ctx)
	to.Dataverse = supertypes.NewSingleNestedObjectValueOfNull[dataverse](ctx)
	to.GoogleCloudStorage = supertypes.NewSingleNestedObjectValueOfNull[targetDataSourceModel](ctx)
	to.S3Compatible = supertypes.NewSingleNestedObjectValueOfNull[s3Compatible](ctx)
	to.AzureBlobStorage = supertypes.NewSingleNestedObjectValueOfNull[targetDataSourceModel](ctx)
	to.ExternalDataShare = supertypes.NewSingleNestedObjectValueOfNull[externalDataShare](ctx)

	if from.OneLake != nil {
		onelakeModel := &oneLakeModel{
			ItemID:      customtypes.NewUUIDPointerValue(from.OneLake.ItemID),
			Path:        types.StringPointerValue(from.OneLake.Path),
			WorkspaceID: customtypes.NewUUIDPointerValue(from.OneLake.WorkspaceID),
		}
		to.Onelake = supertypes.NewSingleNestedObjectValueOf(ctx, onelakeModel)
	}

	if from.AdlsGen2 != nil {
		adlsGen2Model := &targetDataSourceModel{
			ConnectionID: customtypes.NewUUIDPointerValue(from.AdlsGen2.ConnectionID),
			Location:     customtypes.NewURLPointerValue(from.AdlsGen2.Location),
			Subpath:      types.StringPointerValue(from.AdlsGen2.Subpath),
		}
		to.AdlsGen2 = supertypes.NewSingleNestedObjectValueOf(ctx, adlsGen2Model)
	}

	if from.AmazonS3 != nil {
		amazonS3Model := &targetDataSourceModel{
			ConnectionID: customtypes.NewUUIDPointerValue(from.AmazonS3.ConnectionID),
			Location:     customtypes.NewURLPointerValue(from.AmazonS3.Location),
			Subpath:      types.StringPointerValue(from.AmazonS3.Subpath),
		}
		to.AmazonS3 = supertypes.NewSingleNestedObjectValueOf(ctx, amazonS3Model)
	}

	if from.Dataverse != nil {
		dataverseModel := &dataverse{
			ConnectionID:      customtypes.NewUUIDPointerValue(from.Dataverse.ConnectionID),
			DeltaLakeFolder:   types.StringPointerValue(from.Dataverse.DeltaLakeFolder),
			EnvironmentDomain: customtypes.NewURLPointerValue(from.Dataverse.EnvironmentDomain),
			TableName:         types.StringPointerValue(from.Dataverse.TableName),
		}
		to.Dataverse = supertypes.NewSingleNestedObjectValueOf(ctx, dataverseModel)
	}

	if from.GoogleCloudStorage != nil {
		googleStorageCloudModel := &targetDataSourceModel{
			ConnectionID: customtypes.NewUUIDPointerValue(from.GoogleCloudStorage.ConnectionID),
			Location:     customtypes.NewURLPointerValue(from.GoogleCloudStorage.Location),
			Subpath:      types.StringPointerValue(from.GoogleCloudStorage.Subpath),
		}
		to.GoogleCloudStorage = supertypes.NewSingleNestedObjectValueOf(ctx, googleStorageCloudModel)
	}

	if from.S3Compatible != nil {
		s3CompatibleModel := &s3Compatible{
			ConnectionID: customtypes.NewUUIDPointerValue(from.S3Compatible.ConnectionID),
			Location:     customtypes.NewURLPointerValue(from.S3Compatible.Location),
			Subpath:      types.StringPointerValue(from.S3Compatible.Subpath),
			Bucket:       types.StringPointerValue(from.S3Compatible.Bucket),
		}
		to.S3Compatible = supertypes.NewSingleNestedObjectValueOf(ctx, s3CompatibleModel)
	}

	if from.AzureBlobStorage != nil {
		azureBlobStorageModel := &targetDataSourceModel{
			ConnectionID: customtypes.NewUUIDPointerValue(from.AzureBlobStorage.ConnectionID),
			Location:     customtypes.NewURLPointerValue(from.AzureBlobStorage.Location),
			Subpath:      types.StringPointerValue(from.AzureBlobStorage.Subpath),
		}
		to.AzureBlobStorage = supertypes.NewSingleNestedObjectValueOf(ctx, azureBlobStorageModel)
	}

	if from.ExternalDataShare != nil {
		externalDataShareModel := &externalDataShare{
			ConnectionID: customtypes.NewUUIDPointerValue(from.ExternalDataShare.ConnectionID),
		}
		to.ExternalDataShare = supertypes.NewSingleNestedObjectValueOf(ctx, externalDataShareModel)
	}
}

/*
DATA-SOURCE
*/

type dataSourceShortcutModel struct {
	baseShortcutModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceShortcutsModel struct {
	WorkspaceID customtypes.UUID                                     `tfsdk:"workspace_id"`
	ItemID      customtypes.UUID                                     `tfsdk:"item_id"`
	Values      supertypes.SetNestedObjectValueOf[baseShortcutModel] `tfsdk:"values"`
	Timeouts    timeoutsD.Value                                      `tfsdk:"timeouts"`
}

func (to *dataSourceShortcutsModel) setValues(ctx context.Context, workspaceID, itemID string, from []fabcore.ShortcutTransformFlagged) diag.Diagnostics {
	slice := make([]*baseShortcutModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseShortcutModel

		if diags := entityModel.set(ctx, workspaceID, itemID, toShortcut(entity)); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

func toShortcut(v fabcore.ShortcutTransformFlagged) fabcore.Shortcut {
	return fabcore.Shortcut{
		Path:      v.Path,
		Name:      v.Name,
		Target:    v.Target,
		Transform: v.Transform,
	}
}

type resourceShortcutModel struct {
	baseShortcutModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateShortcut struct {
	fabcore.CreateShortcutRequest
}

func (to *requestCreateShortcut) set(ctx context.Context, from resourceShortcutModel) diag.Diagnostics {
	to.Name = from.Name.ValueStringPointer()
	to.Path = from.Path.ValueStringPointer()

	target, diags := from.Target.Get(ctx)
	if diags.HasError() {
		return diags
	}

	creatableShortcutTarget := &fabcore.CreatableShortcutTarget{}

	if !target.Onelake.IsNull() {
		entity, diags := target.Onelake.Get(ctx)
		if diags.HasError() {
			return diags
		}

		creatableShortcutTarget.OneLake = &fabcore.OneLake{
			ItemID:      entity.ItemID.ValueStringPointer(),
			Path:        entity.Path.ValueStringPointer(),
			WorkspaceID: entity.WorkspaceID.ValueStringPointer(),
		}
	}

	if !target.AdlsGen2.IsNull() {
		entity, diags := target.AdlsGen2.Get(ctx)
		if diags.HasError() {
			return diags
		}

		creatableShortcutTarget.AdlsGen2 = &fabcore.AdlsGen2{
			ConnectionID: entity.ConnectionID.ValueStringPointer(),
			Location:     entity.Location.ValueStringPointer(),
			Subpath:      entity.Subpath.ValueStringPointer(),
		}
	}

	if !target.AmazonS3.IsNull() {
		entity, diags := target.AmazonS3.Get(ctx)
		if diags.HasError() {
			return diags
		}

		creatableShortcutTarget.AmazonS3 = &fabcore.AmazonS3{
			ConnectionID: entity.ConnectionID.ValueStringPointer(),
			Location:     entity.Location.ValueStringPointer(),
			Subpath:      entity.Subpath.ValueStringPointer(),
		}
	}

	if !target.Dataverse.IsNull() {
		entity, diags := target.Dataverse.Get(ctx)
		if diags.HasError() {
			return diags
		}

		creatableShortcutTarget.Dataverse = &fabcore.Dataverse{
			ConnectionID:      entity.ConnectionID.ValueStringPointer(),
			DeltaLakeFolder:   entity.DeltaLakeFolder.ValueStringPointer(),
			EnvironmentDomain: entity.EnvironmentDomain.ValueStringPointer(),
			TableName:         entity.TableName.ValueStringPointer(),
		}
	}

	if !target.GoogleCloudStorage.IsNull() {
		entity, diags := target.GoogleCloudStorage.Get(ctx)
		if diags.HasError() {
			return diags
		}

		creatableShortcutTarget.GoogleCloudStorage = &fabcore.GoogleCloudStorage{
			ConnectionID: entity.ConnectionID.ValueStringPointer(),
			Location:     entity.Location.ValueStringPointer(),
			Subpath:      entity.Subpath.ValueStringPointer(),
		}
	}

	if !target.S3Compatible.IsNull() {
		entity, diags := target.S3Compatible.Get(ctx)
		if diags.HasError() {
			return diags
		}

		creatableShortcutTarget.S3Compatible = &fabcore.S3Compatible{
			ConnectionID: entity.ConnectionID.ValueStringPointer(),
			Location:     entity.Location.ValueStringPointer(),
			Subpath:      entity.Subpath.ValueStringPointer(),
			Bucket:       entity.Bucket.ValueStringPointer(),
		}
	}

	if !target.AzureBlobStorage.IsNull() {
		entity, diags := target.AzureBlobStorage.Get(ctx)
		if diags.HasError() {
			return diags
		}

		creatableShortcutTarget.AzureBlobStorage = &fabcore.AzureBlobStorage{
			ConnectionID: entity.ConnectionID.ValueStringPointer(),
			Location:     entity.Location.ValueStringPointer(),
			Subpath:      entity.Subpath.ValueStringPointer(),
		}
	}

	to.Target = creatableShortcutTarget

	return nil
}
