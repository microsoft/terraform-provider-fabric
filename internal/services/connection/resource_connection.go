// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure        = (*resourceConnection)(nil)
	_ resource.ResourceWithModifyPlan       = (*resourceConnection)(nil)
	_ resource.ResourceWithConfigValidators = (*resourceConnection)(nil)
)

type resourceConnection struct {
	Name        string
	TFName      string
	IsPreview   bool
	pConfigData *pconfig.ProviderData
	client      *fabcore.ConnectionsClient
	clientGw    *fabcore.GatewaysClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceConnection() resource.Resource {
	return &resourceConnection{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceConnection) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceConnection) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema(ctx, false).GetResource(ctx)
}

func (r *resourceConnection) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("credential_details").AtName("basic_credentials"),
			path.MatchRoot("credential_details").AtName("key_credentials"),
			path.MatchRoot("credential_details").AtName("service_principal_credentials"),
			path.MatchRoot("credential_details").AtName("shared_access_signature_credentials"),
		),
	}
}

func (r *resourceConnection) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	pConfigData, ok := req.ProviderData.(*pconfig.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			common.ErrorResourceConfigType,
			fmt.Sprintf(common.ErrorFabricClientType, req.ProviderData),
		)

		return
	}

	r.pConfigData = pConfigData
	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewConnectionsClient()
	r.clientGw = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewGatewaysClient()

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(r.Name, r.IsPreview, r.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceConnection) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	tflog.Debug(ctx, "MODIFY PLAN", map[string]any{
		"action": "start",
	})

	//nolint:nestif
	if !req.Plan.Raw.IsNull() {
		var plan resourceConnectionModel[rsConnectionDetailsModel, rsCredentialDetailsModel]

		if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
			return
		}

		connectionDetails, diags := plan.getConnectionDetails(ctx)
		if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
			return
		}

		var supportedConnectionType fabcore.ConnectionCreationMetadata

		if resp.Diagnostics.Append(r.getConnectionTypeMetadata(ctx, *connectionDetails, &supportedConnectionType)...); resp.Diagnostics.HasError() {
			return
		}

		if resp.Diagnostics.Append(r.validateCreationMethod(*connectionDetails, supportedConnectionType.CreationMethods)...); resp.Diagnostics.HasError() {
			return
		}

		credentialDetails, diags := plan.getCredentialDetails(ctx)
		if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
			return
		}

		if resp.Diagnostics.Append(r.validateConnectionEncryption(*credentialDetails, supportedConnectionType.SupportedConnectionEncryptionTypes)...); resp.Diagnostics.HasError() {
			return
		}

		if resp.Diagnostics.Append(r.validateCredentialType(*credentialDetails, supportedConnectionType.SupportedCredentialTypes)...); resp.Diagnostics.HasError() {
			return
		}

		if resp.Diagnostics.Append(r.validateSkipTestConnection(*credentialDetails, *supportedConnectionType.SupportsSkipTestConnection)...); resp.Diagnostics.HasError() {
			return
		}

		if resp.Diagnostics.Append(r.validateCreationMethodParameters(ctx, *connectionDetails, supportedConnectionType.CreationMethods)...); resp.Diagnostics.HasError() {
			return
		}

		if resp.Diagnostics.Append(r.setConnectionPerametersDataType(ctx, supportedConnectionType.CreationMethods, &plan)...); resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	}

	tflog.Debug(ctx, "MODIFY PLAN", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceConnection) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan, config resourceConnectionModel[rsConnectionDetailsModel, rsCredentialDetailsModel]

	if resp.Diagnostics.Append(req.Config.Get(ctx, &config)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqCreate requestCreateConnection

	if resp.Diagnostics.Append(reqCreate.set(ctx, plan, config)...); resp.Diagnostics.HasError() {
		return
	}

	respCreate, err := r.client.CreateConnection(ctx, reqCreate.CreateConnectionRequestClassification, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, respCreate.ConnectionClassification)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceConnection) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceConnectionModel[rsConnectionDetailsModel, rsCredentialDetailsModel]

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Read(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	diags = r.get(ctx, &state)
	if utils.IsErrNotFound(state.ID.ValueString(), &diags, fabcore.ErrWorkspace.WorkspaceNotFound) {
		resp.State.RemoveResource(ctx)

		resp.Diagnostics.Append(diags...)

		return
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceConnection) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan, config resourceConnectionModel[rsConnectionDetailsModel, rsCredentialDetailsModel]

	if resp.Diagnostics.Append(req.Config.Get(ctx, &config)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateConnection

	if resp.Diagnostics.Append(reqUpdate.set(ctx, plan, config)...); resp.Diagnostics.HasError() {
		return
	}

	respUpdate, err := r.client.UpdateConnection(ctx, plan.ID.ValueString(), reqUpdate.UpdateConnectionRequestClassification, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, respUpdate.ConnectionClassification)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceConnection) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	var state resourceConnectionModel[rsConnectionDetailsModel, rsCredentialDetailsModel]

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.DeleteConnection(ctx, state.ID.ValueString(), nil)
	resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...)

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceConnection) get(ctx context.Context, model *resourceConnectionModel[rsConnectionDetailsModel, rsCredentialDetailsModel]) diag.Diagnostics {
	tflog.Trace(ctx, "GET", map[string]any{
		"id": model.ID.ValueString(),
	})

	respGet, err := r.client.GetConnection(ctx, model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrWorkspace.WorkspaceNotFound); diags.HasError() {
		return diags
	}

	return model.set(ctx, respGet.ConnectionClassification)
}

