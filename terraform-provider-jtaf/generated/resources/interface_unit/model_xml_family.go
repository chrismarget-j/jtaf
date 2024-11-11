package resourceinterfaceunit

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type xmlModelFamily struct {
	Inet *xmlModelFamilyInet `xml:"inet,omitempty"`
}

func (x *xmlModelFamily) loadTfData(ctx context.Context, tfObj types.Object, diags *diag.Diagnostics) {
	var tfData tfModelFamily
	diags.Append(tfObj.As(ctx, &tfData, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return
	}

	if !tfData.Inet.IsNull() {
		x.Inet = new(xmlModelFamilyInet)
		x.Inet.loadTfData(ctx, tfData.Inet, diags)
	}
}
