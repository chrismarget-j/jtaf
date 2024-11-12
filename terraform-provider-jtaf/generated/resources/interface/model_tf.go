package resourceinterface

import (
	"context"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ common.AttrTyper = (*tfModel)(nil)

type tfModel struct {
	Id          types.String `tfsdk:"id"`
	XPath       types.String `tfsdk:"xpath"`
	Name        types.String `tfsdk:"name"`
	ParentXPath types.String `tfsdk:"parent_xpath"`
	Description types.String `tfsdk:"description"`
	Mtu         types.Int64  `tfsdk:"mtu"`
}

func (t *tfModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":           types.StringType,
		"xpath":        types.StringType,
		"name":         types.StringType,
		"parent_xpath": types.StringType,
		"description":  types.StringType,
		"mtu":          types.Int64Type,
	}
}

func (t *tfModel) attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"xpath":        schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"name":         schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"parent_xpath": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Validators: []validator.String{stringvalidator.RegexMatches(common.XPathRegex, common.XPathRegexMsg)}},
		"description":  schema.StringAttribute{Optional: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"mtu":          schema.Int64Attribute{Optional: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()}},
	}
}

func (t *tfModel) loadXmlData(ctx context.Context, x *xmlModel, diags *diag.Diagnostics) {
	if x == nil {
		return
	}

	t.Name = types.StringPointerValue(x.Name)
	t.Description = types.StringPointerValue(x.Description)
	t.Mtu = types.Int64PointerValue(x.Mtu)
}
