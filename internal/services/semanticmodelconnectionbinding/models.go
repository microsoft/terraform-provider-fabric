// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package semanticmodelconnectionbinding

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabsemanticmodel "github.com/microsoft/fabric-sdk-go/fabric/semanticmodel"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type connectionDetailsModel struct {
	Path types.String `tfsdk:"path"`
	Type types.String `tfsdk:"type"`
}

type connectionBindingPropertiesModel struct {
	ID                types.String                                                 `tfsdk:"id"`
	ConnectionID      customtypes.UUID                                             `tfsdk:"connection_id"`
	ConnectivityType  types.String                                                 `tfsdk:"connectivity_type"`
	ConnectionDetails supertypes.SingleNestedObjectValueOf[connectionDetailsModel] `tfsdk:"connection_details"`
	DisplayName       types.String                                                 `tfsdk:"display_name"`
	GatewayID         customtypes.UUID                                             `tfsdk:"gateway_id"`
}

func (to *connectionBindingPropertiesModel) set(ctx context.Context, workspaceID, semanticModelID string, from fabcore.ItemConnection) diag.Diagnostics {
	if from.ConnectionDetails == nil || from.ConnectionDetails.Path == nil || from.ConnectionDetails.Type == nil || from.ConnectivityType == nil {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic("Invalid connection binding response", "The API response is missing required connection details or connectivity type."),
		}
	}

	to.ID = types.StringValue(connectionBindingIdentifier(workspaceID, semanticModelID, *from.ConnectionDetails.Path, *from.ConnectionDetails.Type))
	to.ConnectionID = customtypes.NewUUIDPointerValue(from.ID)
	to.ConnectivityType = types.StringPointerValue((*string)(from.ConnectivityType))
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.GatewayID = customtypes.NewUUIDPointerValue(from.GatewayID)

	details := connectionDetailsModel{
		Path: types.StringPointerValue(from.ConnectionDetails.Path),
		Type: types.StringPointerValue(from.ConnectionDetails.Type),
	}
	to.ConnectionDetails = supertypes.NewSingleNestedObjectValueOfNull[connectionDetailsModel](ctx)

	return to.ConnectionDetails.Set(ctx, &details)
}

type resourceSemanticModelConnectionBindingModel struct {
	connectionBindingPropertiesModel

	WorkspaceID     customtypes.UUID `tfsdk:"workspace_id"`
	SemanticModelID customtypes.UUID `tfsdk:"semantic_model_id"`
	Timeouts        timeoutsR.Value  `tfsdk:"timeouts"`
}

type dataSourceSemanticModelConnectionBindingModel struct {
	connectionBindingPropertiesModel

	WorkspaceID     customtypes.UUID `tfsdk:"workspace_id"`
	SemanticModelID customtypes.UUID `tfsdk:"semantic_model_id"`
	Timeouts        timeoutsD.Value  `tfsdk:"timeouts"`
}

type dataSourceSemanticModelConnectionBindingsModel struct {
	WorkspaceID     customtypes.UUID                                                    `tfsdk:"workspace_id"`
	SemanticModelID customtypes.UUID                                                    `tfsdk:"semantic_model_id"`
	Values          supertypes.SetNestedObjectValueOf[connectionBindingPropertiesModel] `tfsdk:"values"`
	Timeouts        timeoutsD.Value                                                     `tfsdk:"timeouts"`
}

func (to *dataSourceSemanticModelConnectionBindingsModel) setValues(ctx context.Context, from []fabcore.ItemConnection) diag.Diagnostics {
	values := make([]*connectionBindingPropertiesModel, 0, len(from))

	for _, binding := range from {
		var bindingModel connectionBindingPropertiesModel
		if diags := bindingModel.set(ctx, to.WorkspaceID.ValueString(), to.SemanticModelID.ValueString(), binding); diags.HasError() {
			return diags
		}

		values = append(values, &bindingModel)
	}

	return to.Values.Set(ctx, values)
}

type requestBindSemanticModelConnection struct {
	fabsemanticmodel.BindSemanticModelConnectionRequest
}

func (to *requestBindSemanticModelConnection) set(ctx context.Context, from resourceSemanticModelConnectionBindingModel, connectivityType fabsemanticmodel.ConnectivityType) diag.Diagnostics {
	details, diags := from.ConnectionDetails.Get(ctx)
	if diags.HasError() {
		return diags
	}

	var connectionID *string
	if connectivityType != fabsemanticmodel.ConnectivityTypeNone && connectivityType != fabsemanticmodel.ConnectivityTypeAutomatic {
		connectionID = from.ConnectionID.ValueStringPointer()
	}

	to.ConnectionBinding = &fabsemanticmodel.ConnectionBinding{
		ConnectionDetails: &fabsemanticmodel.ListConnectionDetails{
			Path: details.Path.ValueStringPointer(),
			Type: details.Type.ValueStringPointer(),
		},
		ConnectivityType: &connectivityType,
		ID:               connectionID,
	}

	return nil
}

func connectionBindingIdentifier(workspaceID, semanticModelID, path, connectionType string) string {
	encodedPath := base64.RawURLEncoding.EncodeToString([]byte(path))
	encodedType := base64.RawURLEncoding.EncodeToString([]byte(connectionType))

	return strings.Join([]string{workspaceID, semanticModelID, encodedPath, encodedType}, "/")
}

type connectionBindingImportIdentifier struct {
	WorkspaceID     string
	SemanticModelID string
	ConnectionPath  string
	ConnectionType  string
}

func parseConnectionBindingIdentifier(identifier string) (connectionBindingImportIdentifier, error) {
	parts := strings.Split(identifier, "/")
	if len(parts) != 4 {
		return connectionBindingImportIdentifier{}, errors.New("expected four slash-separated identifier segments")
	}

	pathBytes, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return connectionBindingImportIdentifier{}, fmt.Errorf("decode connection path: %w", err)
	}

	typeBytes, err := base64.RawURLEncoding.DecodeString(parts[3])
	if err != nil {
		return connectionBindingImportIdentifier{}, fmt.Errorf("decode connection type: %w", err)
	}

	if base64.RawURLEncoding.EncodeToString(pathBytes) != parts[2] || base64.RawURLEncoding.EncodeToString(typeBytes) != parts[3] {
		return connectionBindingImportIdentifier{}, errors.New("identifier contains non-canonical base64url values")
	}

	return connectionBindingImportIdentifier{
		WorkspaceID:     parts[0],
		SemanticModelID: parts[1],
		ConnectionPath:  string(pathBytes),
		ConnectionType:  string(typeBytes),
	}, nil
}
