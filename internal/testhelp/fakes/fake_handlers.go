// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
)

func handleDeleteWithSimpleID[TEntity, TOptions, TDeleteResponse any](
	handler *typedHandler[TEntity],
	f *func(ctx context.Context, id string, options *TOptions) (resp azfake.Responder[TDeleteResponse], errResp azfake.ErrorResponder),
) {
	if f == nil {
		return
	}

	*f = func(_ context.Context, id string, _ *TOptions) (azfake.Responder[TDeleteResponse], azfake.ErrorResponder) {
		return deleteByID[TEntity, TDeleteResponse](handler, id)
	}
}

func handleDeleteWithParentID[TEntity, TOptions, TResponse any](
	handler *typedHandler[TEntity],
	f *func(ctx context.Context, parentID, childID string, options *TOptions) (resp azfake.Responder[TResponse], errResp azfake.ErrorResponder),
) {
	if f == nil {
		return
	}

	*f = func(_ context.Context, parentID, childID string, _ *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder) {
		id := generateID(parentID, childID)

		return deleteByID[TEntity, TResponse](handler, id)
	}
}

func handleGetWithSimpleID[TEntity, TOptions, TResponse any](
	handler *typedHandler[TEntity],
	getTransformer getTransformer[TEntity, TResponse],
	f *func(ctx context.Context, id string, options *TOptions) (resp azfake.Responder[TResponse], errResp azfake.ErrorResponder),
) {
	if f == nil {
		return
	}

	*f = func(_ context.Context, id string, _ *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder) {
		return getByID(handler, id, getTransformer)
	}
}

func handleGetWithParentID[TEntity, TOptions, TResponse any](
	handler *typedHandler[TEntity],
	getTransformer getTransformer[TEntity, TResponse],
	f *func(ctx context.Context, parentID, childID string, options *TOptions) (resp azfake.Responder[TResponse], errResp azfake.ErrorResponder),
) {
	if f == nil {
		return
	}

	*f = func(_ context.Context, parentID, childID string, _ *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder) {
		id := generateID(parentID, childID)

		return getByID(handler, id, getTransformer)
	}
}

func handleGetOnelakeShortcut[TEntity, TOptions, TResponse any](
	handler *typedHandler[TEntity],
	getTransformer getTransformer[TEntity, TResponse],
	f *func(ctx context.Context, workspaceID, itemID, path, name string, options *TOptions) (resp azfake.Responder[TResponse], errResp azfake.ErrorResponder),
) {
	if f == nil {
		return
	}

	*f = func(_ context.Context, workspaceID, itemID, path, name string, _ *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder) {
		id := generateID(path, name)

		return getByID(handler, id, getTransformer)
	}
}

func handleGetDefinition[TEntity, TDefinition, TOptions, TResponse any](
	handler *typedHandler[TEntity],
	definitionTransformer definitionTransformer[TDefinition, TResponse],
	function *func(ctx context.Context, parentID, childID string, options *TOptions) (azfake.PollerResponder[TResponse], azfake.ErrorResponder),
) {
	if function == nil {
		return
	}

	*function = func(_ context.Context, parentID, childID string, _ *TOptions) (azfake.PollerResponder[TResponse], azfake.ErrorResponder) { //nolint:unparam
		var resp azfake.PollerResponder[TResponse]

		var errResp azfake.ErrorResponder

		id := generateID(parentID, childID)

		if definition, ok := handler.definitions[id]; ok {
			typedDefinition, ok := definition.(TDefinition)

			if !ok {
				panic("Definition not of the expected type") // lintignore:R009
			}

			respValue := definitionTransformer.TransformDefinition(&typedDefinition)
			resp.SetTerminalResponse(http.StatusOK, respValue, nil)
		} else {
			respValue := definitionTransformer.TransformDefinition(nil)
			resp.SetTerminalResponse(http.StatusOK, respValue, nil)
		}

		return resp, errResp
	}
}

func handleUpdateWithSimpleID[TEntity, TOptions, TUpdateRequest, TResponse any](
	handler *typedHandler[TEntity],
	updater updater[TUpdateRequest, TEntity],
	updateTransformer updateTransformer[TEntity, TResponse],
	f *func(ctx context.Context, id string, updateRequest TUpdateRequest, options *TOptions) (resp azfake.Responder[TResponse], errResp azfake.ErrorResponder),
) {
	if f == nil {
		return
	}

	*f = func(_ context.Context, id string, updateRequest TUpdateRequest, _ *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder) {
		return updateByID(handler, id, updateRequest, updater, updateTransformer)
	}
}

