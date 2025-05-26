// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelakeshortcut

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
	Path        customtypes.PathString                            `tfsdk:"path"`
	Target      supertypes.SingleNestedObjectValueOf[targetModel] `tfsdk:"target"`
	WorkspaceID customtypes.UUID                                  `tfsdk:"workspace_id"`
	ItemID      customtypes.UUID                                  `tfsdk:"item_id"`
}

type targetModel struct {
	Type               types.String                                             `tfsdk:"type"`
	Onelake            supertypes.SingleNestedObjectValueOf[oneLakeModel]       `tfsdk:"onelake"`
	AdlsGen2           supertypes.SingleNestedObjectValueOf[adlsGen2]           `tfsdk:"adls_gen2"`
	AmazonS3           supertypes.SingleNestedObjectValueOf[amazonS3]           `tfsdk:"amazon_s3"`
	Dataverse          supertypes.SingleNestedObjectValueOf[dataverse]          `tfsdk:"dataverse"`
	GoogleCloudStorage supertypes.SingleNestedObjectValueOf[googleCloudStorage] `tfsdk:"google_cloud_storage"`
	S3Compatible       supertypes.SingleNestedObjectValueOf[s3Compatible]       `tfsdk:"s3_compatible"`
	AzureBlobStorage   supertypes.SingleNestedObjectValueOf[azureBlobStorage]   `tfsdk:"azure_blob_storage"`
	ExternalDataShare  supertypes.SingleNestedObjectValueOf[externalDataShare]  `tfsdk:"external_data_share"`
}

type oneLakeModel struct {
	ItemID      customtypes.UUID       `tfsdk:"item_id"`
	Path        customtypes.PathString `tfsdk:"path"`
	WorkspaceID customtypes.UUID       `tfsdk:"workspace_id"`
}
type adlsGen2 struct {
	ConnectionID customtypes.UUID `tfsdk:"connection_id"`
	Location     customtypes.URL  `tfsdk:"location"`
	Subpath      types.String     `tfsdk:"subpath"`
}
type amazonS3 struct {
	ConnectionID customtypes.UUID `tfsdk:"connection_id"`
	Location     customtypes.URL  `tfsdk:"location"`
	Subpath      types.String     `tfsdk:"subpath"`
}
type dataverse struct {
	ConnectionID      customtypes.UUID `tfsdk:"connection_id"`
	DeltaLakeFolder   types.String     `tfsdk:"deltalake_folder"`
	EnvironmentDomain customtypes.URL  `tfsdk:"environment_domain"`
	TableName         types.String     `tfsdk:"table_name"`
}

type googleCloudStorage struct {
	ConnectionID customtypes.UUID `tfsdk:"connection_id"`
	Location     customtypes.URL  `tfsdk:"location"`
	Subpath      types.String     `tfsdk:"subpath"`
}

type s3Compatible struct {
	ConnectionID customtypes.UUID `tfsdk:"connection_id"`
	Location     customtypes.URL  `tfsdk:"location"`
	Subpath      types.String     `tfsdk:"subpath"`
	Bucket       types.String     `tfsdk:"bucket"`
}
type azureBlobStorage struct {
	ConnectionID customtypes.UUID `tfsdk:"connection_id"`
	Location     customtypes.URL  `tfsdk:"location"`
	Subpath      types.String     `tfsdk:"subpath"`
}
type externalDataShare struct {
	ConnectionID customtypes.UUID `tfsdk:"connection_id"`
}

