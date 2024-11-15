package resourceinterfaceunit

import (
	"context"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/nemith/netconf"
)

type xmlModelFamily struct {
	Operation *netconf.MergeStrategy `xml:"operation,attr,omitempty"`
	Inet      *xmlModelFamilyInet    `xml:"inet,omitempty"`
}

func newXmlModelFamily(ctx context.Context, v types.Object, diags *diag.Diagnostics) *xmlModelFamily {
	var tfData tfModelFamily
	if !v.IsNull() {
		diags.Append(v.As(ctx, &tfData, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil
		}
	}

	return &xmlModelFamily{
		Operation: common.MergeStrategyFromValue(v),
		Inet:      newXmlModelFamilyInet(ctx, tfData.Inet, diags),
	}
}

func delXmlModelFamily(ctx context.Context, v types.Object, diags *diag.Diagnostics) *xmlModelFamily {
	var tfData tfModelFamily
	if !v.IsNull() {
		diags.Append(v.As(ctx, &tfData, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil
		}
	}

	return &xmlModelFamily{
		Operation: common.RemoveConfig,
	}
}
