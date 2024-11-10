package resourceinterfaceunit

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type tfModelFamily struct {
	Inet types.Object `tfsdk:"inet"`
}

func (t *tfModelFamily) attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"inet": schema.SingleNestedAttribute{Optional: true, Attributes: (*tfModelFamilyInet)(nil).attributes()},
	}
}

func (t tfModelFamily) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"inet": types.ObjectType{AttrTypes: (*tfModelFamilyInet)(nil).attrTypes()},
	}
}
