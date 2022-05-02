package validators

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ tfsdk.AttributeValidator = StringInSliceValidator{}

type StringInSliceValidator []string

func (v StringInSliceValidator) Description(ctx context.Context) string {
	return "must be one of: " + strings.Join(v, ", ")
}

func (v StringInSliceValidator) MarkdownDescription(ctx context.Context) string {
	return "must be one of: `" + strings.Join(v, "`, `") + "`"
}

func (v StringInSliceValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var attr types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &attr)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	for _, v := range v {
		if v == attr.Value {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(req.AttributePath, "Invalid attribute value", "must be one of: "+strings.Join(v.Values, ", "))
}
