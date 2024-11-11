package resourceinterfaceunit

import (
	"context"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ common.AttrTyper = (*tfModelFamilyInet)(nil)

type tfModelFamilyInet struct {
	ArpMaxCache types.Int64 `tfsdk:"arp_max_cache"`
}

func (t *tfModelFamilyInet) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"arp_max_cache": types.Int64Type,
	}
}

func (t *tfModelFamilyInet) attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"arp_max_cache": schema.Int64Attribute{Optional: true},
	}
}

func (t *tfModelFamilyInet) loadXmlData(ctx context.Context, x *xmlModelFamilyInet, diags *diag.Diagnostics) {
	if x == nil {
		return
	}

	t.ArpMaxCache = types.Int64PointerValue(x.ArpMaxCache)
}
