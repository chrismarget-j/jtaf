package interfacesinterfaceunit

import (
	"context"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common/values"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/nemith/netconf"
)

type xmlModelFamilyInet struct {
	Operation   *netconf.MergeStrategy `xml:"operation,attr,omitempty"`
	ArpMaxCache *values.XmlInt64       `xml:"arp-max-cache,omitempty"`
}

func newXmlModelFamilyInet(ctx context.Context, v types.Object, diags *diag.Diagnostics) *xmlModelFamilyInet {
	var tfData tfModelFamilyInet
	if !v.IsNull() {
		diags.Append(v.As(ctx, &tfData, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil
		}
	}

	return &xmlModelFamilyInet{
		Operation:   common.MergeStrategyFromValue(v),
		ArpMaxCache: values.NewXmlInt64(ctx, tfData.ArpMaxCache, diags),
	}
}

func delXmlModelFamilyInet(ctx context.Context, v types.Object, diags *diag.Diagnostics) *xmlModelFamilyInet {
	var tfData tfModelFamilyInet
	if !v.IsNull() {
		diags.Append(v.As(ctx, &tfData, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil
		}
	}

	return &xmlModelFamilyInet{
		Operation: common.RemoveConfig,
	}
}
