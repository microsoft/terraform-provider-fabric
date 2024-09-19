// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package testhelp

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/microsoft/terraform-provider-fabric/internal/auth"
)

func CheckAuthMethod(expected auth.AuthenticationMethod, testState *TestState) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		resp, err := auth.NewCredential(*testState.Config.Auth)
		if err != nil {
			return err
		}

		if resp.AuthMethod != expected {
			return fmt.Errorf("expected %s got %s", expected, resp.AuthMethod)
		}

		return nil
	}
}

func CheckAuthConfig(property, value string, testState *TestState) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		field := getConfigField(testState.Config.Auth, property)
		if field != value {
			return fmt.Errorf("%s not set correctly, expected %s got %s", property, value, field)
		}

		return nil
	}
}

func getConfigField(v *auth.Config, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)

	return f.String()
}
