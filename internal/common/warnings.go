// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package common //revive:disable-line:package-naming

const (
	WarningCapacityNotActiveHeader     = "Fabric Capacity is not active"
	WarningCapacityNotActiveDetails    = "Could not refresh Fabric Item %s because its capacity is not active. Keeping the previous Terraform state; remote changes have not been checked. Resume the capacity and run terraform plan again to refresh the item. Item create or update operations may still fail while the capacity is inactive."
	WarningItemDefinitionUpdateHeader  = "Fabric Item definition update"
	WarningItemDefinitionUpdateDetails = "%s definition update operation will overwrite the existing definition on the Fabric side."
	WarningPreviewModeHeader           = "'%s' preview mode"
	WarningPreviewModeDetails          = "The '%s' behavior may change in future releases without notice or compatibility guarantees."
)