func handleUpdateWithParentID[TEntity, TOptions, TUpdateRequest, TResponse any](
	handler *typedHandler[TEntity],
	updater updater[TUpdateRequest, TEntity],
	updateTransformer updateTransformer[TEntity, TResponse],
	f *func(ctx context.Context, parentID, childID string, updateRequest TUpdateRequest, options *TOptions) (resp azfake.Responder[TResponse], errResp azfake.ErrorResponder),
) {
	if f == nil {
		return
	}

	*f = func(_ context.Context, parentID, childID string, updateRequest TUpdateRequest, _ *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder) {
		id := generateID(parentID, childID)

		return updateByID(handler, id, updateRequest, updater, updateTransformer)
	}
}

func getDefinition[TDefinition, TEntity any](handler *typedHandler[TEntity], id string) *TDefinition {
	if definition, ok := handler.definitions[id]; ok {
		typedDefinition, ok := definition.(TDefinition)

		if !ok {
			panic("Definition not of the expected type") // lintignore:R009
		}

		return &typedDefinition
	}

	return nil
}

func upsertDefinition[TDefinition, TEntity any](handler *typedHandler[TEntity], id string, definition *TDefinition) {
	if definition != nil {
		handler.definitions[id] = *definition
	} else {
		delete(handler.definitions, id)
	}
}

func handleUpdateDefinition[TEntity, TDefinition, TRequest, TOptions, TResponse any](
	handler *typedHandler[TEntity],
	definitionUpdater definitionUpdater[TRequest, TDefinition],
	function *func(ctx context.Context, parentID, childID string, request TRequest, options *TOptions) (azfake.PollerResponder[TResponse], azfake.ErrorResponder),
) {
	if function == nil {
		return
	}

	*function = func(_ context.Context, parentID, childID string, request TRequest, _ *TOptions) (azfake.PollerResponder[TResponse], azfake.ErrorResponder) { //nolint:unparam
		var resp azfake.PollerResponder[TResponse]

		var errResp azfake.ErrorResponder

		id := generateID(parentID, childID)

		typedDefinition := getDefinition[TDefinition](handler, id)

		updatedDefinition := definitionUpdater.UpdateDefinition(typedDefinition, request)
		upsertDefinition[TDefinition](handler, id, updatedDefinition)

		var respValue TResponse

		resp.SetTerminalResponse(http.StatusOK, respValue, nil)

		return resp, errResp
	}
}

func handleCreateWithoutWorkspace[TEntity, TOptions, TCreateRequest, TResponse any](
	handler *typedHandler[TEntity],
	creator creator[TCreateRequest, TEntity],
	validator validator[TEntity],
	createTransformer createTransformer[TEntity, TResponse],
	f *func(ctx context.Context, createRequest TCreateRequest, options *TOptions) (resp azfake.Responder[TResponse], errResp azfake.ErrorResponder),
) {
	if f == nil {
		return
	}

	*f = func(_ context.Context, createRequest TCreateRequest, _ *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder) {
		var resp azfake.Responder[TResponse]

		var errResp azfake.ErrorResponder

		newEntity := creator.Create(createRequest)

		if statusCode, err := validator.Validate(newEntity, handler.Elements()); err != nil {
			var empty TEntity
			respValue := createTransformer.TransformCreate(empty)
			resp.SetResponse(statusCode, respValue, nil)

			errResp.SetError(err)
			// errResp.SetResponseError(statusCode, err.Error())
		} else {
			handler.Upsert(newEntity)

			respValue := createTransformer.TransformCreate(newEntity)
			resp.SetResponse(statusCode, respValue, nil)
		}

		return resp, errResp
	}
}

func handleListOnelakeShortcut[TEntity, TOptions, TResponse any](
	handler *typedHandler[TEntity],
	listTransformer listTransformer[TEntity, TResponse],
	f *func(workspaceID, itemID string, options *TOptions) azfake.PagerResponder[TResponse],
) {
	if f == nil {
		return
	}

	*f = func(workspaceID, itemID string, options *TOptions) azfake.PagerResponder[TResponse] {
		var resp azfake.PagerResponder[TResponse]

		elements := handler.Elements()

		respValue := listTransformer.TransformList(elements)
		resp.AddPage(http.StatusOK, respValue, nil)

		return resp
	}
}

