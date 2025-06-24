// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package folder_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

func fakeMoveFolder(
	folder fabcore.Folder,
) func(ctx context.Context, workspaceID, folderID string, moveFolderRequest fabcore.MoveFolderRequest, options *fabcore.FoldersClientMoveFolderOptions) (resp azfake.Responder[fabcore.FoldersClientMoveFolderResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, folderID string, moveReq fabcore.MoveFolderRequest, _ *fabcore.FoldersClientMoveFolderOptions) (resp azfake.Responder[fabcore.FoldersClientMoveFolderResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.FoldersClientMoveFolderResponse]{}
		updatedFolder := fabcore.Folder{
			ID:             to.Ptr(folderID),
			WorkspaceID:    to.Ptr(workspaceID),
			DisplayName:    folder.DisplayName,
			ParentFolderID: moveReq.TargetFolderID,
		}

		resp.SetResponse(http.StatusOK, fabcore.FoldersClientMoveFolderResponse{
			Folder: updatedFolder,
		}, nil)
		return
	}
}

func fakeGetFolder(
	folder fabcore.Folder,
) func(ctx context.Context, workspaceID, folderID string, options *fabcore.FoldersClientGetFolderOptions) (resp azfake.Responder[fabcore.FoldersClientGetFolderResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, folderID string, _ *fabcore.FoldersClientGetFolderOptions) (resp azfake.Responder[fabcore.FoldersClientGetFolderResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.FoldersClientGetFolderResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.FoldersClientGetFolderResponse{
			Folder: fabcore.Folder{
				ID:             to.Ptr(folderID),
				WorkspaceID:    to.Ptr(workspaceID),
				DisplayName:    folder.DisplayName,
				ParentFolderID: folder.ParentFolderID,
			},
		}, nil)

		return
	}
}
