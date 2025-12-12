// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package tenantsettings

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseTenantSettingsModel struct {
	SettingName              types.String                                                   `tfsdk:"setting_name"`
	TenantSettingGroup       types.String                                                   `tfsdk:"tenant_setting_group"`
	Title                    types.String                                                   `tfsdk:"title"`
	CanSpecifySecurityGroups types.Bool                                                     `tfsdk:"can_specify_security_groups"`
	Enabled                  types.Bool                                                     `tfsdk:"enabled"`
	DelegateToCapacity       types.Bool                                                     `tfsdk:"delegate_to_capacity"`
	DelegateToDomain         types.Bool                                                     `tfsdk:"delegate_to_domain"`
	DelegateToWorkspace      types.Bool                                                     `tfsdk:"delegate_to_workspace"`
	DeleteBehaviour          types.String                                                   `tfsdk:"delete_behaviour"`
	EnabledSecurityGroups    supertypes.SetNestedObjectValueOf[tenantSettingsSecurityGroup] `tfsdk:"enabled_security_groups"`
	ExcludedSecurityGroups   supertypes.SetNestedObjectValueOf[tenantSettingsSecurityGroup] `tfsdk:"excluded_security_groups"`
	Properties               supertypes.SetNestedObjectValueOf[tenantSettingsProperty]      `tfsdk:"properties"`
}

func (to *baseTenantSettingsModel) set(ctx context.Context, from fabadmin.TenantSetting) diag.Diagnostics {
	to.SettingName = types.StringPointerValue(from.SettingName)
	to.TenantSettingGroup = types.StringPointerValue(from.TenantSettingGroup)
	to.Title = types.StringPointerValue(from.Title)
	to.CanSpecifySecurityGroups = types.BoolPointerValue(from.CanSpecifySecurityGroups)
	to.Enabled = types.BoolPointerValue(from.Enabled)
	to.DelegateToCapacity = types.BoolPointerValue(from.DelegateToCapacity)
	to.DelegateToDomain = types.BoolPointerValue(from.DelegateToDomain)
	to.DelegateToWorkspace = types.BoolPointerValue(from.DelegateToWorkspace)

	if from.EnabledSecurityGroups != nil {
		slice := make([]*tenantSettingsSecurityGroup, 0, len(from.EnabledSecurityGroups))
		for _, securityGroup := range from.EnabledSecurityGroups {
			sg := tenantSettingsSecurityGroup{}
			sg.set(securityGroup)
			slice = append(slice, &sg)
		}
		if diags := to.EnabledSecurityGroups.Set(ctx, slice); diags.HasError() {
			return diags
		}
	} else {
		if diags := to.EnabledSecurityGroups.Set(ctx, []*tenantSettingsSecurityGroup{}); diags.HasError() {
			return diags
		}
	}

	if from.ExcludedSecurityGroups != nil {
		slice := make([]*tenantSettingsSecurityGroup, 0, len(from.ExcludedSecurityGroups))
		for _, securityGroup := range from.ExcludedSecurityGroups {
			sg := tenantSettingsSecurityGroup{}
			sg.set(securityGroup)
			slice = append(slice, &sg)
		}
		if diags := to.ExcludedSecurityGroups.Set(ctx, slice); diags.HasError() {
			return diags
		}
	} else {
		if diags := to.ExcludedSecurityGroups.Set(ctx, []*tenantSettingsSecurityGroup{}); diags.HasError() {
			return diags
		}
	}

	if from.Properties != nil {
		slice := make([]*tenantSettingsProperty, 0, len(from.Properties))
		for _, property := range from.Properties {
			prop := tenantSettingsProperty{}
			prop.set(property)
			slice = append(slice, &prop)
		}
		if diags := to.Properties.Set(ctx, slice); diags.HasError() {
			return diags
		}
	} else {
		if diags := to.Properties.Set(ctx, []*tenantSettingsProperty{}); diags.HasError() {
			return diags
		}
	}

	return nil
}

/*
DATA-SOURCE
*/

