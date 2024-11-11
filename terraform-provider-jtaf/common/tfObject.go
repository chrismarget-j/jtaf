package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func ObjectValueFromAttrTyper(ctx context.Context, attrTyper AttrTyper, diags *diag.Diagnostics) basetypes.ObjectValue {
	attrTypes := attrTyper.AttrTypes()
	result, d := types.ObjectValueFrom(ctx, attrTypes, attrTyper)
	diags.Append(d...)
	return result
}
