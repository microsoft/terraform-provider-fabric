// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseOneLakeDataAccessSecurityModel struct {
	WorkspaceID   customtypes.UUID                                  `tfsdk:"workspace_id"`
	ItemID        customtypes.UUID                                  `tfsdk:"item_id"`
	ID            customtypes.UUID                                  `tfsdk:"id"`
	Name          types.String                                      `tfsdk:"name"`
	Kind          types.String                                      `tfsdk:"kind"`
	DecisionRules supertypes.SetNestedObjectValueOf[decisionRule]   `tfsdk:"decision_rules"`
	Members       supertypes.SingleNestedObjectValueOf[memberModel] `tfsdk:"members"`
}

type decisionRule struct {
	Effect          types.String                                       `tfsdk:"effect"`
	PermissionScope supertypes.SetNestedObjectValueOf[permissionScope] `tfsdk:"permission"`
}

type permissionScope struct {
	AttributeName            types.String                        `tfsdk:"attribute_name"`
	AttributeValueIncludedIn supertypes.SetValueOf[types.String] `tfsdk:"attribute_value_included_in"`
}

type memberModel struct {
	FabricItemMembers     supertypes.SetNestedObjectValueOf[fabricItemMember]     `tfsdk:"fabric_item_members"`
	MicrosoftEntraMembers supertypes.SetNestedObjectValueOf[microsoftEntraMember] `tfsdk:"microsoft_entra_members"`
}

type fabricItemMember struct {
	ItemAccess supertypes.SetValueOf[types.String] `tfsdk:"item_access"`
	SourcePath types.String                        `tfsdk:"source_path"`
}

type microsoftEntraMember struct {
	ObjectID   customtypes.UUID `tfsdk:"object_id"`
	ObjectType types.String     `tfsdk:"object_type"`
	TenantID   customtypes.UUID `tfsdk:"tenant_id"`
}

func (to *baseOneLakeDataAccessSecurityModel) setFromBase(
	ctx context.Context,
	workspaceID, itemID string,
	name, id *string,
	kind *fabcore.DataAccessRoleKind,
	decisionRules []fabcore.DecisionRule,
	members *fabcore.Members,
) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.ItemID = customtypes.NewUUIDValue(itemID)
	to.ID = customtypes.NewUUIDPointerValue(id)
	to.Name = types.StringPointerValue(name)
	to.Kind = types.StringPointerValue((*string)(kind))

	if diags := to.setDecisionRules(ctx, decisionRules); diags.HasError() {
		return diags
	}

	return to.setMembers(ctx, members)
}

func (to *baseOneLakeDataAccessSecurityModel) set(ctx context.Context, workspaceID, itemID string, from fabcore.DataAccessRoleBase) diag.Diagnostics {
	return to.setFromBase(ctx, workspaceID, itemID, from.Name, nil, from.Kind, from.DecisionRules, from.Members)
}

func (to *baseOneLakeDataAccessSecurityModel) setFromListItem(ctx context.Context, workspaceID, itemID string, from fabcore.DataAccessRoleListItem) diag.Diagnostics {
	return to.setFromBase(ctx, workspaceID, itemID, from.Name, from.ID, from.Kind, from.DecisionRules, from.Members)
}

func (to *baseOneLakeDataAccessSecurityModel) setDecisionRules(ctx context.Context, from []fabcore.DecisionRule) diag.Diagnostics {
	to.DecisionRules = supertypes.NewSetNestedObjectValueOfNull[decisionRule](ctx)

	if from == nil {
		return nil
	}

	rules := make([]*decisionRule, 0, len(from))

	for _, rule := range from {
		r := &decisionRule{}
		if diags := r.set(ctx, rule); diags.HasError() {
			return diags
		}

		rules = append(rules, r)
	}

	return to.DecisionRules.Set(ctx, rules)
}

func (to *baseOneLakeDataAccessSecurityModel) setMembers(ctx context.Context, from *fabcore.Members) diag.Diagnostics {
	to.Members = supertypes.NewSingleNestedObjectValueOfNull[memberModel](ctx)

	if from == nil {
		return nil
	}

	m := &memberModel{}
	if diags := m.set(ctx, from); diags.HasError() {
		return diags
	}

	return to.Members.Set(ctx, m)
}

