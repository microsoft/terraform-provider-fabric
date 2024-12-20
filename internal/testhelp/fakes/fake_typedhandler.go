// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"context"
	"reflect"
	"strings"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

var errItemNotFound = fabcore.ErrItem.ItemNotFound.Error()

// typedHandler is a handler for a specific entity type.
type typedHandler[TEntity any] struct {
	*fakeServer
	identifier identifier[TEntity]
}

// newTypedHandler creates a new typedHandler.
func newTypedHandler[TEntity any](server *fakeServer, identifier identifier[TEntity]) *typedHandler[TEntity] {
	typedHandler := &typedHandler[TEntity]{
		fakeServer: server,
		identifier: identifier,
	}

	return typedHandler
}

// ConfigureEntityWithSimpleID configures an entity with a simple ID.
func configureEntityPagerWithSimpleID[TEntity, TGetOutput, TUpdateOutput, TCreateOutput, TListOutput, TCreationData, TUpdateData, TGetOptions, TUpdateOptions, TCreateOptions, TListOptions, TDeleteOptions, TDeleteResponse any](
	handler *typedHandler[TEntity],
	operations simpleIDOperations[TEntity, TGetOutput, TUpdateOutput, TCreateOutput, TListOutput, TCreationData, TUpdateData],
	getFunction *func(ctx context.Context, id string, options *TGetOptions) (resp azfake.Responder[TGetOutput], errResp azfake.ErrorResponder),
	updateFunction *func(ctx context.Context, id string, updateRequest TUpdateData, options *TUpdateOptions) (resp azfake.Responder[TUpdateOutput], errResp azfake.ErrorResponder),
	createFunction *func(ctx context.Context, createRequest TCreationData, options *TCreateOptions) (resp azfake.Responder[TCreateOutput], errResp azfake.ErrorResponder),
	listFunction *func(options *TListOptions) (resp azfake.PagerResponder[TListOutput]),
	deleteFunction *func(ctx context.Context, id string, options *TDeleteOptions) (resp azfake.Responder[TDeleteResponse], errResp azfake.ErrorResponder),
) {
	handleGetWithSimpleID(handler, operations, getFunction)
	handleUpdateWithSimpleID(handler, operations, operations, updateFunction)
	handleCreate(handler, operations, operations, operations, createFunction)
	handleListPager(handler, operations, listFunction)
	handleDeleteWithSimpleID(handler, deleteFunction)
}

func configureEntityWithSimpleID[TEntity, TGetOutput, TUpdateOutput, TCreateOutput, TListOutput, TCreationData, TUpdateData, TGetOptions, TUpdateOptions, TCreateOptions, TListOptions, TDeleteOptions, TDeleteResponse any](
	handler *typedHandler[TEntity],
	operations simpleIDOperations[TEntity, TGetOutput, TUpdateOutput, TCreateOutput, TListOutput, TCreationData, TUpdateData],
	getFunction *func(ctx context.Context, id string, options *TGetOptions) (resp azfake.Responder[TGetOutput], errResp azfake.ErrorResponder),
	updateFunction *func(ctx context.Context, id string, updateRequest TUpdateData, options *TUpdateOptions) (resp azfake.Responder[TUpdateOutput], errResp azfake.ErrorResponder),
	createFunction *func(ctx context.Context, createRequest TCreationData, options *TCreateOptions) (resp azfake.Responder[TCreateOutput], errResp azfake.ErrorResponder),
	listFunction *func(ctx context.Context, options *TListOptions) (resp azfake.Responder[TListOutput], errResp azfake.ErrorResponder),
	deleteFunction *func(ctx context.Context, id string, options *TDeleteOptions) (resp azfake.Responder[TDeleteResponse], errResp azfake.ErrorResponder),
) {
	handleGetWithSimpleID(handler, operations, getFunction)
	handleUpdateWithSimpleID(handler, operations, operations, updateFunction)
	handleCreate(handler, operations, operations, operations, createFunction)
	handleList(handler, operations, listFunction)
	handleDeleteWithSimpleID(handler, deleteFunction)
}