type dataSourceTenantSettingModel struct {
	baseTenantSettingsModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (LIST)
*/

type dataSourceTenantSettingsModel struct {
	Values supertypes.SetNestedObjectValueOf[baseTenantSettingsModel] `tfsdk:"values"`

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

func (to *dataSourceTenantSettingsModel) set(ctx context.Context, from []fabadmin.TenantSetting) diag.Diagnostics {
	slice := make([]*baseTenantSettingsModel, 0, len(from))
	for _, ts := range from {
		tenantSetting := baseTenantSettingsModel{}
		tenantSetting.set(ctx, ts)
		slice = append(slice, &tenantSetting)
	}
	if diags := to.Values.Set(ctx, slice); diags.HasError() {
		return diags
	}

	return nil
}

/*
RESOURCE
*/

type resourceTenantSettingsModel struct {
	baseTenantSettingsModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestUpdateTenantSettings struct {
	fabadmin.UpdateTenantSettingRequest
}

func (to *requestUpdateTenantSettings) set(ctx context.Context, from resourceTenantSettingsModel) diag.Diagnostics {
	if !from.Enabled.IsNull() && !from.Enabled.IsUnknown() {
		to.Enabled = from.Enabled.ValueBoolPointer()
	}

	if !from.DelegateToCapacity.IsNull() && !from.DelegateToCapacity.IsUnknown() {
		to.DelegateToCapacity = from.DelegateToCapacity.ValueBoolPointer()
	}

	if !from.DelegateToDomain.IsNull() && !from.DelegateToDomain.IsUnknown() {
		to.DelegateToDomain = from.DelegateToDomain.ValueBoolPointer()
	}

	if !from.DelegateToWorkspace.IsNull() && !from.DelegateToWorkspace.IsUnknown() {
		to.DelegateToWorkspace = from.DelegateToWorkspace.ValueBoolPointer()
	}

	if !from.EnabledSecurityGroups.IsNull() && !from.EnabledSecurityGroups.IsUnknown() {
		sgs, diags := from.EnabledSecurityGroups.Get(ctx)
		if diags.HasError() {
			return diags
		}

		to.EnabledSecurityGroups = toTenantSettingSecurityGroups(sgs)
	}

	if !from.ExcludedSecurityGroups.IsNull() && !from.ExcludedSecurityGroups.IsUnknown() {
		sgs, diags := from.EnabledSecurityGroups.Get(ctx)
		if diags.HasError() {
			return diags
		}

		to.ExcludedSecurityGroups = toTenantSettingSecurityGroups(sgs)
	}

	if !from.Properties.IsNull() && !from.Properties.IsUnknown() {
		props, diags := from.Properties.Get(ctx)
		if diags.HasError() {
			return diags
		}

		to.Properties = toTenantSettingProperties(props)
	}

	return nil
}

func (to *baseTenantSettingsModel) setUpdate(ctx context.Context, from []fabadmin.TenantSetting) diag.Diagnostics {
	for _, entity := range from {
		if entity.SettingName != nil && *entity.SettingName == to.SettingName.ValueString() {
			return to.set(ctx, entity)
		}
	}
	var diags diag.Diagnostics

	diags.AddError(
		common.ErrorReadHeader,
		"Error during updating Tenant Settings.",
	)

	return diags
}

/*
HELPER MODELS
*/

type tenantSettingsSecurityGroup struct {
	GraphID customtypes.UUID `tfsdk:"graph_id"`
	Name    types.String     `tfsdk:"name"`
}

func (to *tenantSettingsSecurityGroup) set(from fabadmin.TenantSettingSecurityGroup) {
	to.GraphID = customtypes.NewUUIDPointerValue(from.GraphID)
	to.Name = types.StringPointerValue(from.Name)
}

type tenantSettingsProperty struct {
	Name  types.String `tfsdk:"name"`
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

func (to *tenantSettingsProperty) set(from fabadmin.TenantSettingProperty) {
	to.Name = types.StringPointerValue(from.Name)
	to.Type = types.StringPointerValue((*string)(from.Type))
	to.Value = types.StringPointerValue(from.Value)
}

type DeleteBehaviour string

const (
	// NoChange - The setting is not changed during delete operation.
	NoChange DeleteBehaviour = "NoChange"
	// Disable - The setting is disabled during delete operation.
	Disable DeleteBehaviour = "Disable"
)

// PossibleDeleteBehaviourValues returns the possible values for the DeleteBehaviour const type.
func PossibleDeleteBehaviourValues() []DeleteBehaviour {
	return []DeleteBehaviour{
		NoChange,
		Disable,
	}
}

func toTenantSettingSecurityGroups(sgs []*tenantSettingsSecurityGroup) []fabadmin.TenantSettingSecurityGroup {
	if len(sgs) == 0 {
		return nil
	}

	slice := make([]fabadmin.TenantSettingSecurityGroup, 0, len(sgs))

	for _, sg := range sgs {
		securityGroup := fabadmin.TenantSettingSecurityGroup{
			GraphID: sg.GraphID.ValueStringPointer(),
			Name:    sg.Name.ValueStringPointer(),
		}
		slice = append(slice, securityGroup)
	}

	return slice
}

func toTenantSettingProperties(props []*tenantSettingsProperty) []fabadmin.TenantSettingProperty {
	if len(props) == 0 {
		return nil
	}

	slice := make([]fabadmin.TenantSettingProperty, 0, len(props))

	for _, prop := range props {
		properties := fabadmin.TenantSettingProperty{
			Name:  prop.Name.ValueStringPointer(),
			Type:  (*fabadmin.TenantSettingPropertyType)(prop.Type.ValueStringPointer()),
			Value: prop.Value.ValueStringPointer(),
		}
		slice = append(slice, properties)
	}

	return slice
}
