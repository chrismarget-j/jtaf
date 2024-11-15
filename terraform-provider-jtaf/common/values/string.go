package values

import (
	"context"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nemith/netconf"
)

type XmlString struct {
	Operation *netconf.MergeStrategy `xml:"operation,attr,omitempty"`
	Value     *string                `xml:",chardata"`
}

func (o *XmlString) ValuePointer() *string {
	if o == nil {
		return nil
	}

	return o.Value
}

func NewXmlString(_ context.Context, v types.String, _ *diag.Diagnostics) *XmlString {
	return &XmlString{
		Operation: common.MergeStrategyFromValue(v),
		Value:     v.ValueStringPointer(),
	}
}
