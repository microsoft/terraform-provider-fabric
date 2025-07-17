// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package typeutils

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
)

const nullStr = "null"

// DynamicToJSON converts dynamic types to JSON.
func DynamicToJSON(d types.Dynamic) ([]byte, error) {
	if d.IsNull() {
		return nil, nil
	}

	return attrValueToJSON(d.UnderlyingValue())
}

func attrListToJSON(in []attr.Value) ([]json.RawMessage, error) {
	l := make([]json.RawMessage, 0)

	for _, v := range in {
		vv, err := attrValueToJSON(v)
		if err != nil {
			return nil, err
		}

		l = append(l, json.RawMessage(vv))
	}

	return l, nil
}

func attrMapToJSON(in map[string]attr.Value) (map[string]json.RawMessage, error) {
	m := map[string]json.RawMessage{}

	for k, v := range in {
		vv, err := attrValueToJSON(v)
		if err != nil {
			return nil, err
		}

		m[k] = json.RawMessage(vv)
	}

	return m, nil
}

func attrValueToJSON(val attr.Value) ([]byte, error) {
	if val.IsNull() {
		return json.Marshal(nil)
	}

	switch value := val.(type) {
	case types.Bool:
		return json.Marshal(value.ValueBool())
	case types.String:
		return json.Marshal(value.ValueString())
	case types.Int64:
		return json.Marshal(value.ValueInt64())
	case types.Float64:
		return json.Marshal(value.ValueFloat64())
	case types.Number:
		v, _ := value.ValueBigFloat().Float64()

		return json.Marshal(v)
	case types.List:
		l, err := attrListToJSON(value.Elements())
		if err != nil {
			return nil, err
		}

		return json.Marshal(l)
	case types.Set:
		l, err := attrListToJSON(value.Elements())
		if err != nil {
			return nil, err
		}

		return json.Marshal(l)
	case types.Tuple:
		l, err := attrListToJSON(value.Elements())
		if err != nil {
			return nil, err
		}

		return json.Marshal(l)
	case types.Map:
		m, err := attrMapToJSON(value.Elements())
		if err != nil {
			return nil, err
		}

		return json.Marshal(m)
	case types.Object:
		m, err := attrMapToJSON(value.Attributes())
		if err != nil {
			return nil, err
		}

		return json.Marshal(m)
	default:
		return nil, fmt.Errorf("unhandled type: %T", value)
	}
}

// JSONToDynamic converts JSON to dynamic types.
func JSONToDynamic(b []byte, typ attr.Type) (types.Dynamic, error) {
	v, err := attrValueFromJSON(b, typ)
	if err != nil {
		return types.Dynamic{}, err
	}

	return types.DynamicValue(v), nil
}

func attrListFromJSON(b []byte, etyp attr.Type) ([]attr.Value, error) {
	var l []json.RawMessage

	err := json.Unmarshal(b, &l)
	if err != nil {
		return nil, err
	}

	vals := make([]attr.Value, 0)

	for _, b := range l {
		val, err := attrValueFromJSON(b, etyp)
		if err != nil {
			return nil, err
		}

		vals = append(vals, val)
	}

	return vals, nil
}

