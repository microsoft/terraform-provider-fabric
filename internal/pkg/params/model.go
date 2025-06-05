package params

import "github.com/hashicorp/terraform-plugin-framework/types"

type ParametersModel struct {
	Type  types.String `tfsdk:"type"`
	Find  types.String `tfsdk:"find"`
	Value types.String `tfsdk:"value"`
}
