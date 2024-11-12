package resourceinterfaceunit

import (
	"context"
	"encoding/xml"
	"github.com/nemith/netconf"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	providerdata "github.com/chrismarget-j/jtaf/terraform-provider-jtaf/provider_data"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	plan.Id = common.XPathHash(plan.XPath)

	var x xmlModel
	x.loadTfData(ctx, common.ObjectValueFromAttrTyper(ctx, &plan, &resp.Diagnostics), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
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
	if b == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	var x xmlModel
	err := xml.Unmarshal(b, &x)
	if err != nil {
		resp.Diagnostics.AddError("cannot unmarshal config XML", err.Error())
		return
	}

	state.loadXmlData(ctx, &x, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	panic("not implemented")
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tfModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var x xmlModel
	x.XmlAttrOperation = netconf.DeleteConfig

	r.client.SetConfig(ctx, state.ParentXPath, x, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
}
