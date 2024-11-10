package resourceinterfaceunit

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type tfModelFamilyInet struct {
	ArpMaxCache types.Int64 `tfsdk:"arp_max_cache"`
}

func (t *tfModelFamilyInet) attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"arp_max_cache": schema.Int64Attribute{Optional: true},
	}
}

func (t *tfModelFamilyInet) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"arp_max_cache": types.Int64Type,
	}
}

func (t *tfModelFamilyInet) toXmlStruct(ctx context.Context, diags *diag.Diagnostics) *xmlModelFamilyInet {
	if t == nil {
		return nil
	}

	x := new(xmlModelFamilyInet)

	if !t.ArpMaxCache.IsNull() {
		x.ArpMaxCache = t.ArpMaxCache.ValueInt64Pointer()
	}

	return x
}

func tfModelFamilyInetFromTypesObject(ctx context.Context, in types.Object, diags *diag.Diagnostics) *tfModelFamilyInet {
	var result tfModelFamilyInet
	diags.Append(in.As(ctx, &result, basetypes.ObjectAsOptions{})...)
	return &result
}
