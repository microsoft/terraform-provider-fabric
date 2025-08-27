// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (o *OIDCConfig) getAssertion(ctx context.Context) (string, error) {
	if o.Token != "" {
		return o.Token, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, o.RequestURL, http.NoBody)
	if err != nil {
		return "", errors.New("getAssertion: failed to build request")
	}

	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return "", errors.New("getAssertion: cannot parse URL query")
	}

	if query.Get("audience") == "" {
		query.Set("audience", "api://AzureADTokenExchange")
		req.URL.RawQuery = query.Encode()
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.RequestToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("getAssertion: cannot request token: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("getAssertion: cannot parse response: %w", err)
	}

	if c := resp.StatusCode; c < http.StatusOK || c > http.StatusIMUsed {
		return "", fmt.Errorf("getAssertion: received HTTP status %d with response: %s", resp.StatusCode, body)
	}

	var tokenResp struct {
		Count *int    `json:"count"`
		Value *string `json:"value"`
	}

	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return "", fmt.Errorf("getAssertion: cannot unmarshal response: %w", err)
	}

	if tokenResp.Value == nil {
		return "", errors.New("getAssertion: nil JWT assertion received from OIDC provider")
	}

	return *tokenResp.Value, nil
}