func attrValueFromJSON(b []byte, typ attr.Type) (attr.Value, error) { //nolint:gocyclo, gocognit, maintidx
	switch typ := typ.(type) {
	case basetypes.BoolType:
		if b == nil || string(b) == nullStr {
			return types.BoolNull(), nil
		}

		var v bool

		err := json.Unmarshal(b, &v)
		if err != nil {
			return nil, err
		}

		return types.BoolValue(v), nil
	case basetypes.StringType:
		if b == nil || string(b) == nullStr {
			return types.StringNull(), nil
		}

		var v string

		err := json.Unmarshal(b, &v)
		if err != nil {
			return nil, err
		}

		return types.StringValue(v), nil
	case basetypes.Int64Type:
		if b == nil || string(b) == nullStr {
			return types.Int64Null(), nil
		}

		var v int64

		err := json.Unmarshal(b, &v)
		if err != nil {
			return nil, err
		}

		return types.Int64Value(v), nil
	case basetypes.Float64Type:
		if b == nil || string(b) == nullStr {
			return types.Float64Null(), nil
		}

		var v float64

		err := json.Unmarshal(b, &v)
		if err != nil {
			return nil, err
		}

		return types.Float64Value(v), nil
	case basetypes.NumberType:
		if b == nil || string(b) == nullStr {
			return types.NumberNull(), nil
		}

		var v float64

		err := json.Unmarshal(b, &v)
		if err != nil {
			return nil, err
		}

		return types.NumberValue(big.NewFloat(v)), nil
	case basetypes.ListType:
		if b == nil || string(b) == nullStr {
			return types.ListNull(typ.ElemType), nil
		}

		vals, err := attrListFromJSON(b, typ.ElemType)
		if err != nil {
			return nil, err
		}

		vv, diags := types.ListValue(typ.ElemType, vals)
		if diags.HasError() {
			diag := diags.Errors()[0]

			return nil, fmt.Errorf("%s: %s", diag.Summary(), diag.Detail())
		}

		return vv, nil
	case basetypes.SetType:
		if b == nil || string(b) == nullStr {
			return types.SetNull(typ.ElemType), nil
		}

		vals, err := attrListFromJSON(b, typ.ElemType)
		if err != nil {
			return nil, err
		}

		vv, diags := types.SetValue(typ.ElemType, vals)
		if diags.HasError() {
			diag := diags.Errors()[0]

			return nil, fmt.Errorf("%s: %s", diag.Summary(), diag.Detail())
		}

		return vv, nil
	case basetypes.TupleType:
		if b == nil || string(b) == nullStr {
			return types.TupleNull(typ.ElemTypes), nil
		}

		var l []json.RawMessage

		err := json.Unmarshal(b, &l)
		if err != nil {
			return nil, err
		}

		if len(l) != len(typ.ElemTypes) {
			return nil, fmt.Errorf("tuple element size not match: json=%d, type=%d", len(l), len(typ.ElemTypes))
		}

		vals := make([]attr.Value, 0)

		for i, b := range l {
			val, err := attrValueFromJSON(b, typ.ElemTypes[i])
			if err != nil {
				return nil, err
			}

			vals = append(vals, val)
		}

		vv, diags := types.TupleValue(typ.ElemTypes, vals)
		if diags.HasError() {
			diag := diags.Errors()[0]

			return nil, fmt.Errorf("%s: %s", diag.Summary(), diag.Detail())
		}

		return vv, nil
	case basetypes.MapType:
		if b == nil || string(b) == nullStr {
			return types.MapNull(typ.ElemType), nil
		}

		var m map[string]json.RawMessage

		err := json.Unmarshal(b, &m)
		if err != nil {
			return nil, err
		}

		vals := map[string]attr.Value{}

		for k, v := range m {
			val, err := attrValueFromJSON(v, typ.ElemType)
			if err != nil {
				return nil, err
			}

			vals[k] = val
		}

		vv, diags := types.MapValue(typ.ElemType, vals)
		if diags.HasError() {
			diag := diags.Errors()[0]

			return nil, fmt.Errorf("%s: %s", diag.Summary(), diag.Detail())
		}

		return vv, nil
	case basetypes.ObjectType:
		if b == nil || string(b) == nullStr {
			return types.ObjectNull(typ.AttributeTypes()), nil
		}

		var m map[string]json.RawMessage

		err := json.Unmarshal(b, &m)
		if err != nil {
			return nil, err
		}

		vals := map[string]attr.Value{}
		attrTypes := typ.AttributeTypes()

		for k, attrType := range attrTypes {
			val, err := attrValueFromJSON(m[k], attrType)
			if err != nil {
				return nil, err
			}

			vals[k] = val
		}

		vv, diags := types.ObjectValue(attrTypes, vals)
		if diags.HasError() {
			diag := diags.Errors()[0]

			return nil, fmt.Errorf("%s: %s", diag.Summary(), diag.Detail())
		}

		return vv, nil
	case basetypes.DynamicType:
		if b == nil || string(b) == nullStr {
			return types.DynamicNull(), nil
		}

		return JSONToDynamicImplied(b)
	default:
		return nil, fmt.Errorf("unhandled type: %T", typ)
	}
}

