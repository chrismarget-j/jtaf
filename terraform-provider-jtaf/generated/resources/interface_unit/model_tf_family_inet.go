package resourceinterfaceunit

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
