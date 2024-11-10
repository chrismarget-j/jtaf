package resourceinterfaceunit

import (
	"context"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	providerdata "github.com/chrismarget-j/jtaf/terraform-provider-jtaf/provider_data"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

const (
	myType     = "unit"
	parentType = "interface"
)

var _ resource.ResourceWithConfigure = (*Resource)(nil)

type Resource struct {
	client providerdata.Client
}

func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + parentType + "_" + myType
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
	x.Description = plan.Description.ValueStringPointer()
	x.NativeInnerVlanId = plan.NativeInnerVlanId.ValueInt64Pointer()

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

	//myPath, err := common.NewPathFromString(state.Path.ValueString())
	//if err != nil {
	//	resp.Diagnostics.AddError(fmt.Sprintf("failed to parse resource XML path %s during Read", state.Path), err.Error())
	//	return
	//}
	//
	//b := r.client.GetConfig(ctx, myPath, &resp.Diagnostics)
	//if resp.Diagnostics.HasError() {
	//	return
	//}
	//_ = b
	//
	// TODO implement me
	panic("implement me")
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO implement me
	panic("implement me")
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// TODO implement me
	panic("implement me")
}
