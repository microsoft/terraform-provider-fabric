// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package environment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"
)

type resourceEnvironmentModel struct {
	baseEnvironmentPropertiesModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type requestCreateEnvironment struct {
	fabenvironment.CreateEnvironmentRequest
}

func (to *requestCreateEnvironment) set(from resourceEnvironmentModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
}

type requestUpdateEnvironment struct {
	fabenvironment.UpdateEnvironmentRequest
}

func (to *requestUpdateEnvironment) set(from resourceEnvironmentModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
}
