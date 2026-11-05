// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package client_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	"github.com/stretchr/testify/require"

	pclient "github.com/microsoft/terraform-provider-fabric/internal/provider/client"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

type cloudCredential struct {
	scopes []string
}

func (c *cloudCredential) GetToken(_ context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.scopes = options.Scopes

	return azcore.AccessToken{Token: "test-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

type cloudTransport struct {
	request *http.Request
}

func (c *cloudTransport) Do(req *http.Request) (*http.Response, error) {
	c.request = req

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"value":[]}`)),
		Request:    req,
	}, nil
}

func TestUnit_FabricClient_Environment(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		environment   fabric.Environment
		endpoint      string
		privateLinks  bool
		expectedHost  string
		expectedScope string
	}{
		{
			name:          "Public defaults",
			expectedHost:  "api.fabric.microsoft.com",
			expectedScope: "https://api.fabric.microsoft.com/.default",
		},
		{
			name:          "GCC High",
			environment:   fabric.EnvironmentGCCHigh,
			expectedHost:  "highapi.fabric.microsoft.us",
			expectedScope: "https://api.fabric.microsoft.us/.default",
		},
		{
			name:          "Public endpoint override",
			endpoint:      "https://override.example.com",
			expectedHost:  "override.example.com",
			expectedScope: "https://api.fabric.microsoft.com/.default",
		},
		{
			name:          "GCC High endpoint override",
			environment:   fabric.EnvironmentGCCHigh,
			endpoint:      "https://override.example.com",
			expectedHost:  "override.example.com",
			expectedScope: "https://api.fabric.microsoft.us/.default",
		},
		{
			name:          "Public workspace private links",
			privateLinks:  true,
			expectedHost:  "12345678123412341234123456789abc.z12.w.api.fabric.microsoft.com",
			expectedScope: "https://api.fabric.microsoft.com/.default",
		},
		{
			name:          "GCC High workspace private links",
			environment:   fabric.EnvironmentGCCHigh,
			privateLinks:  true,
			expectedHost:  "12345678123412341234123456789abc.z12.w.highapi.fabric.microsoft.us",
			expectedScope: "https://api.fabric.microsoft.us/.default",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			cred := &cloudCredential{}
			transport := &cloudTransport{}
			options := &fabric.ClientOptions{Transport: transport}
			cfg := &pconfig.ProviderConfig{ProviderData: &pconfig.ProviderData{
				Environment:                     testCase.environment,
				Endpoint:                        testCase.endpoint,
				UseWorkspacePrivateLinkEndpoint: testCase.privateLinks,
			}}

			client, err := pclient.NewFabricClient(cred, cfg, options)
			require.NoError(t, err)
			require.Equal(t, testCase.privateLinks, client.UseWorkspacePrivateLinks)

			pager := fabcore.NewClientFactoryWithClient(*client).NewItemsClient().NewListItemsPager("12345678-1234-1234-1234-123456789abc", nil)
			_, err = pager.NextPage(t.Context())
			require.NoError(t, err)
			require.NotNil(t, transport.request)
			require.Equal(t, testCase.expectedHost, transport.request.URL.Host)
			require.Equal(t, []string{testCase.expectedScope}, cred.scopes)
			require.Equal(t, "Bearer test-token", transport.request.Header.Get("Authorization"))
			require.Equal(t, fabric.EnvironmentPublic, options.Environment, "caller options must not be mutated")
		})
	}
}

func TestUnit_FabricClient_InvalidConfiguration(t *testing.T) {
	t.Parallel()

	for name, data := range map[string]pconfig.ProviderData{
		"unsupported environment": {Environment: fabric.Environment(99)},
		"relative endpoint":       {Endpoint: "invalid"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client, err := pclient.NewFabricClient(&cloudCredential{}, &pconfig.ProviderConfig{ProviderData: &data}, nil)
			require.Error(t, err)
			require.Nil(t, client)
		})
	}
}