func (r *resourceConnection) getConnectionTypeMetadata(ctx context.Context, model rsConnectionDetailsModel, supportedConnectionType *fabcore.ConnectionCreationMetadata) diag.Diagnostics {
	pager := r.client.NewListSupportedConnectionTypesPager(&fabcore.ConnectionsClientListSupportedConnectionTypesOptions{
		ShowAllCreationMethods: azto.Ptr(true),
	})

	var allConnections []fabcore.ConnectionCreationMetadata

	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		allConnections = append(allConnections, resp.Value...)
	}

	vNames := make([]string, 0, len(allConnections))

	for _, v := range allConnections {
		if *v.Type == model.Type.ValueString() {
			*supportedConnectionType = v

			return nil
		}

		vNames = append(vNames, *v.Type)
	}

	var diags diag.Diagnostics

	diags.AddAttributeError(
		path.Root("connection_details").AtName("type"),
		"Unsupported connection type",
		fmt.Sprintf("The connection type '%s' is not supported. Supported values: %s", model.Type.ValueString(), utils.ConvertStringSlicesToString(vNames, true, true)),
	)

	return diags
}

func (r *resourceConnection) validateCreationMethod(model rsConnectionDetailsModel, elements []fabcore.ConnectionCreationMethod) diag.Diagnostics {
	vNames := make([]string, 0, len(elements))

	for _, v := range elements {
		if *v.Name == model.CreationMethod.ValueString() {
			return nil
		}

		vNames = append(vNames, *v.Name)
	}

	var diags diag.Diagnostics

	diags.AddAttributeError(
		path.Root("connection_details").AtName("creation_method"),
		"Unsupported creation method",
		fmt.Sprintf("The creation method '%s' is not supported. Supported values: %s", model.CreationMethod.ValueString(), utils.ConvertStringSlicesToString(vNames, true, true)),
	)

	return diags
}

func (r *resourceConnection) validateConnectionEncryption(model rsCredentialDetailsModel, elements []fabcore.ConnectionEncryption) diag.Diagnostics {
	for _, v := range elements {
		if v == fabcore.ConnectionEncryption(model.ConnectionEncryption.ValueString()) {
			return nil
		}
	}

	var diags diag.Diagnostics

	diags.AddAttributeError(
		path.Root("credential_details").AtName("connection_encryption"),
		"Unsupported connection encryption",
		fmt.Sprintf(
			"The connection encryption '%s' is not supported. Supported values: %s",
			model.ConnectionEncryption.ValueString(),
			utils.ConvertStringSlicesToString(utils.ConvertEnumsToStringSlices(elements, true), true, false),
		),
	)

	return diags
}

