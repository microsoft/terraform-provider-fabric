// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package folder

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseFolderModel struct {
	ID             customtypes.UUID `tfsdk:"id"`
	DisplayName    types.String     `tfsdk:"display_name"`
	ParentFolderID customtypes.UUID `tfsdk:"parent_folder_id"`
	WorkspaceID    customtypes.UUID `tfsdk:"workspace_id"`
}

func (to *baseFolderModel) set(from fabcore.Folder) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.ParentFolderID = customtypes.NewUUIDPointerValue(from.ParentFolderID)
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
}

/*
DATA-SOURCE
*/

type dataSourceFolderModel struct {
	baseFolderModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceFoldersModel struct {
	WorkspaceID customtypes.UUID                                   `tfsdk:"workspace_id"`
	Values      supertypes.SetNestedObjectValueOf[baseFolderModel] `tfsdk:"values"`
	Timeouts    timeoutsD.Value                                    `tfsdk:"timeouts"`
}

func (to *dataSourceFoldersModel) setValues(ctx context.Context, from []fabcore.Folder) {
	slice := make([]*baseFolderModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseFolderModel

		entityModel.set(entity)

		slice = append(slice, &entityModel)
	}

	to.Values.Set(ctx, slice)
}

/*
RESOURCE
*/

type resourceFolderModel struct {
	baseFolderModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateFolder struct {
	fabcore.CreateFolderRequest
}

func (to *requestCreateFolder) set(from resourceFolderModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.ParentFolderID = from.ParentFolderID.ValueStringPointer()
}

type requestUpdateFolder struct {
	fabcore.UpdateFolderRequest
}

func (to *requestUpdateFolder) set(from resourceFolderModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
}

type requestMoveFolder struct {
	fabcore.MoveFolderRequest
}

func (to *requestMoveFolder) set(from resourceFolderModel) {
	to.TargetFolderID = from.ParentFolderID.ValueStringPointer()
}
