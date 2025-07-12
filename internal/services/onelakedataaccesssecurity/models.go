// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceOneLakeDataAccessSecurityModel struct {
	baseOneLakeDataAccessSecurityModel

	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
	ItemID      customtypes.UUID `tfsdk:"item_id"`
}

type baseOneLakeDataAccessSecurityModel struct {
	Value supertypes.SetNestedObjectValueOf[dataAccessRole] `tfsdk:"value"`
}

type dataAccessRole struct {
	ID            customtypes.UUID                                `tfsdk:"id"`
	Name          types.String                                    `tfsdk:"name"`
	DecisionRules supertypes.SetNestedObjectValueOf[decisionRule] `tfsdk:"decision_rules"`
	Members       supertypes.SingleNestedObjectValueOf[Member]    `tfsdk:"members"`
}

type decisionRule struct {
	Effect          types.String                                       `tfsdk:"effect"`
	PermissionScope supertypes.SetNestedObjectValueOf[permissionScope] `tfsdk:"permission"`
}

type permissionScope struct {
	AttributeName            types.String                        `tfsdk:"attribute_name"`
	AttributeValueIncludedIn supertypes.SetValueOf[types.String] `tfsdk:"attribute_value_included_in"`
}

type Member struct {
	FabricItemMembers     supertypes.SetNestedObjectValueOf[FabricItemMember]     `tfsdk:"fabric_item_members"`
	MicrosoftEntraMembers supertypes.SetNestedObjectValueOf[MicrosoftEntraMember] `tfsdk:"microsoft_entra_members"`
}

type FabricItemMember struct {
	ItemAccess supertypes.SetValueOf[types.String] `tfsdk:"item_access"`
	SourcePath types.String                        `tfsdk:"source_path"`
}

type MicrosoftEntraMember struct {
	ObjectID   customtypes.UUID `tfsdk:"object_id"`
	ObjectType types.String     `tfsdk:"object_type"`
	TenantID   types.String     `tfsdk:"tenant_id"`
}

func (to *dataSourceOneLakeDataAccessSecurityModel) set(ctx context.Context, from fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse) diag.Diagnostics {
	slice := make([]*dataAccessRole, 0, len(from.Value))

	for _, item := range from.Value {
		role := &dataAccessRole{
			ID:   customtypes.NewUUIDPointerValue(item.ID),
			Name: types.StringPointerValue(item.Name),
		}

		role.DecisionRules = supertypes.NewSetNestedObjectValueOfNull[decisionRule](ctx)

		if item.DecisionRules != nil {
			decisionRules := make([]*decisionRule, 0, len(item.DecisionRules))

			for _, rule := range item.DecisionRules {
				decisionRule := decisionRule{}
				if diags := decisionRule.set(ctx, rule); diags.HasError() {
					return diags
				}

				decisionRules = append(decisionRules, &decisionRule)
			}

			if diags := role.DecisionRules.Set(ctx, decisionRules); diags.HasError() {
				return diags
			}
		}

		role.Members.SetNull(ctx)

		members := &Member{}
		if diags := members.set(ctx, item.Members); diags.HasError() {
			return diags
		}

		if diags := role.Members.Set(ctx, members); diags.HasError() {
			return diags
		}

		slice = append(slice, role)
	}

	if diags := to.Value.Set(ctx, slice); diags.HasError() {
		return diags
	}

	return nil
}

func (to *decisionRule) set(ctx context.Context, from fabcore.DecisionRule) diag.Diagnostics {
	to.Effect = types.StringPointerValue((*string)(from.Effect))
	to.PermissionScope.SetNull(ctx)

	if from.Permission != nil {
		permissions := make([]*permissionScope, 0, len(from.Permission))

		for _, perm := range from.Permission {
			permission := permissionScope{}
			if diags := permission.set(ctx, perm); diags.HasError() {
				return diags
			}

			permissions = append(permissions, &permission)
		}

		if diags := to.PermissionScope.Set(ctx, permissions); diags.HasError() {
			return diags
		}
	}

	return nil
}

func (to *permissionScope) set(ctx context.Context, from fabcore.PermissionScope) diag.Diagnostics {
	to.AttributeValueIncludedIn.SetNull(ctx)

	if from.AttributeName != nil {
		to.AttributeName = types.StringPointerValue((*string)(from.AttributeName))
	} else {
		to.AttributeName = types.StringNull()
	}

	if from.AttributeValueIncludedIn != nil {
		values := make([]types.String, 0, len(from.AttributeValueIncludedIn))
		for _, value := range from.AttributeValueIncludedIn {
			values = append(values, types.StringValue(value))
		}

		if diags := to.AttributeValueIncludedIn.Set(ctx, values); diags.HasError() {
			return diags
		}
	}

	return nil
}

func (to *Member) set(ctx context.Context, from *fabcore.Members) diag.Diagnostics {
	if from.FabricItemMembers != nil {
		fabricItemMembers := make([]*FabricItemMember, 0, len(from.FabricItemMembers))

		for _, item := range from.FabricItemMembers {
			fabricItemMember := &FabricItemMember{
				SourcePath: types.StringPointerValue(item.SourcePath),
			}

			fabricItemMember.ItemAccess.SetNull(ctx)

			if item.ItemAccess != nil {
				itemAccess := make([]types.String, 0, len(item.ItemAccess))

				for _, access := range item.ItemAccess {
					itemAccessString := types.StringValue((string)(access))
					itemAccess = append(itemAccess, itemAccessString)
				}

				if diags := fabricItemMember.ItemAccess.Set(ctx, itemAccess); diags.HasError() {
					return diags
				}
			}

			fabricItemMembers = append(fabricItemMembers, fabricItemMember)
		}

		if diags := to.FabricItemMembers.Set(ctx, fabricItemMembers); diags.HasError() {
			return diags
		}
	}

	to.MicrosoftEntraMembers.SetNull(ctx)

	if from.MicrosoftEntraMembers != nil {
		microsoftEntraMembers := make([]*MicrosoftEntraMember, 0, len(from.MicrosoftEntraMembers))

		for _, item := range from.MicrosoftEntraMembers {
			microsoftEntraMember := &MicrosoftEntraMember{
				ObjectID: customtypes.NewUUIDPointerValue(item.ObjectID),
				TenantID: types.StringPointerValue(item.TenantID),
			}

			if item.ObjectType != nil {
				microsoftEntraMember.ObjectType = types.StringPointerValue((*string)(item.ObjectType))
			} else {
				microsoftEntraMember.ObjectType = types.StringNull()
			}

			microsoftEntraMembers = append(microsoftEntraMembers, microsoftEntraMember)
		}

		if diags := to.MicrosoftEntraMembers.Set(ctx, microsoftEntraMembers); diags.HasError() {
			return diags
		}
	}

	return nil
}
