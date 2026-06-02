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
	RoleName      types.String                                      `tfsdk:"role_name"`
	Kind          types.String                                      `tfsdk:"kind"`
	DecisionRules supertypes.SetNestedObjectValueOf[decisionRule]   `tfsdk:"decision_rules"`
	Members       supertypes.SingleNestedObjectValueOf[memberModel] `tfsdk:"members"`
}

type decisionRule struct {
	Effect          types.String                                           `tfsdk:"effect"`
	PermissionScope supertypes.SetNestedObjectValueOf[permissionScope]     `tfsdk:"permission"`
	Constraints     supertypes.SingleNestedObjectValueOf[constraintsModel] `tfsdk:"constraints"`
}

type constraintsModel struct {
	Columns supertypes.SetNestedObjectValueOf[columnConstraint] `tfsdk:"columns"`
	Rows    supertypes.SetNestedObjectValueOf[rowConstraint]    `tfsdk:"rows"`
}

type columnConstraint struct {
	ColumnAction supertypes.SetValueOf[types.String] `tfsdk:"column_action"`
	ColumnEffect types.String                        `tfsdk:"column_effect"`
	ColumnNames  supertypes.SetValueOf[types.String] `tfsdk:"column_names"`
	TablePath    types.String                        `tfsdk:"table_path"`
}

type rowConstraint struct {
	TablePath types.String `tfsdk:"table_path"`
	Value     types.String `tfsdk:"value"`
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
	roleName *string,
	kind *fabcore.DataAccessRoleKind,
	decisionRules []fabcore.DecisionRule,
	members *fabcore.Members,
) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.ItemID = customtypes.NewUUIDValue(itemID)
	to.RoleName = types.StringPointerValue(roleName)
	to.Kind = types.StringPointerValue((*string)(kind))

	if diags := to.setDecisionRules(ctx, decisionRules); diags.HasError() {
		return diags
	}

	return to.setMembers(ctx, members)
}

func (to *baseOneLakeDataAccessSecurityModel) set(ctx context.Context, workspaceID, itemID string, from fabcore.DataAccessRoleBase) diag.Diagnostics {
	return to.setFromBase(ctx, workspaceID, itemID, from.Name, from.Kind, from.DecisionRules, from.Members)
}

