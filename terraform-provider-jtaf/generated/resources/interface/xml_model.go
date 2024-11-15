package resourceinterface

import (
	"context"
	"encoding/xml"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common/values"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/nemith/netconf"
)

type xmlModel struct {
	XMLName     xml.Name               `xml:"interface"`
	Operation   *netconf.MergeStrategy `xml:"operation,attr,omitempty"`
	Name        *values.XmlString      `xml:"name"`
	Description *values.XmlString      `xml:"description,omitempty"`
	Mtu         *values.XmlInt64       `xml:"mtu,omitempty"`
}

func newXmlModel(ctx context.Context, v types.Object, diags *diag.Diagnostics) *xmlModel {
	var tfData tfModel
	if !v.IsNull() {
		diags.Append(v.As(ctx, &tfData, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil
		}
	}

	return &xmlModel{
		Name:        values.NewXmlString(ctx, tfData.Name, diags),
		Description: values.NewXmlString(ctx, tfData.Description, diags),
		Mtu:         values.NewXmlInt64(ctx, tfData.Mtu, diags),
	}
}