// JSONToDynamicImplied is similar to FromJSON, while it is for typeless case.
// In which case, the following type conversion rules are applied (Go -> TF):
// - bool: bool
// - float64: number
// - string: string
// - []interface{}: tuple
// - map[string]interface{}: object
// - nil: null (dynamic)
// In case the input json is of zero-length, it returns null (dynamic).
func JSONToDynamicImplied(b []byte) (types.Dynamic, error) {
	if len(b) == 0 {
		return types.DynamicNull(), nil
	}

	_, v, err := attrValueFromJSONImplied(b)
	if err != nil {
		return types.Dynamic{}, err
	}

	return types.DynamicValue(v), nil
}

func attrValueFromJSONImplied(b []byte) (attr.Type, attr.Value, error) {
	if string(b) == nullStr {
		return types.DynamicType, types.DynamicNull(), nil
	}

	var object map[string]json.RawMessage

	err := json.Unmarshal(b, &object)
	if err == nil {
		attrTypes := map[string]attr.Type{}
		attrVals := map[string]attr.Value{}

		for k, v := range object {
			attrTypes[k], attrVals[k], err = attrValueFromJSONImplied(v)
			if err != nil {
				return nil, nil, err
			}
		}

		typ := types.ObjectType{AttrTypes: attrTypes}
		val, diags := types.ObjectValue(attrTypes, attrVals)

		if diags.HasError() {
			diag := diags.Errors()[0]

			return nil, nil, fmt.Errorf("%s: %s", diag.Summary(), diag.Detail())
		}

		return typ, val, nil
	}

	var array []json.RawMessage

	err = json.Unmarshal(b, &array)
	if err == nil {
		eTypes := []attr.Type{}
		eVals := []attr.Value{}

		for _, e := range array {
			eType, eVal, err := attrValueFromJSONImplied(e)
			if err != nil {
				return nil, nil, err
			}

			eTypes = append(eTypes, eType)
			eVals = append(eVals, eVal)
		}

		typ := types.TupleType{ElemTypes: eTypes}
		val, diags := types.TupleValue(eTypes, eVals)

		if diags.HasError() {
			diag := diags.Errors()[0]

			return nil, nil, fmt.Errorf("%s: %s", diag.Summary(), diag.Detail())
		}

		return typ, val, nil
	}

	// Primitives
	var v any

	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal %s: %w", string(b), err)
	}

	switch v := v.(type) {
	case bool:
		return types.BoolType, types.BoolValue(v), nil
	case float64:
		return types.NumberType, types.NumberValue(big.NewFloat(v)), nil
	case string:
		return types.StringType, types.StringValue(v), nil
	case nil:
		return types.DynamicType, types.DynamicNull(), nil
	default:
		return nil, nil, fmt.Errorf("unhandled type: %T", v)
	}
}

// IsFullyKnown returns true if `val` is known. If `val` is an aggregate type,
// IsFullyKnown only returns true if all elements and attributes are known, as well.
//
//nolint:gocognit
func IsFullyKnown(val attr.Value) bool {
	if val == nil {
		return true
	}

	if val.IsUnknown() {
		return false
	}

	switch v := val.(type) {
	case types.Dynamic:
		return IsFullyKnown(v.UnderlyingValue())
	case types.List:
		for _, e := range v.Elements() {
			if !IsFullyKnown(e) {
				return false
			}
		}

		return true
	case types.Set:
		for _, e := range v.Elements() {
			if !IsFullyKnown(e) {
				return false
			}
		}

		return true
	case types.Tuple:
		for _, e := range v.Elements() {
			if !IsFullyKnown(e) {
				return false
			}
		}

		return true
	case types.Map:
		for _, e := range v.Elements() {
			if !IsFullyKnown(e) {
				return false
			}
		}

		return true
	case types.Object:
		for _, e := range v.Attributes() {
			if !IsFullyKnown(e) {
				return false
			}
		}

		return true
	default:
		return true
	}
}

func SemanticallyEqual(a, b types.Dynamic) bool {
	aJSON, err := DynamicToJSON(a)
	if err != nil {
		return false
	}

	bJSON, err := DynamicToJSON(b)
	if err != nil {
		return false
	}

	aJSONStr, err := transforms.JSONNormalize(string(aJSON))
	if err != nil {
		return false
	}

	aJSONValue := jsontypes.NewNormalizedValue(aJSONStr)

	bJSONStr, err := transforms.JSONNormalize(string(bJSON))
	if err != nil {
		return false
	}

	bJSONValue := jsontypes.NewNormalizedValue(bJSONStr)

	return aJSONValue.Equal(bJSONValue)
}
