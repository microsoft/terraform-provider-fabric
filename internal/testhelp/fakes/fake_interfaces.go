// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

type identifier[TEntity any] interface {
	// GetID returns the ID of the entity.
	GetID(entity TEntity) string
}

type getTransformer[TEntity, TOutput any] interface {
	// TransformGet transforms an entity into a response.
	TransformGet(entity TEntity) TOutput
}

type updateTransformer[TEntity, TOutput any] interface {
	// TransformUpdate transforms an entity into a response.
	TransformUpdate(entity TEntity) TOutput
}

type createTransformer[TEntity, TOutput any] interface {
	// TransformCreate transforms an entity into a response.
	TransformCreate(entity TEntity) TOutput
}

type listTransformer[TEntity, TOutput any] interface {
	// TransformList transforms a list of entities into a response.
	TransformList(list []TEntity) TOutput
}
type validator[TEntity any] interface {
	// Validate validates the entity against existing entities in the server.
	Validate(newEntity TEntity, existing []TEntity) (statusCode int, err error)
}

type updater[TUpdateData, TEntity any] interface {
	// Update updates the entity with the given data.
	Update(base TEntity, data TUpdateData) TEntity
}

type creator[TCreationData, TEntity any] interface {
	// Create creates an entity with the given data.
	Create(data TCreationData) TEntity
}

type creatorWithParentID[TCreationData, TEntity any] interface {
	// CreateWithParentID creates an entity with the given data and parent ID.
	CreateWithParentID(parentID string, data TCreationData) TEntity
}

type parentFilter[TEntity any] interface {
	// Filter filters the elements based on the parent ID.
	Filter(elements []TEntity, parentID string) []TEntity
}

type definitionCreator[TDefinitionCreationData, TDefinition any] interface {
	// CreateDefinition creates a definition with the given data.
	CreateDefinition(data TDefinitionCreationData) *TDefinition
}

type definitionUpdater[TDefinitionUpdateData, TDefinition any] interface {
	// UpdateDefinition updates the definition with the given data.
	UpdateDefinition(base *TDefinition, data TDefinitionUpdateData) *TDefinition
}

type definitionTransformer[TDefinition, TOutput any] interface {
	// TransformDefinition transforms a definition into a response.
	TransformDefinition(entity *TDefinition) TOutput
}
