package resourceinterfaceunit

import (
	"context"
	"encoding/xml"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	providerdata "github.com/chrismarget-j/jtaf/terraform-provider-jtaf/provider_data"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	myType     = "unit"      // todo: generated
	parentType = "interface" // todo: generated
)

var _ resource.ResourceWithConfigure = (*Resource)(nil)

type Resource struct {
	client providerdata.Client
}

func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = common.JoinNonEmptyPartsWithUnderscores(req.ProviderTypeName, parentType, myType)
}

func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: (*tfModel)(nil).attributes()}
}

func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*providerdata.ResourceData).Client
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tfModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.XPath = common.AddPath(plan.ParentXPath, myType, plan.Name, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Id = plan.XPath

	var x xmlModel
	x.Name = plan.Name.ValueString()

	x.Description = plan.Description.ValueStringPointer()            // todo: generated
	x.NativeInnerVlanId = plan.NativeInnerVlanId.ValueInt64Pointer() // todo: generated

	if !plan.Family.IsNull() {
		x.Family = new(xmlModelFamily)
		var planFamily tfModelFamily
		resp.Diagnostics.Append(plan.Family.As(ctx, &planFamily, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		if !planFamily.Inet.IsNull() {
			x.Family.Inet = new(xmlModelFamilyInet)
			var planFamilyInet tfModelFamilyInet
			resp.Diagnostics.Append(planFamily.Inet.As(ctx, &planFamilyInet, basetypes.ObjectAsOptions{})...)
			if resp.Diagnostics.HasError() {
				return
			}

			x.Family.Inet.ArpMaxCache = planFamilyInet.ArpMaxCache.ValueInt64Pointer()
		}
	}

	r.client.SetConfig(ctx, plan.ParentXPath, x, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tfModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	b := r.client.GetConfig(ctx, state.XPath, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var x xmlModel
	err := xml.Unmarshal(b, &x)
	if err != nil {
		resp.Diagnostics.AddError("cannot unmarshal config XML", err.Error())
		return
	}

	state.Description = types.StringPointerValue(x.Description)            // todo: generated
	state.NativeInnerVlanId = types.Int64PointerValue(x.NativeInnerVlanId) // todo: generated

	if x.Family != nil {
		stateFamily := new(tfModelFamily)

		if x.Family.Inet != nil {
			stateFamilyInet := new(tfModelFamilyInet)

			if x.Family.Inet.ArpMaxCache != nil {
				stateFamilyInet.ArpMaxCache = types.Int64PointerValue(x.Family.Inet.ArpMaxCache)
			}

			stateFamily.Inet = common.ObjectValueFromWithDiags(ctx, (*tfModelFamilyInet)(nil).attrTypes(), stateFamilyInet, &resp.Diagnostics)
		}

		state.Family = common.ObjectValueFromWithDiags(ctx, (*tfModelFamily)().attrTypes(), stateFamily, &resp.Diagnostics)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO implement me
	panic("implement me")
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// TODO implement me
	panic("implement me")
}
