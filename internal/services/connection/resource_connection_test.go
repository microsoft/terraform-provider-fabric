// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/services/connection"
)

func TestUnit_NewResourceConnection(t *testing.T) {
	ctx := t.Context()
	r := connection.NewResourceConnection()
	require.NotNil(t, r)

	resp := resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, &resp)
	require.Empty(t, resp.Diagnostics)
	require.NotNil(t, resp.Schema)

	attributes := resp.Schema.Attributes

	// Basic attribute presence checks
	require.Contains(t, attributes, "id")
	require.Contains(t, attributes, "display_name")
	require.Contains(t, attributes, "connectivity_type")
	require.Contains(t, attributes, "gateway_id")
	require.Contains(t, attributes, "privacy_level")
	require.Contains(t, attributes, "connection_details")
	require.Contains(t, attributes, "credential_details")

	// Test id attribute
	idAttr, ok := attributes["id"].(schema.StringAttribute)
	require.True(t, ok, "id should be a StringAttribute")
	require.True(t, idAttr.Computed, "id should be computed")

	// Test display_name attribute
	displayNameAttr, ok := attributes["display_name"].(schema.StringAttribute)
	require.True(t, ok, "display_name should be a StringAttribute")
	require.True(t, displayNameAttr.Required, "display_name should be required")

	// Test connectivity_type attribute
	connectivityTypeAttr, ok := attributes["connectivity_type"].(schema.StringAttribute)
	require.True(t, ok, "connectivity_type should be a StringAttribute")
	require.True(t, connectivityTypeAttr.Required, "connectivity_type should be required")

	// Test gateway_id attribute
	gatewayIDAttr, ok := attributes["gateway_id"].(schema.StringAttribute)
	require.True(t, ok, "gateway_id should be a StringAttribute")
	require.True(t, gatewayIDAttr.Optional, "gateway_id should be optional")

	// Test privacy_level attribute
	privacyLevelAttr, ok := attributes["privacy_level"].(schema.StringAttribute)
	require.True(t, ok, "privacy_level should be a StringAttribute")
	require.False(t, privacyLevelAttr.Required, "privacy_level should be required")

	// Test connection_details nested attribute
	connectionDetailsAttr, ok := attributes["connection_details"].(schema.SingleNestedAttribute)
	require.True(t, ok, "connection_details should be a SingleNestedAttribute")
	require.True(t, connectionDetailsAttr.Required, "connection_details should be required")

	// Test connection_details nested attributes
	connDetailsAttrs := connectionDetailsAttr.Attributes
	require.Contains(t, connDetailsAttrs, "type")
	require.Contains(t, connDetailsAttrs, "creation_method")
	require.Contains(t, connDetailsAttrs, "parameters")

	// Test connection_details.type
	typeAttr, ok := connDetailsAttrs["type"].(schema.StringAttribute)
	require.True(t, ok, "connection_details.type should be a StringAttribute")
	require.True(t, typeAttr.Required, "connection_details.type should be required")

	// Test connection_details.parameters
	paramsAttr, ok := connDetailsAttrs["parameters"].(schema.ListNestedAttribute)
	require.True(t, ok, "connection_details.parameters should be a ListNestedAttribute")
	require.True(t, paramsAttr.Required, "connection_details.parameters should be optional")

	// Test credential_details nested attribute
	credentialDetailsAttr, ok := attributes["credential_details"].(schema.SingleNestedAttribute)
	require.True(t, ok, "credential_details should be a SingleNestedAttribute")
	require.True(t, credentialDetailsAttr.Required, "credential_details should be required")

	// Test credential_details nested attributes
	credDetailsAttrs := credentialDetailsAttr.Attributes
	require.Contains(t, credDetailsAttrs, "single_sign_on_type")
	require.Contains(t, credDetailsAttrs, "connection_encryption")
	require.Contains(t, credDetailsAttrs, "skip_test_connection")
	require.Contains(t, credDetailsAttrs, "credentials")

	// Test credential_details.credentials
	credentialsAttr, ok := credDetailsAttrs["credentials"].(schema.SingleNestedAttribute)
	require.True(t, ok, "credential_details.credentials should be a SingleNestedAttribute")
	require.True(t, credentialsAttr.Required, "credential_details.credentials should be required")

	// Test credential_details.credentials attributes
	credsAttrs := credentialsAttr.Attributes
	require.Contains(t, credsAttrs, "credential_type")

	// Test for credential types based on credential_type
	credTypeAttr, ok := credsAttrs["credential_type"].(schema.StringAttribute)
	require.True(t, ok, "credential_details.credentials.credential_type should be a StringAttribute")
	require.True(t, credTypeAttr.Required, "credential_details.credentials.credential_type should be required")

	// Check for specific credential fields based on credential type
	require.Contains(t, credsAttrs, "username")
	require.Contains(t, credsAttrs, "password")
}