func (to *baseOneLakeDataAccessSecurityModel) setFromListItem(ctx context.Context, workspaceID, itemID string, from fabcore.DataAccessRoleListItem) diag.Diagnostics {
	return to.setFromBase(ctx, workspaceID, itemID, from.Name, from.Kind, from.DecisionRules, from.Members)
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
	to.Constraints = supertypes.NewSingleNestedObjectValueOfNull[constraintsModel](ctx)

	if from.Constraints != nil {
		c := &constraintsModel{}
		if diags := c.set(ctx, *from.Constraints); diags.HasError() {
			return diags
		}

		if diags := to.Constraints.Set(ctx, c); diags.HasError() {
			return diags
		}
	}

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

func (to *constraintsModel) set(ctx context.Context, from fabcore.DecisionRuleConstraints) diag.Diagnostics {
	to.Columns = supertypes.NewSetNestedObjectValueOfNull[columnConstraint](ctx)
	to.Rows = supertypes.NewSetNestedObjectValueOfNull[rowConstraint](ctx)

	if from.Columns != nil {
		columns := make([]*columnConstraint, 0, len(from.Columns))

		for _, col := range from.Columns {
			c := &columnConstraint{}
			if diags := c.set(ctx, col); diags.HasError() {
				return diags
			}

			columns = append(columns, c)
		}

		if diags := to.Columns.Set(ctx, columns); diags.HasError() {
			return diags
		}
	}

	if from.Rows != nil {
		rows := make([]*rowConstraint, 0, len(from.Rows))

		for _, row := range from.Rows {
			rows = append(rows, &rowConstraint{
				TablePath: types.StringPointerValue(row.TablePath),
				Value:     types.StringPointerValue(row.Value),
			})
		}

		if diags := to.Rows.Set(ctx, rows); diags.HasError() {
			return diags
		}
	}

	return nil
}

func (to *columnConstraint) set(ctx context.Context, from fabcore.ColumnConstraint) diag.Diagnostics {
	to.ColumnEffect = types.StringPointerValue((*string)(from.ColumnEffect))
	to.TablePath = types.StringPointerValue(from.TablePath)

	to.ColumnAction = supertypes.NewSetValueOfNull[types.String](ctx)

	if from.ColumnAction != nil {
		actions := make([]types.String, 0, len(from.ColumnAction))
		for _, a := range from.ColumnAction {
			actions = append(actions, types.StringValue(string(a)))
		}

		if diags := to.ColumnAction.Set(ctx, actions); diags.HasError() {
			return diags
		}
	}

	to.ColumnNames = supertypes.NewSetValueOfNull[types.String](ctx)

	if from.ColumnNames != nil {
		names := make([]types.String, 0, len(from.ColumnNames))
		for _, n := range from.ColumnNames {
			names = append(names, types.StringValue(n))
		}

		if diags := to.ColumnNames.Set(ctx, names); diags.HasError() {
			return diags
		}
	}

	return nil
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
	to.Name = from.RoleName.ValueStringPointer()

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

		constraints, diags := buildDecisionRuleConstraints(ctx, rule.Constraints)
		if diags.HasError() {
			return diags
		}

		decisionRule.Constraints = constraints

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

func buildDecisionRuleConstraints(ctx context.Context, src supertypes.SingleNestedObjectValueOf[constraintsModel]) (*fabcore.DecisionRuleConstraints, diag.Diagnostics) {
	if src.IsNull() || src.IsUnknown() {
		return nil, nil
	}

	constraints, diags := src.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	out := &fabcore.DecisionRuleConstraints{}

	columns, diags := buildColumnConstraints(ctx, constraints.Columns)
	if diags.HasError() {
		return nil, diags
	}

	out.Columns = columns

	rows, diags := buildRowConstraints(ctx, constraints.Rows)
	if diags.HasError() {
		return nil, diags
	}

	out.Rows = rows

	return out, nil
}

func buildColumnConstraints(ctx context.Context, src supertypes.SetNestedObjectValueOf[columnConstraint]) ([]fabcore.ColumnConstraint, diag.Diagnostics) {
	columns, diags := src.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	if len(columns) == 0 {
		return nil, nil
	}

	out := make([]fabcore.ColumnConstraint, 0, len(columns))
	for _, col := range columns {
		cc := fabcore.ColumnConstraint{
			ColumnEffect: (*fabcore.ColumnEffect)(col.ColumnEffect.ValueStringPointer()),
			TablePath:    col.TablePath.ValueStringPointer(),
		}

		columnActions, diags := col.ColumnAction.Get(ctx)
		if diags.HasError() {
			return nil, diags
		}

		if len(columnActions) > 0 {
			cc.ColumnAction = make([]fabcore.ColumnAction, 0, len(columnActions))
			for _, a := range columnActions {
				cc.ColumnAction = append(cc.ColumnAction, fabcore.ColumnAction(a.ValueString()))
			}
		}

		columnNames, diags := col.ColumnNames.Get(ctx)
		if diags.HasError() {
			return nil, diags
		}

		if len(columnNames) > 0 {
			cc.ColumnNames = make([]string, 0, len(columnNames))
			for _, n := range columnNames {
				cc.ColumnNames = append(cc.ColumnNames, n.ValueString())
			}
		}

		out = append(out, cc)
	}

	return out, nil
}

func buildRowConstraints(ctx context.Context, src supertypes.SetNestedObjectValueOf[rowConstraint]) ([]fabcore.RowConstraint, diag.Diagnostics) {
	rows, diags := src.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	if len(rows) == 0 {
		return nil, nil
	}

	out := make([]fabcore.RowConstraint, 0, len(rows))
	for _, row := range rows {
		out = append(out, fabcore.RowConstraint{
			TablePath: row.TablePath.ValueStringPointer(),
			Value:     row.Value.ValueStringPointer(),
		})
	}

	return out, nil
}