func handleDeleteOnelakeShortcut[TEntity, TOptions, TResponse any](
	handler *typedHandler[TEntity],
	f *func(ctx context.Context, workspaceID, itemID, path, name string, options *TOptions) (resp azfake.Responder[TResponse], errResp azfake.ErrorResponder),
) {
	if f == nil {
		return
	}

	*f = func(_ context.Context, workspaceID, itemID, path, name string, _ *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder) {
		id := generateID(path, name)
		return deleteByID[TEntity, TResponse](handler, id)
	}
}

func handleNonLROCreateOnelakeShortcut[TEntity, TOptions, TCreateRequest, TResponse any](
	handler *typedHandler[TEntity],
	creator creatorWithWorkspaceIDAndItemID[TCreateRequest, TEntity],
	validator validator[TEntity],
	createTransformer createTransformer[TEntity, TResponse],
	f *func(ctx context.Context, workspaceID, itemID string, createRequest TCreateRequest, options *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder),
) {
	if f == nil {
		return
	}

	*f = func(_ context.Context, workspaceID, itemID string, createRequest TCreateRequest, _ *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder) {
		var resp azfake.Responder[TResponse]
		var errResp azfake.ErrorResponder

		newEntity := creator.CreateWithWorkspaceIDAndItemID(workspaceID, itemID, createRequest)

		if statusCode, err := validator.Validate(newEntity, handler.Elements()); err != nil {
			var empty TEntity
			respValue := createTransformer.TransformCreate(empty)
			resp.SetResponse(statusCode, respValue, nil)
			errResp.SetError(err)
		} else {
			handler.Upsert(newEntity)
			respValue := createTransformer.TransformCreate(newEntity)
			resp.SetResponse(statusCode, respValue, nil)
		}

		return resp, errResp
	}
}

func handleNonLROCreate[TEntity, TOptions, TCreateRequest, TResponse any](
	handler *typedHandler[TEntity],
	creator creatorWithParentID[TCreateRequest, TEntity],
	validator validator[TEntity],
	createTransformer createTransformer[TEntity, TResponse],
	f *func(ctx context.Context, parentID string, createRequest TCreateRequest, options *TOptions) (resp azfake.Responder[TResponse], errResp azfake.ErrorResponder),
) {
	if f == nil {
		return
	}

	*f = func(_ context.Context, parentID string, createRequest TCreateRequest, _ *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder) {
		var resp azfake.Responder[TResponse]
		var errResp azfake.ErrorResponder

		newEntity := creator.CreateWithParentID(parentID, createRequest)

		if statusCode, err := validator.Validate(newEntity, handler.Elements()); err != nil {
			var empty TEntity
			respValue := createTransformer.TransformCreate(empty)
			resp.SetResponse(statusCode, respValue, nil)
			errResp.SetError(err)
		} else {
			handler.Upsert(newEntity)
			respValue := createTransformer.TransformCreate(newEntity)
			resp.SetResponse(statusCode, respValue, nil)
		}

		return resp, errResp
	}
}

// func handleNonLROUpdateDefinition[TEntity, TDefinition, TRequest, TOptions, TResponse any](handler *typedHandler[TEntity], definitionUpdater definitionUpdater[TRequest, TDefinition], function *func(ctx context.Context, parentID, childID string, request TRequest, options *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder)) {
// 	if function == nil {
// 		return
// 	}

// 	*function = func(_ context.Context, parentID, childID string, request TRequest, _ *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder) { //nolint:unparam
// 		var resp azfake.Responder[TResponse]

// 		var errResp azfake.ErrorResponder

// 		id := generateID(parentID, childID)

// 		typedDefinition := getDefinition[TDefinition](handler, id)

// 		updatedDefinition := definitionUpdater.UpdateDefinition(typedDefinition, request)
// 		upsertDefinition(handler, id, updatedDefinition)

// 		var respValue TResponse

// 		resp.SetResponse(http.StatusOK, respValue, nil)

// 		return resp, errResp
// 	}
// }

// func handleNonLROGetDefinition[TEntity, TDefinition, TOptions, TResponse any](handler *typedHandler[TEntity], definitionTransformer definitionTransformer[TDefinition, TResponse], function *func(ctx context.Context, parentID, childID string, options *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder)) {
// 	if function == nil {
// 		return
// 	}

