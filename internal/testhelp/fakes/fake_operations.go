// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

// General operations that apply to every entity.
type operationsBase[
	TEntity,
	TGetEntityResponse,
	TUpdateEntityResponse,
	TCreateEntityResponse,
	TListEntityResponse,
	TCreateEntityRequest,
	TUpdateEntityRequest any] interface {
	identifier[TEntity]
	getTransformer[TEntity, TGetEntityResponse]
	updateTransformer[TEntity, TUpdateEntityResponse]
	createTransformer[TEntity, TCreateEntityResponse]
	listTransformer[TEntity, TListEntityResponse]
	updater[TUpdateEntityRequest, TEntity]
	validator[TEntity]
}

type operationsBaseOneLake[
	TEntity,
	TGetEntityResponse,
	TCreateEntityResponse,
	TListEntityResponse,
	TCreateEntityRequest any] interface {
	identifier[TEntity]
	getTransformer[TEntity, TGetEntityResponse]
	createTransformer[TEntity, TCreateEntityResponse]
	listTransformer[TEntity, TListEntityResponse]
	validator[TEntity]
}

// Operations that apply to entities with a simple ID.
type simpleIDOperations[
	TEntity,
	TGetEntityResponse,
	TUpdateEntityResponse,
	TCreateEntityResponse,
	TListEntityResponse,
	TCreateEntityRequest,
	TUpdateEntityRequest any] interface {
	operationsBase[TEntity, TGetEntityResponse, TUpdateEntityResponse, TCreateEntityResponse, TListEntityResponse, TCreateEntityRequest, TUpdateEntityRequest]
	creator[TCreateEntityRequest, TEntity]
}

// Operations that apply to entities with a parent ID + simple ID.
type parentIDOperations[
	TEntity,
	TGetEntityResponse,
	TUpdateEntityResponse,
	TCreateEntityResponse,
	TListEntityResponse,
	TCreateEntityRequest,
	TUpdateEntityRequest any] interface {
	operationsBase[TEntity, TGetEntityResponse, TUpdateEntityResponse, TCreateEntityResponse, TListEntityResponse, TCreateEntityRequest, TUpdateEntityRequest]
	creatorWithParentID[TCreateEntityRequest, TEntity]
	parentFilter[TEntity]
}

// Operations that apply to entities with a definition.
type definitionOperations[
	TDefinition,
	TCreateEntityRequest,
	TDefinitionUpdateRequest,
	TGetDefinitionResponse,
	TUpdateDefinitionResponse any] interface {
	definitionCreator[TCreateEntityRequest, TDefinition]
	definitionUpdater[TDefinitionUpdateRequest, TDefinition]
	definitionTransformer[TDefinition, TGetDefinitionResponse]
}

// Operations that apply to entities with a definition non LRO creation.
type definitionOperationsNonLROCreation[
	TDefinition,
	TDefinitionUpdateRequest,
	TGetDefinitionResponse,
	TUpdateDefinitionResponse any] interface {
	definitionUpdater[TDefinitionUpdateRequest, TDefinition]
	definitionTransformer[TDefinition, TGetDefinitionResponse]
}

type onelakeOperations[
	TEntity,
	TGetOutput,
	TCreateOutput,
	TListOutput,
	TCreateInput any,
] interface {
	identifier[TEntity]
	getTransformer[TEntity, TGetOutput]
	createTransformer[TEntity, TCreateOutput]
	listTransformer[TEntity, TListOutput]
	validator[TEntity]
	creatorWithWorkspaceIDAndItemID[TCreateInput, TEntity]
}
