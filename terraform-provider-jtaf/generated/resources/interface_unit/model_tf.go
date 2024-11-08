package resourceinterfaceunit

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type tfModel struct {
	Id                types.String `tfsdk:"id"`
	Path              types.String `tfsdk:"path"`
	Name              types.String `tfsdk:"name"`
	ParentPath        types.String `tfsdk:"parent_path"`
	Description       types.String `tfsdk:"description"`
	NativeInnerVlanId types.Int64  `tfsdk:"native_inner_vlan_id"`
}

func (t *tfModel) attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":                   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"path":                 schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"name":                 schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"parent_path":          schema.StringAttribute{Required: true},
		"description":          schema.StringAttribute{Optional: true},
		"native_inner_vlan_id": schema.Int64Attribute{Optional: true},
	}
}
