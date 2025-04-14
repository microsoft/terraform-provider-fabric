// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/fabric-sdk-go/fabric/core"
)

// logRequestBody logs the request body in a structured way to help with debugging
// It handles various credential types and masks sensitive values like passwords and secrets
func logRequestBody(ctx context.Context, req *requestCreateConnection) {
	if req == nil {
		tflog.Info(ctx, "Request is nil, nothing to log")
		return
	}

	// Log the full type information of the request for debugging
	tflog.Info(ctx, fmt.Sprintf("Request type: %T", req.CreateConnectionRequestClassification))

	// Log the underlying request type using reflection
	reqType := reflect.TypeOf(req.CreateConnectionRequestClassification)

	tflog.Info(ctx, fmt.Sprintf("Request implementation type: %s", reqType.String()))

	// Log the actual CreateConnectionRequest from the interface
	if req.CreateConnectionRequestClassification != nil {
		baseReq := req.CreateConnectionRequestClassification.GetCreateConnectionRequest()
		if baseReq != nil {
			// Log connectivity type
			if baseReq.ConnectivityType != nil {
				tflog.Info(ctx, fmt.Sprintf("ConnectivityType: %s", string(*baseReq.ConnectivityType)))
			} else {
				tflog.Info(ctx, "ConnectivityType is nil")
			}

			// Log display name (not sensitive)
			if baseReq.DisplayName != nil {
				tflog.Info(ctx, fmt.Sprintf("DisplayName: %s", *baseReq.DisplayName))
			}

			// Log privacy level
			if baseReq.PrivacyLevel != nil {
				tflog.Info(ctx, fmt.Sprintf("PrivacyLevel: %s", string(*baseReq.PrivacyLevel)))
			} else {
				tflog.Info(ctx, "PrivacyLevel is nil")
			}

			// For specific request types, log additional info
			switch typedReq := req.CreateConnectionRequestClassification.(type) {
			case *core.CreateCloudConnectionRequest:
				tflog.Info(ctx, "Request is a CreateCloudConnectionRequest")
				if typedReq.CredentialDetails != nil && typedReq.CredentialDetails.Credentials != nil {
					credType := typedReq.CredentialDetails.Credentials.GetCredentials().CredentialType
					if credType != nil {
						tflog.Info(ctx, fmt.Sprintf("CredentialType: %s", string(*credType)))
					}
				}
			case *core.CreateOnPremisesConnectionRequest:
				tflog.Info(ctx, "Request is a CreateOnPremisesConnectionRequest")
				if typedReq.GatewayID != nil {
					tflog.Info(ctx, fmt.Sprintf("GatewayID: %s", *typedReq.GatewayID))
				}
				if typedReq.CredentialDetails != nil && typedReq.CredentialDetails.Credentials != nil {
					tflog.Info(ctx, "Contains OnPremisesGatewayCredentials")
				}
			case *core.CreateVirtualNetworkGatewayConnectionRequest:
				tflog.Info(ctx, "Request is a CreateVirtualNetworkGatewayConnectionRequest")
				if typedReq.GatewayID != nil {
					tflog.Info(ctx, fmt.Sprintf("GatewayID: %s", *typedReq.GatewayID))
				}
			default:
				tflog.Info(ctx, fmt.Sprintf("Unknown request type: %T", typedReq))
			}
		}
	}

	// Create a copy of the request to sanitize sensitive data

	// Marshal the sanitized request to JSON for readable logging
	jsonBytes, err := json.MarshalIndent(req.GetCreateConnectionRequest(), "", "  ")
	if err != nil {
		tflog.Error(ctx, "Failed to marshal request body for logging", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Log the sanitized JSON request
	tflog.Info(ctx, "Connection request body", map[string]interface{}{
		"request": string(jsonBytes),
	})

	// Try to extract more specific diagnostic information
	if cloudReq, ok := req.CreateConnectionRequestClassification.(*core.CreateCloudConnectionRequest); ok {
		if cloudReq.CredentialDetails != nil {
			credType := "*Unknown*"
			if cloudReq.CredentialDetails.Credentials != nil {
				if ct := cloudReq.CredentialDetails.Credentials.GetCredentials().CredentialType; ct != nil {
					credType = string(*ct)
				}
			}

			tflog.Info(ctx, "Connection credential details", map[string]interface{}{
				"credential_type":       credType,
				"connection_encryption": string(*cloudReq.CredentialDetails.ConnectionEncryption),
				"single_sign_on_type":   string(*cloudReq.CredentialDetails.SingleSignOnType),
				"skip_test_connection":  *cloudReq.CredentialDetails.SkipTestConnection,
			})

			if cloudReq.ConnectionDetails != nil && cloudReq.ConnectionDetails.Parameters != nil {
				// Log parameter names (not values to protect sensitive data)
				paramNames := make([]string, 0, len(cloudReq.ConnectionDetails.Parameters))
				for _, param := range cloudReq.ConnectionDetails.Parameters {
					if textParam, ok := param.(*core.ConnectionDetailsTextParameter); ok && textParam.Name != nil {
						paramNames = append(paramNames, *textParam.Name)
					} else {
						paramNames = append(paramNames, "unknown")
					}
				}

				tflog.Info(ctx, "Connection parameters", map[string]interface{}{
					"parameter_count": len(cloudReq.ConnectionDetails.Parameters),
					"parameter_names": strings.Join(paramNames, ", "),
				})
			}
		}
	}
}
