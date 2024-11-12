package resourceinterfaceunit

import (
	"context"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ common.AttrTyper = (*tfModelFamily)(nil)

type tfModelFamily struct {
	Inet types.Object `tfsdk:"inet"`
}

func (t *tfModelFamily) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"inet": types.ObjectType{AttrTypes: (*tfModelFamilyInet)(nil).AttrTypes()},
	}
}

func (t *tfModelFamily) attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"inet": schema.SingleNestedAttribute{Optional: true, Attributes: (*tfModelFamilyInet)(nil).attributes(), PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()}},
	}
}

func (t *tfModelFamily) loadXmlData(ctx context.Context, x *xmlModelFamily, diags *diag.Diagnostics) {
	if x == nil {
		return
	}

	if x.Inet != nil {
		var inet tfModelFamilyInet
		inet.loadXmlData(ctx, x.Inet, diags)
		t.Inet = common.ObjectValueFromAttrTyper(ctx, &inet, diags)
	}
}
