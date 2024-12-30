// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

type resourceConnectionModel struct {
	baseResourceConnectionModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

func (m resourceConnectionModel) getConnectionDetails(ctx context.Context) (*rsConnectionDetailsModel, diag.Diagnostics) {
	if !m.ConnectionDetails.IsNull() && !m.ConnectionDetails.IsUnknown() {
		return m.ConnectionDetails.Get(ctx)
	}

	return nil, nil
}

func (m resourceConnectionModel) getCredentialDetails(ctx context.Context) (*rsCredentialDetailsModel, diag.Diagnostics) {
	if !m.CredentialDetails.IsNull() && !m.CredentialDetails.IsUnknown() {
		return m.CredentialDetails.Get(ctx)
	}

	return nil, nil
}

type requestCreateConnection struct {
	fabcore.CreateConnectionRequestClassification
}

func (to *requestCreateConnection) set(ctx context.Context, from resourceConnectionModel, supportedConnectionType fabcore.ConnectionCreationMetadata) diag.Diagnostics {
	// to.DisplayName = from.DisplayName.ValueStringPointer()
	// to.PrivacyLevel = (*fabcore.PrivacyLevel)(from.PrivacyLevel.ValueStringPointer())
	connectivityType := (fabcore.ConnectivityType)(from.ConnectivityType.ValueString())

	// connectionDetails, diags := from.ConnectionDetails.Get(ctx)
	// if diags.HasError() {
	// 	return diags
	// }

	// connectionDetailsParameters, diags := connectionDetails.Parameters.Get(ctx)
	// if diags.HasError() {
	// 	return diags
	// }

	// credentialDetails, diags := from.CredentialDetails.Get(ctx)
	// if diags.HasError() {
	// 	return diags
	// }

	var requestCreateConnectionDetails requestCreateConnectionDetails
	requestCreateConnectionDetails.set(ctx, from.ConnectionDetails)

	switch connectivityType {
	case fabcore.ConnectivityTypeShareableCloud:
		aaa := &fabcore.CreateCloudConnectionRequest{
			DisplayName:       from.DisplayName.ValueStringPointer(),
			PrivacyLevel:      (*fabcore.PrivacyLevel)(from.PrivacyLevel.ValueStringPointer()),
			ConnectivityType:  &connectivityType,
			ConnectionDetails: &requestCreateConnectionDetails.CreateConnectionDetails,
			CredentialDetails: &fabcore.CreateCredentialDetails{},
		}

		to.CreateConnectionRequestClassification = aaa
	case fabcore.ConnectivityTypeVirtualNetworkGateway:

		bbb := &fabcore.CreateVirtualNetworkGatewayConnectionRequest{}

		to.CreateConnectionRequestClassification = bbb
	}

	return nil
}

type requestCreateConnectionDetails struct {
	fabcore.CreateConnectionDetails
}

func (to *requestCreateConnectionDetails) set(ctx context.Context, from supertypes.SingleNestedObjectValueOf[rsConnectionDetailsModel]) diag.Diagnostics {
	if !from.IsNull() && !from.IsUnknown() {
		connectionDetails, diags := from.Get(ctx)
		if diags.HasError() {
			return diags
		}

		if !connectionDetails.Parameters.IsNull() && !connectionDetails.Parameters.IsUnknown() {
			_, diags := connectionDetails.Parameters.Get(ctx)
			if diags.HasError() {
				return diags
			}
		}

		to.CreationMethod = connectionDetails.CreationMethod.ValueStringPointer()
		to.Type = connectionDetails.Type.ValueStringPointer()
	}

	return nil
}
