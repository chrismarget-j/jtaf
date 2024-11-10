package resourceinterfaceunit

import (
	"context"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type xmlModelFamily struct {
	Inet *xmlModelFamilyInet `xml:"inet,omitempty"`
}

func (x *xmlModelFamily) toTF(ctx context.Context, diags *diag.Diagnostics) types.Object {
	if x == nil {
		return types.ObjectNull((*tfModelFamily)(nil).attrTypes())
	}

	var tf tfModelFamily

	if x.Inet != nil {
		tf.Inet = x.Inet.toTF(ctx, diags)
	}

	return common.ObjectValueFromWithDiags(ctx, (*tfModelFamily)(nil).attrTypes(), tf, diags)
}