// ConfigureEntityWithParentID configures an entity with a parent ID.
func configureEntityWithParentID[TEntity, TGetOutput, TUpdateOutput, TCreateOutput, TListOutput, TCreationData, TUpdateData, TGetOptions, TUpdateOptions, TCreateOptions, TListOptions, TDeleteOptions, TDeleteResponse any](
	handler *typedHandler[TEntity],
	operations parentIDOperations[TEntity, TGetOutput, TUpdateOutput, TCreateOutput, TListOutput, TCreationData, TUpdateData],
	getFunction *func(ctx context.Context, parentID, childID string, options *TGetOptions) (resp azfake.Responder[TGetOutput], errResp azfake.ErrorResponder),
	updateFunction *func(ctx context.Context, parentID, childID string, updateRequest TUpdateData, options *TUpdateOptions) (resp azfake.Responder[TUpdateOutput], errResp azfake.ErrorResponder),
	createFunction *func(ctx context.Context, parentID string, createRequest TCreationData, options *TCreateOptions) (resp azfake.PollerResponder[TCreateOutput], errResp azfake.ErrorResponder),
	listFunction *func(parentID string, options *TListOptions) (resp azfake.PagerResponder[TListOutput]),
	deleteFunction *func(ctx context.Context, parentID, childID string, options *TDeleteOptions) (resp azfake.Responder[TDeleteResponse], errResp azfake.ErrorResponder),
) {
	handleGetWithParentID(handler, operations, getFunction)
	handleUpdateWithParentID(handler, operations, operations, updateFunction)
	handleCreateLRO(handler, operations, operations, operations, createFunction)
	handleListPagerWithParentID(handler, operations, operations, listFunction)
	handleDeleteWithParentID(handler, deleteFunction)
}

// ConfigureDefinitions configures the definitions for an entity.
func configureDefinitions[TEntity, TGetOutput, TUpdateOutput, TCreateOutput, TListOutput, TCreationData, TUpdateData, TCreateOptions, TDefinition, TUpdateDefinitionOptions, TDefinitionUpdateData, TDefinitionTransformerOutput, TUpdateDefinitionTransformerOutput, TGetDefinitionsOptions any](
	handler *typedHandler[TEntity],
	entityOperations parentIDOperations[TEntity, TGetOutput, TUpdateOutput, TCreateOutput, TListOutput, TCreationData, TUpdateData],
	definitionOperations definitionOperations[TDefinition, TCreationData, TDefinitionUpdateData, TDefinitionTransformerOutput, TUpdateDefinitionTransformerOutput],
	createFunction *func(ctx context.Context, parentID string, createRequest TCreationData, options *TCreateOptions) (resp azfake.PollerResponder[TCreateOutput], errResp azfake.ErrorResponder),
	getDefinitionsFunction *func(ctx context.Context, parentID, childID string, options *TGetDefinitionsOptions) (resp azfake.PollerResponder[TDefinitionTransformerOutput], errResp azfake.ErrorResponder),
	updateDefinitionsFunction *func(ctx context.Context, parentID, childID string, updateRequest TDefinitionUpdateData, options *TUpdateDefinitionOptions) (resp azfake.PollerResponder[TUpdateDefinitionTransformerOutput], errResp azfake.ErrorResponder),
) {
	handleCreateLROWithDefinitions(handler, entityOperations, definitionOperations, entityOperations, entityOperations, createFunction)
	handleGetDefinition(handler, definitionOperations, getDefinitionsFunction)
	handleUpdateDefinition(handler, definitionOperations, updateDefinitionsFunction)
}

// GenerateID generates an ID from a parent and child ID.
func generateID(parentID, childID string) string {
	return parentID + "/" + childID
}

// Elements returns all the elements in the for the type.
func (h *typedHandler[TEntity]) Elements() []TEntity {
	ret := make([]TEntity, 0)

	// if it is a FabricItem, return all the elements as fabric items
	if h.entityTypeIsFabricItem() {
		for _, element := range h.elements {
			item := asFabricItem(element)
			ret = append(ret, h.getFabricItemAsTEntity(item))
		}
	} else {
		for _, element := range h.elements {
			if element, ok := element.(TEntity); ok {
				ret = append(ret, element)
			}
		}
	}

	return ret
}

// Delete deletes an element by ID.
func (h *typedHandler[TEntity]) Delete(id string) {
	newElements := make([]any, 0)

	for _, element := range h.elements {
		if typedElement, ok := element.(TEntity); ok {
			if h.identifier.GetID(typedElement) != id {
				newElements = append(newElements, element)
			}
		} else if h.entityTypeCanBeConvertedToFabricItem() {
			// if it wasn't found, try to find it as fabric item
			item := asFabricItem(element)
			if !strings.HasSuffix(id, *item.ID) {
				newElements = append(newElements, element)
			}
		} else {
			newElements = append(newElements, element)
		}
	}

	h.elements = newElements
}

// Upsert inserts or updates an element, replacing the existing element if it exists based on the ID.
func (h *typedHandler[TEntity]) Upsert(element TEntity) {
	id := h.identifier.GetID(element)

	// first, try to delete the element if it exists
	h.Delete(id)

	// add to the elements list
	h.elements = append(h.elements, element)
}

