package resourceinterfaceunit

import (
	"context"
	"encoding/xml"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type xmlModel struct {
	XMLName           xml.Name        `xml:"unit"`
	Name              *string         `xml:"name"`
	Description       *string         `xml:"description,omitempty"`
	NativeInnerVlanId *int64          `xml:"native-inner-vlan-id,omitempty"`
	Family            *xmlModelFamily `xml:"family,omitempty"`
}

func (x *xmlModel) loadTfData(ctx context.Context, tfObj types.Object, diags *diag.Diagnostics) {
	var tfData tfModel
	diags.Append(tfObj.As(ctx, &tfData, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return
	}

	x.Name = tfData.Name.ValueStringPointer()
	x.Description = tfData.Description.ValueStringPointer()
	x.NativeInnerVlanId = tfData.NativeInnerVlanId.ValueInt64Pointer()
	if !tfData.Family.IsNull() {
		x.Family = new(xmlModelFamily)
		x.Family.loadTfData(ctx, tfData.Family, diags)
	}
}
