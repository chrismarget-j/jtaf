// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package interfacesinterfaceunitfamilyinetaddress

import (
	"context"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	Primary     types.Bool   `tfsdk:"primary"`
}

func (t *tfModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":           types.StringType,
		"xpath":        types.StringType,
		"name":         types.StringType,
		"parent_xpath": types.StringType,
		"primary":      types.BoolType,
	}
}

func (t *tfModel) attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"xpath":        schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"name":         schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"parent_xpath": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Validators: []validator.String{stringvalidator.RegexMatches(common.XPathRegex, common.XPathRegexMsg)}},
		"primary":      schema.BoolAttribute{Optional: true},
	}
}

func (t *tfModel) loadXmlData(ctx context.Context, x *xmlModel, diags *diag.Diagnostics) {
	if x == nil {
		return
	}

	t.Name = types.StringPointerValue(x.Name.ValuePointer())
	t.Primary = types.BoolPointerValue(x.Primary.ValuePointer())
}

func tfModelNull() types.Object {
	return types.ObjectNull((*tfModel)(nil).AttrTypes())
}
