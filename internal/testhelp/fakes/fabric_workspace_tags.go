// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
)

func configureWorkspaceTags(server *fakeServer) {
	server.ServerFactory.Core.WorkspacesServer.ApplyWorkspaceTags = fakeApplyWorkspaceTags(server)
	server.ServerFactory.Core.WorkspacesServer.UnapplyWorkspaceTags = fakeUnapplyWorkspaceTags(server)
}

func fakeApplyWorkspaceTags(
	server *fakeServer,
) func(ctx context.Context, workspaceID string, applyTagsRequest fabcore.ApplyWorkspaceTagsRequest, options *fabcore.WorkspacesClientApplyWorkspaceTagsOptions) (resp azfake.Responder[fabcore.WorkspacesClientApplyWorkspaceTagsResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID string, req fabcore.ApplyWorkspaceTagsRequest, _ *fabcore.WorkspacesClientApplyWorkspaceTagsOptions) (azfake.Responder[fabcore.WorkspacesClientApplyWorkspaceTagsResponse], azfake.ErrorResponder) {
		var resp azfake.Responder[fabcore.WorkspacesClientApplyWorkspaceTagsResponse]
		var errResp azfake.ErrorResponder

		idx, workspace := findWorkspaceByID(server, workspaceID)
		if workspace == nil {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, fabcore.ErrWorkspace.WorkspaceNotFound.Error(), fabcore.ErrWorkspace.WorkspaceNotFound.Error()))

			return resp, errResp
		}

		for _, tagID := range req.Tags {
			displayName := "Tag " + tagID

			workspace.Tags = append(workspace.Tags, fabcore.WorkspaceAppliedTag{
				ID:          &tagID,
				DisplayName: &displayName,
			})
		}

		server.elements[idx] = *workspace

		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientApplyWorkspaceTagsResponse{}, nil)

		return resp, errResp
	}
}

func fakeUnapplyWorkspaceTags(
	server *fakeServer,
) func(ctx context.Context, workspaceID string, unapplyTagsRequest fabcore.UnapplyWorkspaceTagsRequest, options *fabcore.WorkspacesClientUnapplyWorkspaceTagsOptions) (resp azfake.Responder[fabcore.WorkspacesClientUnapplyWorkspaceTagsResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID string, req fabcore.UnapplyWorkspaceTagsRequest, _ *fabcore.WorkspacesClientUnapplyWorkspaceTagsOptions) (azfake.Responder[fabcore.WorkspacesClientUnapplyWorkspaceTagsResponse], azfake.ErrorResponder) {
		var resp azfake.Responder[fabcore.WorkspacesClientUnapplyWorkspaceTagsResponse]
		var errResp azfake.ErrorResponder

		idx, workspace := findWorkspaceByID(server, workspaceID)
		if workspace == nil {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, fabcore.ErrWorkspace.WorkspaceNotFound.Error(), fabcore.ErrWorkspace.WorkspaceNotFound.Error()))

			return resp, errResp
		}

		toRemove := make(map[string]bool, len(req.Tags))
		for _, tagID := range req.Tags {
			toRemove[tagID] = true
		}

		remaining := make([]fabcore.WorkspaceAppliedTag, 0, len(workspace.Tags))

		for _, tag := range workspace.Tags {
			if tag.ID != nil && !toRemove[*tag.ID] {
				remaining = append(remaining, tag)
			}
		}

		workspace.Tags = remaining
		server.elements[idx] = *workspace

		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientUnapplyWorkspaceTagsResponse{}, nil)

		return resp, errResp
	}
}

// findWorkspaceByID finds a workspace in the fake server's elements by workspaceID.
// Returns the index and a copy of the workspace.
func findWorkspaceByID(server *fakeServer, workspaceID string) (int, *fabcore.WorkspaceInfo) {
	for i, element := range server.elements {
		if workspace, ok := element.(fabcore.WorkspaceInfo); ok && workspace.ID != nil && *workspace.ID == workspaceID {
			return i, &workspace
		}
	}

	return -1, nil
}