// 	*function = func(_ context.Context, parentID, childID string, _ *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder) { //nolint:unparam
// 		var resp azfake.Responder[TResponse]

// 		var errResp azfake.ErrorResponder

// 		id := generateID(parentID, childID)

// 		if definition, ok := handler.definitions[id]; ok {
// 			typedDefinition, ok := definition.(TDefinition)

// 			if !ok {
// 				panic("Definition not of the expected type") // lintignore:R009
// 			}

// 			respValue := definitionTransformer.TransformDefinition(&typedDefinition)
// 			resp.SetResponse(http.StatusOK, respValue, nil)
// 		} else {
// 			respValue := definitionTransformer.TransformDefinition(nil)
// 			resp.SetResponse(http.StatusOK, respValue, nil)
// 		}

// 		return resp, errResp
// 	}
// }

func handleCreateLRO[TEntity, TOptions, TCreateRequest, TResponse any](
	h *typedHandler[TEntity],
	creator creatorWithParentID[TCreateRequest, TEntity],
	validator validator[TEntity],
	createTransformer createTransformer[TEntity, TResponse],
	f *func(ctx context.Context, parentID string, createRequest TCreateRequest, options *TOptions) (resp azfake.PollerResponder[TResponse], errResp azfake.ErrorResponder),
) {
	if f == nil {
		return
	}

	*f = createLROWithOptionalDefinition[TEntity, TOptions, TCreateRequest, any](h, creator, nil, validator, createTransformer)
}

func handleCreateLROWithDefinitions[TEntity, TOptions, TCreateRequest, TDefinition, TResponse any](
	h *typedHandler[TEntity],
	creator creatorWithParentID[TCreateRequest, TEntity],
	definitionCreator definitionCreator[TCreateRequest, TDefinition],
	validator validator[TEntity],
	createTransformer createTransformer[TEntity, TResponse],
	f *func(ctx context.Context, parentID string, createRequest TCreateRequest, options *TOptions) (resp azfake.PollerResponder[TResponse], errResp azfake.ErrorResponder),
) {
	if f == nil {
		return
	}

	*f = createLROWithOptionalDefinition[TEntity, TOptions](h, creator, definitionCreator, validator, createTransformer)
}

func handleListPagerWithParentID[TEntity, TOptions, TResponse any](
	h *typedHandler[TEntity],
	filter parentFilter[TEntity],
	listTransformer listTransformer[TEntity, TResponse],
	f *func(parentID string, options *TOptions) (resp azfake.PagerResponder[TResponse]),
) {
	if f == nil {
		return
	}

	*f = listPagerWithFilter[TEntity, TOptions](h, filter, listTransformer)
}

func handleListPager[TEntity, TOptions, TResponse any](
	handler *typedHandler[TEntity],
	listTransformer listTransformer[TEntity, TResponse],
	f *func(options *TOptions) (resp azfake.PagerResponder[TResponse]),
) {
	if f == nil {
		return
	}

	*f = func(options *TOptions) azfake.PagerResponder[TResponse] {
		return listPagerWithFilter[TEntity, TOptions](handler, nil, listTransformer)("", options)
	}
}

func handleList[TEntity, TOptions, TResponse any](
	handler *typedHandler[TEntity],
	listTransformer listTransformer[TEntity, TResponse],
	f *func(ctx context.Context, options *TOptions) (resp azfake.Responder[TResponse], errResp azfake.ErrorResponder),
) {
	if f == nil {
		return
	}

	*f = func(_ context.Context, options *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder) {
		return listWithFilter[TEntity, TOptions](handler, nil, listTransformer)("", options)
	}
}

func createLROWithOptionalDefinition[TEntity, TOptions, TCreateRequest, TDefinition, TResponse any](
	handler *typedHandler[TEntity],
	creator creatorWithParentID[TCreateRequest, TEntity],
	definitionCreator definitionCreator[TCreateRequest, TDefinition],
	validator validator[TEntity],
	createTransformer createTransformer[TEntity, TResponse],
) func(ctx context.Context, parentID string, createRequest TCreateRequest, options *TOptions) (resp azfake.PollerResponder[TResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, parentID string, createRequest TCreateRequest, _ *TOptions) (azfake.PollerResponder[TResponse], azfake.ErrorResponder) {
		var resp azfake.PollerResponder[TResponse]

		var errResp azfake.ErrorResponder

		newEntity := creator.CreateWithParentID(parentID, createRequest)

		if statusCode, err := validator.Validate(newEntity, handler.Elements()); err != nil {
			var empty TEntity
			respValue := createTransformer.TransformCreate(empty)
			resp.SetTerminalResponse(statusCode, respValue, nil)

			errResp.SetError(err)
			// errResp.SetResponseError(statusCode, err.Error())
		} else {
			handler.Upsert(newEntity)

			if definitionCreator != nil {
				definition := definitionCreator.CreateDefinition(createRequest)
				if definition != nil {
					id := handler.identifier.GetID(newEntity)
					upsertDefinition[TDefinition](handler, id, definition)
				}
			}

			respValue := createTransformer.TransformCreate(newEntity)
			resp.SetTerminalResponse(statusCode, respValue, nil)
		}

		return resp, errResp
	}
}

