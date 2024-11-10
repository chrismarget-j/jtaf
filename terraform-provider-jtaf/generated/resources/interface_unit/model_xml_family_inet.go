package resourceinterfaceunit

import (
	"context"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type xmlModelFamilyInet struct {
	ArpMaxCache *int64 `xml:"arp-max-cache,omitempty"`
}

func (x *xmlModelFamilyInet) toTF(ctx context.Context, diags *diag.Diagnostics) types.Object {
	if x == nil {
		return types.ObjectNull((*tfModelFamilyInet)(nil).attrTypes())
	}

	var tf tfModelFamilyInet

	if x.ArpMaxCache != nil {
		tf.ArpMaxCache = types.Int64PointerValue(x.ArpMaxCache)
	}

	return common.ObjectValueFromWithDiags(ctx, (*tfModelFamilyInet)(nil).attrTypes(), tf, diags)
}
