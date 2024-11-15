package values

import (
	"context"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nemith/netconf"
)

type XmlBool struct {
	Operation *netconf.MergeStrategy `xml:"operation,attr,omitempty"`
}

func (o *XmlBool) ValuePointer() *bool {
	if o == nil {
		return nil
	}

	return common.ToPtr(true)
}

func NewXmlBool(_ context.Context, v types.Bool, _ *diag.Diagnostics) *XmlBool {
	return &XmlBool{
		Operation: common.MergeStrategyFromValue(v),
	}
}
