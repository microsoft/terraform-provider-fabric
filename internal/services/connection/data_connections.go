// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"
	"fmt"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceConnections)(nil)

type dataSourceConnections struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.ConnectionsClient
}

func NewDataSourceConnections() datasource.DataSource {
	return &dataSourceConnections{}
}

func (d *dataSourceConnections) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemsTFName
}

func (d *dataSourceConnections) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List a Fabric " + ItemsName + ".\n\n" +
			"Use this data source to list [" + ItemsName + "](" + ItemDocsURL + ") for more information.\n\n" +
			ItemDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"values": schema.ListNestedAttribute{
				MarkdownDescription: fmt.Sprintf("The list of %s.", ItemsName),
				Computed:            true,
				CustomType:          supertypes.NewListNestedObjectTypeOf[baseConnectionModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The object ID of the connection.",
							Computed:            true,
							CustomType:          customtypes.UUIDType{},
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the connection.",
							Computed:            true,
						},
						"gateway_id": schema.StringAttribute{
							MarkdownDescription: "The gateway object ID of the connection.",
							Computed:            true,
							CustomType:          customtypes.UUIDType{},
						},
						"connectivity_type": schema.StringAttribute{
							MarkdownDescription: "The connectivity type of the connection. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleConnectivityTypeValues(), true, true),
							Computed:            true,
						},
						"privacy_level": schema.StringAttribute{
							MarkdownDescription: "The privacy level of the connection. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossiblePrivacyLevelValues(), true, true),
							Computed:            true,
						},
						"connection_details": schema.SingleNestedAttribute{
							MarkdownDescription: "The connection details of the connection.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[connectionDetailsModel](ctx),
							Attributes: map[string]schema.Attribute{
								"path": schema.StringAttribute{
									MarkdownDescription: "The path of the connection.",
									Computed:            true,
								},
								"type": schema.StringAttribute{
									MarkdownDescription: "The type of the connection.",
									Computed:            true,
								},
							},
						},
						"credential_details": schema.SingleNestedAttribute{
							MarkdownDescription: "The credential details of the connection.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[credentialDetailsModel](ctx),
							Attributes: map[string]schema.Attribute{
								"connection_encryption": schema.StringAttribute{
									MarkdownDescription: "The connection encryption type of the connection. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleConnectionEncryptionValues(), true, true),
									Computed:            true,
								},
								"credential_type": schema.StringAttribute{
									MarkdownDescription: "The credential type of the connection. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleCredentialTypeValues(), true, true),
									Computed:            true,
								},
								"single_sign_on_type": schema.StringAttribute{
									MarkdownDescription: "The single sign-on type of the connection. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleSingleSignOnTypeValues(), true, true),
									Computed:            true,
								},
								"skip_test_connection": schema.BoolAttribute{
									MarkdownDescription: "Whether the connection should skip the test connection during creation and update. `True` - Skip the test connection, `False` - Do not skip the test connection.",
									Computed:            true,
								},
							},
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceConnections) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	pConfigData, ok := req.ProviderData.(*pconfig.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			common.ErrorDataSourceConfigType,
			fmt.Sprintf(common.ErrorFabricClientType, req.ProviderData),
		)

		return
	}

	d.pConfigData = pConfigData
	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewConnectionsClient()
}

func (d *dataSourceConnections) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceConnectionsModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, d.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(d.list(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dataSourceConnections) list(ctx context.Context, model *dataSourceConnectionsModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting "+ItemsName)

	respList, err := d.client.ListConnections(ctx, nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	return model.setValues(ctx, respList)
}