func (to *decisionRule) set(ctx context.Context, from fabcore.DecisionRule) diag.Diagnostics {
	to.Effect = types.StringPointerValue((*string)(from.Effect))
	to.PermissionScope = supertypes.NewSetNestedObjectValueOfNull[permissionScope](ctx)

	if from.Permission == nil {
		return nil
	}

	permissions := make([]*permissionScope, 0, len(from.Permission))

	for _, perm := range from.Permission {
		p := &permissionScope{}
		if diags := p.set(ctx, perm); diags.HasError() {
			return diags
		}

		permissions = append(permissions, p)
	}

	return to.PermissionScope.Set(ctx, permissions)
}

func (to *permissionScope) set(ctx context.Context, from fabcore.PermissionScope) diag.Diagnostics {
	to.AttributeName = types.StringPointerValue((*string)(from.AttributeName))
	to.AttributeValueIncludedIn = supertypes.NewSetValueOfNull[types.String](ctx)

	if from.AttributeValueIncludedIn == nil {
		return nil
	}

	values := make([]types.String, 0, len(from.AttributeValueIncludedIn))
	for _, value := range from.AttributeValueIncludedIn {
		values = append(values, types.StringValue(value))
	}

	return to.AttributeValueIncludedIn.Set(ctx, values)
}

func (to *memberModel) set(ctx context.Context, from *fabcore.Members) diag.Diagnostics {
	to.FabricItemMembers = supertypes.NewSetNestedObjectValueOfNull[fabricItemMember](ctx)
	to.MicrosoftEntraMembers = supertypes.NewSetNestedObjectValueOfNull[microsoftEntraMember](ctx)

	if from.FabricItemMembers != nil {
		fabricItemMembers := make([]*fabricItemMember, 0, len(from.FabricItemMembers))

		for _, item := range from.FabricItemMembers {
			fim := &fabricItemMember{
				SourcePath: types.StringPointerValue(item.SourcePath),
			}

			fim.ItemAccess = supertypes.NewSetValueOfNull[types.String](ctx)

			if item.ItemAccess != nil {
				itemAccess := make([]types.String, 0, len(item.ItemAccess))
				for _, access := range item.ItemAccess {
					itemAccess = append(itemAccess, types.StringValue(string(access)))
				}

				if diags := fim.ItemAccess.Set(ctx, itemAccess); diags.HasError() {
					return diags
				}
			}

			fabricItemMembers = append(fabricItemMembers, fim)
		}

		if diags := to.FabricItemMembers.Set(ctx, fabricItemMembers); diags.HasError() {
			return diags
		}
	}

	if from.MicrosoftEntraMembers != nil {
		entraMembers := make([]*microsoftEntraMember, 0, len(from.MicrosoftEntraMembers))

		for _, item := range from.MicrosoftEntraMembers {
			entraMembers = append(entraMembers, &microsoftEntraMember{
				ObjectID:   customtypes.NewUUIDPointerValue(item.ObjectID),
				ObjectType: types.StringPointerValue((*string)(item.ObjectType)),
				TenantID:   customtypes.NewUUIDPointerValue(item.TenantID),
			})
		}

		if diags := to.MicrosoftEntraMembers.Set(ctx, entraMembers); diags.HasError() {
			return diags
		}
	}

	return nil
}

/*
DATA-SOURCE (single)
*/

