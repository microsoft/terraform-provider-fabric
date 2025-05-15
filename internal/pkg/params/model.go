package params

type ParametersModel struct {
	Type  string `tfsdk:"type"`
	Find  string `tfsdk:"find"`
	Value string `tfsdk:"value"`
}