func listPagerWithFilter[TEntity, TOptions, TResponse any](
	h *typedHandler[TEntity],
	filter parentFilter[TEntity],
	listTransformer listTransformer[TEntity, TResponse],
) func(parentID string, options *TOptions) (resp azfake.PagerResponder[TResponse]) {
	return func(parentID string, _ *TOptions) azfake.PagerResponder[TResponse] {
		var resp azfake.PagerResponder[TResponse]

		elements := h.Elements()

		if filter != nil {
			elements = filter.Filter(elements, parentID)
		}

		respValue := listTransformer.TransformList(elements)

		resp.AddPage(http.StatusOK, respValue, nil)

		return resp
	}
}

func listWithFilter[TEntity, TOptions, TResponse any](
	h *typedHandler[TEntity],
	filter parentFilter[TEntity],
	listTransformer listTransformer[TEntity, TResponse],
) func(parentID string, options *TOptions) (resp azfake.Responder[TResponse], errResp azfake.ErrorResponder) {
	return func(parentID string, _ *TOptions) (azfake.Responder[TResponse], azfake.ErrorResponder) {
		var resp azfake.Responder[TResponse]
		var errResp azfake.ErrorResponder

		elements := h.Elements()

		if filter != nil {
			elements = filter.Filter(elements, parentID)
		}

		respValue := listTransformer.TransformList(elements)

		resp.SetResponse(http.StatusOK, respValue, nil)

		return resp, errResp
	}
}

func updateByID[TEntity, TUpdateRequest, TResponse any](
	handler *typedHandler[TEntity],
	id string,
	updateRequest TUpdateRequest,
	updater updater[TUpdateRequest, TEntity],
	updateTransformer updateTransformer[TEntity, TResponse],
) (azfake.Responder[TResponse], azfake.ErrorResponder) {
	var resp azfake.Responder[TResponse]

	var errResp azfake.ErrorResponder

	var empty TResponse

	if handler.Contains(id) {
		element := handler.Get(id)
		updatedElement := updater.Update(element, updateRequest)
		handler.Upsert(updatedElement)

		respValue := updateTransformer.TransformUpdate(updatedElement)
		resp.SetResponse(http.StatusOK, respValue, nil)
	} else {
		resp.SetResponse(http.StatusNotFound, empty, nil)

		errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errItemNotFound, "Item not found"))
	}

	return resp, errResp
}

func deleteByID[TEntity, TResponse any](handler *typedHandler[TEntity], id string) (azfake.Responder[TResponse], azfake.ErrorResponder) {
	var resp azfake.Responder[TResponse]

	var errResp azfake.ErrorResponder

	var empty TResponse

	if handler.Contains(id) {
		handler.Delete(id)
		resp.SetResponse(http.StatusOK, empty, nil)
	} else {
		resp.SetResponse(http.StatusNotFound, empty, nil)

		errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errItemNotFound, "Item not found"))
	}

	return resp, errResp
}

func getByID[TEntity, TResponse any](handler *typedHandler[TEntity], id string, getTransformer getTransformer[TEntity, TResponse]) (azfake.Responder[TResponse], azfake.ErrorResponder) {
	var resp azfake.Responder[TResponse]

	var errResp azfake.ErrorResponder

	var empty TResponse

	if handler.Contains(id) {
		element := handler.Get(id)
		respValue := getTransformer.TransformGet(element)
		resp.SetResponse(http.StatusOK, respValue, nil)
	} else {
		resp.SetResponse(http.StatusNotFound, empty, nil)

		errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errItemNotFound, "Item not found"))
	}

	return resp, errResp
}
