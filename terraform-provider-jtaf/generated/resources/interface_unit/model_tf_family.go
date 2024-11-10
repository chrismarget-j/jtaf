package resourceinterfaceunit

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type tfModelFamily struct {
	Inet types.Object `tfsdk:"inet"`
}

func (t *tfModelFamily) attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"inet": schema.SingleNestedAttribute{Optional: true, Attributes: (*tfModelFamilyInet)(nil).attributes()},
	}
}

func (t *tfModelFamily) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"inet": types.ObjectType{AttrTypes: (*tfModelFamilyInet)(nil).attrTypes()},
	}
}

func (t *tfModelFamily) toXmlStruct(ctx context.Context, diags *diag.Diagnostics) *xmlModelFamily {
	if t == nil {
		return nil
	}

	x := new(xmlModelFamily)

	if !t.Inet.IsNull() {
		x.Inet = tfModelFamilyInetFromTypesObject(ctx, t.Inet, diags).toXmlStruct(ctx, diags)
	}

	return x
}

func tfModelFamilyFromTypesObject(ctx context.Context, in types.Object, diags *diag.Diagnostics) *tfModelFamily {
	var result tfModelFamily
	diags.Append(in.As(ctx, &result, basetypes.ObjectAsOptions{})...)
	return &result
}
