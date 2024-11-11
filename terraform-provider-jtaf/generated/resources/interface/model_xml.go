package resourceinterface

import (
	"context"
	"encoding/xml"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type xmlModel struct {
	XMLName     xml.Name `xml:"interface"`
	Name        *string  `xml:"name"`
	Description *string  `xml:"description,omitempty"`
	Mtu         *int64   `xml:"mtu,omitempty"`
}

func (x *xmlModel) loadTfData(ctx context.Context, tfObj types.Object, diags *diag.Diagnostics) {
	var tfData tfModel
	diags.Append(tfObj.As(ctx, &tfData, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return
	}

	x.Name = tfData.Name.ValueStringPointer()
	x.Description = tfData.Description.ValueStringPointer()
	x.Mtu = tfData.Mtu.ValueInt64Pointer()
}
