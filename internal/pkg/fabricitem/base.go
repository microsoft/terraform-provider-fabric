// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

const (
	DefinitionFormatDefault          = "Default"
	SPNSupportedDataSource           = "\n\n-> This data-source supports Service Principal authentication."
	SPNNotSupportedDataSource        = "\n\n-> This data-source does not support Service Principal. Please use a User context authentication."
	SPNSupportedResource             = "\n\n-> This resource supports Service Principal authentication."
	SPNNotSupportedResource          = "\n\n-> This resource does not support Service Principal. Please use a User context authentication."
	SPNSupportedEphemeralResource    = "\n\n-> This ephemeral resource supports Service Principal authentication."
	SPNNotSupportedEphemeralResource = "\n\n-> This ephemeral resource does not support Service Principal. Please use a User context authentication."
	PreviewDataSource                = "\n\n~> This data-source is in **preview**. To access it, you must explicitly enable the `preview` mode in the provider level configuration."
	PreviewResource                  = "\n\n~> This resource is in **preview**. To access it, you must explicitly enable the `preview` mode in the provider level configuration."
	PreviewEphemeralResource         = "\n\n~> This ephemeral resource is in **preview**. To access it, you must explicitly enable the `preview` mode in the provider level configuration."
)
