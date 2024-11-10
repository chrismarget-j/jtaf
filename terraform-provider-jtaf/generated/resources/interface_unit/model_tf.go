package resourceinterfaceunit

import (
	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type tfModel struct {
	Id                types.String `tfsdk:"id"`
	XPath             types.String `tfsdk:"xpath"`
	Name              types.String `tfsdk:"name"`
	ParentXPath       types.String `tfsdk:"parent_xpath"`
	Description       types.String `tfsdk:"description"`
	NativeInnerVlanId types.Int64  `tfsdk:"native_inner_vlan_id"`
	Family            types.Object `tfsdk:"family"`
}

func (t *tfModel) attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":                   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"xpath":                schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"name":                 schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"parent_xpath":         schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.RegexMatches(common.XPathRegex, common.XPathRegexMsg)}},
		"description":          schema.StringAttribute{Optional: true},
		"native_inner_vlan_id": schema.Int64Attribute{Optional: true},
		"family":               schema.SingleNestedAttribute{Optional: true, Attributes: (*tfModelFamily)(nil).attributes()},
	}
}

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