func (r *resourceConnection) validateCredentialType(model rsCredentialDetailsModel, elements []fabcore.CredentialType) diag.Diagnostics {
	if slices.Contains(elements, fabcore.CredentialType(model.CredentialType.ValueString())) {
		return nil
	}

	var diags diag.Diagnostics

	diags.AddAttributeError(
		path.Root("credential_details").AtName("credential_type"),
		"Unsupported credential type",
		fmt.Sprintf(
			"The credential type '%s' is not supported. Supported values: %s",
			model.CredentialType.ValueString(),
			utils.ConvertStringSlicesToString(utils.ConvertEnumsToStringSlices(elements, true), true, false),
		),
	)

	return diags
}

func (r *resourceConnection) validateSkipTestConnection(model rsCredentialDetailsModel, supportsSkipTestConnection bool) diag.Diagnostics { //revive:disable-line:flag-parameter
	if model.SkipTestConnection.ValueBool() != supportsSkipTestConnection {
		var diags diag.Diagnostics

		diags.AddAttributeError(
			path.Root("credential_details").AtName("skip_test_connection"),
			"Unsupported skip test connection",
			"The skip test connection value is not supported.",
		)

		return diags
	}

	return nil
}

//nolint:gocognit,gocyclo
func (r *resourceConnection) validateCreationMethodParameters(ctx context.Context, model rsConnectionDetailsModel, elements []fabcore.ConnectionCreationMethod) diag.Diagnostics {
	var diags diag.Diagnostics
	var vParameters []fabcore.ConnectionCreationParameter

	for _, v := range elements {
		if *v.Name == model.CreationMethod.ValueString() {
			vParameters = v.Parameters

			break
		}
	}

	// in reality, this should never happen because the creation method is already validated
	if len(vParameters) == 0 {
		diags.AddAttributeError(
			path.Root("connection_details").AtName("creation_method"),
			"Unsupported creation method",
			"The creation method is not supported.",
		)

		return diags
	}

	connectionDetailsParams, diags := model.getParameters(ctx)
	if diags.HasError() {
		return diags
	}

	// check if all keys of connectionDetailsParams are in vParameters
	var vNames []string

	for k := range connectionDetailsParams {
		var found bool

		for _, v := range vParameters {
			if *v.Name == k {
				found = true

				break
			}

			vNames = append(vNames, *v.Name)
		}

		if !found {
			diags.AddAttributeError(
				path.Root("connection_details").AtName("parameters"),
				"Unsupported connection parameter key",
				fmt.Sprintf("The connection parameter '%s' is not supported. Supported parameters: %s", k, utils.ConvertStringSlicesToString(vNames, true, true)),
			)
		}
	}

	if diags.HasError() {
		return diags
	}

	// check if all required keys of vParameters are in connectionDetailsParams
	for _, v := range vParameters {
		if *v.Required {
			if _, ok := connectionDetailsParams[*v.Name]; !ok {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Missing connection parameter key",
					fmt.Sprintf("The required connection parameter '%s' is missing.", *v.Name),
				)
			}

			if connectionDetailsParams[*v.Name] == "" {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Missing connection parameter value",
					fmt.Sprintf("The required connection parameter '%s' value is missing.", *v.Name),
				)
			}
		}
	}

	if diags.HasError() {
		return diags
	}

	for k, v := range connectionDetailsParams {
		var dataType fabcore.DataType
		var allowedValues []string

		for _, v := range vParameters {
			if *v.Name == k {
				dataType = *v.DataType
				allowedValues = v.AllowedValues

				break
			}
		}

		switch dataType {
		case fabcore.DataTypeBoolean:
			// Use boolean as the parameter input value. False - the value is false, True - the value is true.
			if !strings.EqualFold(v, "true") && !strings.EqualFold(v, "false") {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. Supported values: `True`, `False`", k),
				)
			}
		case fabcore.DataTypeDate:
			// Use date as the parameter input value, using YYYY-MM-DD format.
			if _, err := time.Parse(time.DateOnly, v); err != nil {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. Supported format: `YYYY-MM-DD`", k),
				)
			}
		case fabcore.DataTypeDateTime:
			// Use date time as the parameter input value, using YYYY-MM-DDTHH:mm:ss.FFFZ (ISO 8601) format.
			if _, err := time.Parse("2006-01-02T15:04:05.000Z07:00", v); err != nil {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. Supported format: `YYYY-MM-DDTHH:mm:ss.FFFZ`", k),
				)
			}
		case fabcore.DataTypeDateTimeZone:
			// Use date time zone as the parameter input value, using YYYY-MM-DDTHH:mm:ss.FFF±hh:mm format.
			if _, err := time.Parse("2006-01-02T03:04:05.000-07:00", v); err != nil {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. Supported format: `YYYY-MM-DDTHH:mm:ss.FFF±hh:mm`", k),
				)
			}
		case fabcore.DataTypeDuration:
			// Use duration as the parameter input value, using [-]P(n)DT(n)H(n)M(n)S format. For example: P3DT4H30M10S (for 3 days, 4 hours, 30 minutes, and 10 seconds).
			if _, err := time.ParseDuration(v); err != nil {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf(
						"The connection parameter '%s' value is invalid. Supported format: `[-]P(n)DT(n)H(n)M(n)S`. For example: `P3DT4H30M10S` (for 3 days, 4 hours, 30 minutes, and 10 seconds).",
						k,
					),
				)
			}
		case fabcore.DataTypeNumber:
			// Use number as the parameter input value (integer or floating point).
			if _, err := strconv.ParseFloat(v, 32); err != nil {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. It must be integer or floating point.", k),
				)
			}
		case fabcore.DataTypeText:
			// Use text as the parameter input value.
			if v == "" {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. It must not be empty.", k),
				)
			}
		case fabcore.DataTypeTime:
			// Use time as the parameter input value, using HH:mm:ss.FFFZ format.
			if _, err := time.Parse("15:04:05.000Z07:00", v); err != nil {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. Supported format: `HH:mm:ss.FFFZ`", k),
				)
			}
		}

		if len(allowedValues) > 0 {
			if !slices.Contains(allowedValues, v) {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. Supported values: %s", k, utils.ConvertStringSlicesToString(allowedValues, true, true)),
				)
			}
		}
	}

	return diags
}

func (r *resourceConnection) setConnectionPerametersDataType(
	ctx context.Context,
	elements []fabcore.ConnectionCreationMethod,
	plan *resourceConnectionModel[rsConnectionDetailsModel, rsCredentialDetailsModel],
) diag.Diagnostics {
	var diags diag.Diagnostics
	var vParameters []fabcore.ConnectionCreationParameter

	connectionDetails, diags := plan.ConnectionDetails.Get(ctx)
	if diags.HasError() {
		return diags
	}

	for _, v := range elements {
		if *v.Name == connectionDetails.CreationMethod.ValueString() {
			vParameters = v.Parameters

			break
		}
	}

	parameters, diags := connectionDetails.Parameters.Get(ctx)
	if diags.HasError() {
		return diags
	}

	for _, parameter := range parameters {
		for _, v := range vParameters {
			if strings.EqualFold(*v.Name, parameter.Name.ValueString()) {
				parameter.DataType = types.StringValue(string(*v.DataType))

				break
			}
		}
	}

	if diags := connectionDetails.Parameters.Set(ctx, parameters); diags.HasError() {
		return diags
	}

	if diags := plan.ConnectionDetails.Set(ctx, connectionDetails); diags.HasError() {
		return diags
	}

	return nil
}
