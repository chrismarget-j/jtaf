package resourceinterfaceunit

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type xmlModelFamilyInet struct {
	ArpMaxCache *int64 `xml:"arp-max-cache,omitempty"`
}

func (x *xmlModelFamilyInet) loadTfData(ctx context.Context, tfObj types.Object, diags *diag.Diagnostics) {
	var tfData tfModelFamilyInet
	diags.Append(tfObj.As(ctx, &tfData, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return
	}

	x.ArpMaxCache = tfData.ArpMaxCache.ValueInt64Pointer()
}