type dataSourceOneLakeDataAccessSecurityModel struct {
	baseOneLakeDataAccessSecurityModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceOneLakeDataAccessSecuritiesModel struct {
	WorkspaceID customtypes.UUID                                                      `tfsdk:"workspace_id"`
	ItemID      customtypes.UUID                                                      `tfsdk:"item_id"`
	Values      supertypes.SetNestedObjectValueOf[baseOneLakeDataAccessSecurityModel] `tfsdk:"values"`
	Timeouts    timeoutsD.Value                                                       `tfsdk:"timeouts"`
}

func (to *dataSourceOneLakeDataAccessSecuritiesModel) setValues(ctx context.Context, workspaceID, itemID string, from []fabcore.DataAccessRoleListItem) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.ItemID = customtypes.NewUUIDValue(itemID)

	slice := make([]*baseOneLakeDataAccessSecurityModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseOneLakeDataAccessSecurityModel

		if diags := entityModel.setFromListItem(ctx, workspaceID, itemID, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

/*
RESOURCE
*/

type resourceOneLakeDataAccessSecurityModel struct {
	baseOneLakeDataAccessSecurityModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateOrUpdateOneLakeDataAccessSecurity struct {
	fabcore.DataAccessRoleBase
}

func (to *requestCreateOrUpdateOneLakeDataAccessSecurity) set(ctx context.Context, from resourceOneLakeDataAccessSecurityModel) diag.Diagnostics {
	to.Name = from.Name.ValueStringPointer()

	if !from.Kind.IsNull() && !from.Kind.IsUnknown() {
		to.Kind = (*fabcore.DataAccessRoleKind)(from.Kind.ValueStringPointer())
	}

	decisionRules, diags := from.DecisionRules.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.DecisionRules = make([]fabcore.DecisionRule, 0, len(decisionRules))

	for _, rule := range decisionRules {
		decisionRule := fabcore.DecisionRule{
			Effect: (*fabcore.Effect)(rule.Effect.ValueStringPointer()),
		}

		permissions, diags := rule.PermissionScope.Get(ctx)
		if diags.HasError() {
			return diags
		}

		decisionRule.Permission = make([]fabcore.PermissionScope, 0, len(permissions))

		for _, permission := range permissions {
			attributesIncludedIn, diags := permission.AttributeValueIncludedIn.Get(ctx)
			if diags.HasError() {
				return diags
			}

			permissionScope := fabcore.PermissionScope{
				AttributeName: (*fabcore.AttributeName)(permission.AttributeName.ValueStringPointer()),
			}

			if len(attributesIncludedIn) > 0 {
				permissionScope.AttributeValueIncludedIn = make([]string, 0, len(attributesIncludedIn))
				for _, attr := range attributesIncludedIn {
					permissionScope.AttributeValueIncludedIn = append(permissionScope.AttributeValueIncludedIn, attr.ValueString())
				}
			}

			decisionRule.Permission = append(decisionRule.Permission, permissionScope)
		}

		to.DecisionRules = append(to.DecisionRules, decisionRule)
	}

	return to.setMembers(ctx, from)
}

func (to *requestCreateOrUpdateOneLakeDataAccessSecurity) setMembers(ctx context.Context, from resourceOneLakeDataAccessSecurityModel) diag.Diagnostics {
	if from.Members.IsNull() || from.Members.IsUnknown() {
		return nil
	}

	members, diags := from.Members.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.Members = &fabcore.Members{}

	fabricItemMembers, diags := members.FabricItemMembers.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.Members.FabricItemMembers = make([]fabcore.FabricItemMember, 0, len(fabricItemMembers))

	for _, fim := range fabricItemMembers {
		member := fabcore.FabricItemMember{
			SourcePath: fim.SourcePath.ValueStringPointer(),
		}

		itemAccess, diags := fim.ItemAccess.Get(ctx)
		if diags.HasError() {
			return diags
		}

		if len(itemAccess) > 0 {
			member.ItemAccess = make([]fabcore.ItemAccess, 0, len(itemAccess))
			for _, access := range itemAccess {
				member.ItemAccess = append(member.ItemAccess, fabcore.ItemAccess(access.ValueString()))
			}
		}

		to.Members.FabricItemMembers = append(to.Members.FabricItemMembers, member)
	}

	microsoftEntraMembers, diags := members.MicrosoftEntraMembers.Get(ctx)
	if diags.HasError() {
		return diags
	}

	if len(microsoftEntraMembers) > 0 {
		to.Members.MicrosoftEntraMembers = make([]fabcore.MicrosoftEntraMember, 0, len(microsoftEntraMembers))
		for _, mem := range microsoftEntraMembers {
			to.Members.MicrosoftEntraMembers = append(to.Members.MicrosoftEntraMembers, fabcore.MicrosoftEntraMember{
				ObjectID:   mem.ObjectID.ValueStringPointer(),
				TenantID:   mem.TenantID.ValueStringPointer(),
				ObjectType: (*fabcore.ObjectType)(mem.ObjectType.ValueStringPointer()),
			})
		}
	}

	return nil
}