func (to *baseShortcutModel) set(ctx context.Context, workspaceID, itemID string, from fabcore.Shortcut) diag.Diagnostics {
	to.Name = types.StringPointerValue(from.Name)
	to.Path = customtypes.NewPathStringPointerValue(from.Path)
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.ItemID = customtypes.NewUUIDValue(itemID)

	onelakeShortcutComputedID := fmt.Sprintf("%s%s%s%s", workspaceID, itemID, strings.TrimPrefix(to.Path.ValueString(), "/"), to.Name.ValueString())

	to.ID = types.StringPointerValue(&onelakeShortcutComputedID)

	to.Target = supertypes.NewSingleNestedObjectValueOfNull[targetModel](ctx)
	target := supertypes.NewSingleNestedObjectValueOfNull[targetModel](ctx)

	if from.Target != nil {
		targetModel := &targetModel{}
		if diags := targetModel.set(ctx, from.Target); diags.HasError() {
			return diags
		}

		if diags := target.Set(ctx, targetModel); diags.HasError() {
			return diags
		}
	}

	to.Target = target

	return nil
}

func (to *targetModel) set(ctx context.Context, from *fabcore.Target) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	to.Type = types.StringPointerValue((*string)(from.Type))
	to.Onelake = supertypes.NewSingleNestedObjectValueOfNull[oneLakeModel](ctx)
	to.AdlsGen2 = supertypes.NewSingleNestedObjectValueOfNull[adlsGen2](ctx)
	to.AmazonS3 = supertypes.NewSingleNestedObjectValueOfNull[amazonS3](ctx)
	to.Dataverse = supertypes.NewSingleNestedObjectValueOfNull[dataverse](ctx)
	to.GoogleCloudStorage = supertypes.NewSingleNestedObjectValueOfNull[googleCloudStorage](ctx)
	to.S3Compatible = supertypes.NewSingleNestedObjectValueOfNull[s3Compatible](ctx)
	to.AzureBlobStorage = supertypes.NewSingleNestedObjectValueOfNull[azureBlobStorage](ctx)
	to.ExternalDataShare = supertypes.NewSingleNestedObjectValueOfNull[externalDataShare](ctx)

	if from.OneLake != nil {
		onelakeModel := &oneLakeModel{
			ItemID:      customtypes.NewUUIDPointerValue(from.OneLake.ItemID),
			Path:        customtypes.NewPathStringPointerValue(from.OneLake.Path),
			WorkspaceID: customtypes.NewUUIDPointerValue(from.OneLake.WorkspaceID),
		}
		to.Onelake = supertypes.NewSingleNestedObjectValueOf(ctx, onelakeModel)
	}

	if from.AdlsGen2 != nil {
		adlsGen2Model := &adlsGen2{
			ConnectionID: customtypes.NewUUIDPointerValue(from.AdlsGen2.ConnectionID),
			Location:     customtypes.NewURLPointerValue(from.AdlsGen2.Location),
			Subpath:      types.StringPointerValue(from.AdlsGen2.Subpath),
		}
		to.AdlsGen2 = supertypes.NewSingleNestedObjectValueOf(ctx, adlsGen2Model)
	}

	if from.AmazonS3 != nil {
		amazonS3Model := &amazonS3{
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
		googleStorageCloudModel := &googleCloudStorage{
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
		azureBlobStorageModel := &azureBlobStorage{
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

	return diagnostics
}

/*
DATA-SOURCE
*/

type dataSourceOnelakeShortcutModel struct {
	baseShortcutModel
	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceOnelakeShortcutsModel struct {
	WorkspaceID customtypes.UUID                                     `tfsdk:"workspace_id"`
	ItemID      customtypes.UUID                                     `tfsdk:"item_id"`
	Values      supertypes.SetNestedObjectValueOf[baseShortcutModel] `tfsdk:"values"`
	Timeouts    timeoutsD.Value                                      `tfsdk:"timeouts"`
}

func (to *dataSourceOnelakeShortcutsModel) setValues(ctx context.Context, workspaceID, itemID string, from []fabcore.Shortcut) diag.Diagnostics {
	slice := make([]*baseShortcutModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseShortcutModel

		if diags := entityModel.set(ctx, workspaceID, itemID, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

type resourceOneLakeShortcutModel struct {
	baseShortcutModel
	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateOnelakeShortcut struct {
	fabcore.CreateShortcutRequest
}

func (to *requestCreateOnelakeShortcut) set(ctx context.Context, from resourceOneLakeShortcutModel) diag.Diagnostics {
	to.Name = from.Name.ValueStringPointer()
	to.Path = from.Path.ValueStringPointer()

	target, diags := from.Target.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.Target = &fabcore.CreatableShortcutTarget{
		OneLake: func() *fabcore.OneLake {
			onelake, diags := target.Onelake.Get(ctx)
			if diags.HasError() || onelake == nil {
				return nil
			}

			return &fabcore.OneLake{
				ItemID:      onelake.ItemID.ValueStringPointer(),
				Path:        onelake.Path.ValueStringPointer(),
				WorkspaceID: onelake.WorkspaceID.ValueStringPointer(),
			}
		}(),
		AdlsGen2: func() *fabcore.AdlsGen2 {
			adlsGen2, diags := target.AdlsGen2.Get(ctx)
			if diags.HasError() || adlsGen2 == nil {
				return nil
			}

			return &fabcore.AdlsGen2{
				ConnectionID: adlsGen2.ConnectionID.ValueStringPointer(),
				Location:     adlsGen2.Location.ValueStringPointer(),
				Subpath:      adlsGen2.Subpath.ValueStringPointer(),
			}
		}(),
		AmazonS3: func() *fabcore.AmazonS3 {
			amazonS3, diags := target.AmazonS3.Get(ctx)
			if diags.HasError() || amazonS3 == nil {
				return nil
			}

			return &fabcore.AmazonS3{
				ConnectionID: amazonS3.ConnectionID.ValueStringPointer(),
				Location:     amazonS3.Location.ValueStringPointer(),
				Subpath:      amazonS3.Subpath.ValueStringPointer(),
			}
		}(),
		Dataverse: func() *fabcore.Dataverse {
			dataverse, diags := target.Dataverse.Get(ctx)
			if diags.HasError() || dataverse == nil {
				return nil
			}

			return &fabcore.Dataverse{
				ConnectionID:      dataverse.ConnectionID.ValueStringPointer(),
				DeltaLakeFolder:   dataverse.DeltaLakeFolder.ValueStringPointer(),
				EnvironmentDomain: dataverse.EnvironmentDomain.ValueStringPointer(),
				TableName:         dataverse.TableName.ValueStringPointer(),
			}
		}(),
		GoogleCloudStorage: func() *fabcore.GoogleCloudStorage {
			googleCloudStorage, diags := target.GoogleCloudStorage.Get(ctx)
			if diags.HasError() || googleCloudStorage == nil {
				return nil
			}

			return &fabcore.GoogleCloudStorage{
				ConnectionID: googleCloudStorage.ConnectionID.ValueStringPointer(),
				Location:     googleCloudStorage.Location.ValueStringPointer(),
				Subpath:      googleCloudStorage.Subpath.ValueStringPointer(),
			}
		}(),
		S3Compatible: func() *fabcore.S3Compatible {
			s3Compatible, diags := target.S3Compatible.Get(ctx)
			if diags.HasError() || s3Compatible == nil {
				return nil
			}

			return &fabcore.S3Compatible{
				ConnectionID: s3Compatible.ConnectionID.ValueStringPointer(),
				Location:     s3Compatible.Location.ValueStringPointer(),
				Subpath:      s3Compatible.Subpath.ValueStringPointer(),
				Bucket:       s3Compatible.Bucket.ValueStringPointer(),
			}
		}(),
		AzureBlobStorage: func() *fabcore.AzureBlobStorage {
			azureBlobStorage, diags := target.GoogleCloudStorage.Get(ctx)
			if diags.HasError() || azureBlobStorage == nil {
				return nil
			}

			return &fabcore.AzureBlobStorage{
				ConnectionID: azureBlobStorage.ConnectionID.ValueStringPointer(),
				Location:     azureBlobStorage.Location.ValueStringPointer(),
				Subpath:      azureBlobStorage.Subpath.ValueStringPointer(),
			}
		}(),
	}

	return nil
}
