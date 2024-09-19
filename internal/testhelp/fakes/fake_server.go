// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"reflect"

	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
)

var FakeServer = newFakeServer()

// Server is a fake server that can be used to test the provider.
type fakeServer struct {
	ServerFactory *fabfake.ServerFactory
	elements      []any
	definitions   map[string]any
	types         []reflect.Type
}

// NewFakeServer creates a new fake server.
func newFakeServer() *fakeServer {
	server := &fakeServer{
		ServerFactory: &fabfake.ServerFactory{},
		elements:      make([]any, 0),
		definitions:   make(map[string]any),
		types:         make([]reflect.Type, 0),
	}

	// Register entities.
	handleEntity(server, configureItem)
	handleEntity(server, configureCapacity)
	handleEntity(server, configureDataPipeline)
	handleEntity(server, configureDomain)
	handleEntity(server, configureEventhouse)
	handleEntity(server, configureEnvironment)
	handleEntity(server, configureKQLDatabase)
	handleEntity(server, configureLakehouse)
	handleEntity(server, configureNotebook)
	handleEntity(server, configureReport)
	handleEntity(server, configureSemanticModel)
	handleEntity(server, configureWarehouse)
	handleEntity(server, configureWorkspace)

	return server
}

// HandleEntity registers an entity with the server.
// When the configureFunction is called, it is expected to register all the required handles and returns a sample of the entity.
func handleEntity[TEntity any](server *fakeServer, configureFunction func(server *fakeServer) TEntity) {
	sample := configureFunction(server)
	server.types = append(server.types, reflect.TypeOf(sample))
}

// SupportsType returns true if the server supports the given type.
func (s *fakeServer) isSupportedType(t reflect.Type) bool {
	for _, supportedType := range s.types {
		if supportedType == t {
			return true
		}
	}

	return false
}

// Upsert inserts or updates an element in the server.
// It panics if the element type is not supported.
func (s *fakeServer) Upsert(element any) {
	if !s.isSupportedType(reflect.TypeOf(element)) {
		panic("Unsupported type: " + reflect.TypeOf(element).String() + ". Did you forget to call HandleEntity in NewFakeServer?") // lintignore:R009
	}

	for i, e := range s.elements {
		if e == element {
			s.elements[i] = element

			return
		}
	}

	s.elements = append(s.elements, element)
}
