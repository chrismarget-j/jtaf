package resourceinterfaceunit

import (
	"context"
	"encoding/xml"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	providerdata "github.com/chrismarget-j/jtaf/terraform-provider-jtaf/provider_data"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nemith/netconf"
)

const myType = "unit"

var _ resource.ResourceWithConfigure = (*Resource)(nil)

type Resource struct {
	session *netconf.Session
}

func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interface_unit"
}

func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: (*tfModel)(nil).attributes()}
}

func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.session = req.ProviderData.(*providerdata.ResourceData).Session
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tfModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentPath := common.ParseParentPath(plan.ParentPath, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	myPath := append(parentPath, common.NewPathElement(myType, map[string]string{"name": plan.Name.ValueString()}))

	plan.Id = types.StringValue(myPath.String())
	plan.Path = types.StringValue(myPath.String())

	var x xmlModel
	x.Name = plan.Name.ValueString()
	x.Description = plan.Description.ValueStringPointer()
	x.NativeInnerVlanId = plan.NativeInnerVlanId.ValueInt64Pointer()

	header, prefix, footer := common.XmlWrappersFromPath(parentPath, "", common.XmlIndent)
	xmlBytes, err := xml.MarshalIndent(x, prefix, common.XmlIndent)
	if err != nil {
		resp.Diagnostics.AddError("failed marshaling config xml", err.Error())
		return
	}

	payload := header + string(xmlBytes) + "\n" + footer

	err = r.session.EditConfig(ctx, netconf.Candidate, []byte(payload))
	if err != nil {
		resp.Diagnostics.AddError("failed while editing device config", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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
