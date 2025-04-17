// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const HeaderUserAgent = "User-Agent"

type UserAgentPolicy struct {
	UserAgent string
}

func (c UserAgentPolicy) Do(req *policy.Request) (*http.Response, error) {
	req.Raw().Header.Set(HeaderUserAgent, c.UserAgent)
	return req.Next()
}

var _ policy.Policy = UserAgentPolicy{}

// WithUserAgent returns a policy.Policy that adds an HTTP extension header of
// `User-Agent` whose value is passed and has no length limitation
func WithUserAgent(userAgent string) policy.Policy {
	return UserAgentPolicy{UserAgent: userAgent}
}

func BuildUserAgent(terraformVersion, sdkVersion, providerVersion, partnerID string) string {
	if terraformVersion == "" {
		terraformVersion = "0.11+compatible"
	}

	terraformUserAgent := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io)", terraformVersion)
	sdkUserAgent := fmt.Sprintf("fabric-sdk-go/%s", sdkVersion)
	providerUserAgent := fmt.Sprintf("terraform-provider-fabric/%s", providerVersion)
	userAgent := strings.TrimSpace(fmt.Sprintf("%s %s %s", terraformUserAgent, sdkUserAgent, providerUserAgent))

	if partnerID == "" {
		// Microsoft’s Terraform Partner ID is this specific GUID
		partnerID = "222c6c49-1b0a-5959-a213-6608f9eb8820"
	}

	userAgent = fmt.Sprintf("%s pid-%s", userAgent, partnerID)

	return userAgent
}