// Get gets an element by ID.
func (h *typedHandler[TEntity]) Get(id string) TEntity {
	// check if TEntity is FabricItem
	if h.entityTypeIsFabricItem() {
		for _, element := range h.elements {
			item := asFabricItem(element)
			if strings.HasSuffix(id, *item.ID) {
				//nolint
				return element.(TEntity)
			}
		}

		panic("Element not found") // lintignore:R009
	}

	// if it is not a FabricItem, find the element by ID
	pointer := h.getPointer(id)
	if pointer != nil {
		return *pointer
	}

	// if it still wasn't found, try to find it if they were inserted as fabric items
	if h.entityTypeCanBeConvertedToFabricItem() {
		for _, element := range h.elements {
			item := asFabricItem(element)
			if strings.HasSuffix(id, *item.ID) {
				return h.getFabricItemAsTEntity(item)
			}
		}
	}

	// if that didn't work, panic
	panic("Element not found") // lintignore:R009
}

// Contains returns true if the element exists.
func (h *typedHandler[TEntity]) Contains(id string) bool {
	found := h.getPointer(id) != nil

	if found {
		return true
	}

	// if it wasn't found, try to find it as fabric item
	if h.entityTypeCanBeConvertedToFabricItem() {
		for _, element := range h.elements {
			item := asFabricItem(element)
			if strings.HasSuffix(id, *item.ID) {
				return true
			}
		}
	}

	return false
}

// getPointer gets a pointer to an element by ID.
func (h *typedHandler[TEntity]) getPointer(id string) *TEntity {
	for _, element := range h.elements {
		if typedElement, ok := element.(TEntity); ok {
			if h.identifier.GetID(typedElement) == id {
				return &typedElement
			}
		}
	}

	return nil
}

func (h *typedHandler[TEntity]) getFabricItemAsTEntity(item fabcore.Item) TEntity {
	var entity TEntity

	h.setReflectedStringPropertyValue(&entity, "ID", *item.ID)
	h.setReflectedStringPropertyValue(&entity, "WorkspaceID", *item.WorkspaceID)
	h.setReflectedStringPropertyValue(&entity, "DisplayName", *item.DisplayName)
	h.setReflectedStringPropertyValue(&entity, "Description", *item.Description)
	h.setReflectedStringPropertyValue(&entity, "Type", string(*item.Type))

	return entity
}

// asFabricItem converts an element to a fabric item.
func asFabricItem(element any) fabcore.Item {
	itemType := fabcore.ItemType(*getReflectedStringPropertyValue(element, "Type"))

	item := fabcore.Item{
		Type:        &itemType,
		Description: getReflectedStringPropertyValue(element, "Description"),
		DisplayName: getReflectedStringPropertyValue(element, "DisplayName"),
		ID:          getReflectedStringPropertyValue(element, "ID"),
		WorkspaceID: getReflectedStringPropertyValue(element, "WorkspaceID"),
	}

	return item
}

// getReflectedStringPropertyValue gets a string property value from a reflected object.
func getReflectedStringPropertyValue(element any, propertyName string) *string {
	reflectedValue := reflect.ValueOf(element)
	propertyValue := reflectedValue.FieldByName(propertyName)

	str := propertyValue.Elem().String()

	return &str
}

// setReflectedStringPropertyValue sets a string property value on a reflected object.
func (h *typedHandler[TEntity]) setReflectedStringPropertyValue(entity *TEntity, propertyName, value string) {
	reflectedValue := reflect.ValueOf(entity).Elem()
	propertyValue := reflectedValue.FieldByName(propertyName)

	// create a new pointer to the type of the property
	ptr := reflect.New(propertyValue.Type().Elem())
	ptr.Elem().SetString(value)

	// set the value as a pointer
	propertyValue.Set(ptr)
}

func (h *typedHandler[TEntity]) entityTypeIsFabricItem() bool {
	var entity TEntity

	return reflect.TypeOf(entity) == reflect.TypeOf(fabcore.Item{})
}

func (h *typedHandler[TEntity]) entityTypeCanBeConvertedToFabricItem() bool {
	var entity TEntity

	requiredPropertyNames := []string{"ID", "WorkspaceID", "DisplayName", "Description", "Type"}

	for _, propertyName := range requiredPropertyNames {
		// check if the property exists
		if !reflect.ValueOf(entity).FieldByName(propertyName).IsValid() {
			return false
		}
	}

	return true
}
