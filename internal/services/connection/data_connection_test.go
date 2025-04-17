// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/services/connection"
)

func TestUnit_NewDataSourceConnection(t *testing.T) {
	ctx := t.Context()
	ds := connection.NewDataSourceConnection()
	require.NotNil(t, ds)

	resp := datasource.SchemaResponse{}
	ds.Schema(ctx, datasource.SchemaRequest{}, &resp)
	require.Empty(t, resp.Diagnostics)
	require.NotNil(t, resp.Schema)

	attributes := resp.Schema.Attributes

	require.Contains(t, attributes, "id")
	require.Contains(t, attributes, "display_name")
	require.Contains(t, attributes, "connectivity_type")
	require.Contains(t, attributes, "gateway_id")
	require.Contains(t, attributes, "privacy_level")
	require.Contains(t, attributes, "connection_details")
	require.Contains(t, attributes, "credential_details")
}
