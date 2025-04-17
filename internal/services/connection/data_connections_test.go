// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/services/connection"
)

func TestUnit_NewDataSourceConnections(t *testing.T) {
	ctx := t.Context()
	ds := connection.NewDataSourceConnections()
	require.NotNil(t, ds)

	resp := datasource.SchemaResponse{}
	ds.Schema(ctx, datasource.SchemaRequest{}, &resp)
	require.Empty(t, resp.Diagnostics)
	require.NotNil(t, resp.Schema)

	attributes := resp.Schema.Attributes

	require.Contains(t, attributes, "workspace_id")
	require.Contains(t, attributes, "values")
}
