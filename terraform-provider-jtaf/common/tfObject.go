package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func ObjectValueFromWithDiags(ctx context.Context, attributeTypes map[string]attr.Type, attributes any, diags *diag.Diagnostics) basetypes.ObjectValue {
	o, d := types.ObjectValueFrom(ctx, attributeTypes, attributes)
	diags.Append(d...)
	return o
}
