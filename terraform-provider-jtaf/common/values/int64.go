// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package values

import (
	"context"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nemith/netconf"
)

type XmlInt64 struct {
	Operation *netconf.MergeStrategy `xml:"operation,attr,omitempty"`
	Value     *int64                 `xml:",chardata"`
}

func (o *XmlInt64) ValuePointer() *int64 {
	if o == nil {
		return nil
	}

	return o.Value
}

func NewXmlInt64(_ context.Context, v types.Int64, _ *diag.Diagnostics) *XmlInt64 {
	return &XmlInt64{
		Operation: common.MergeStrategyFromValue(v),
		Value:     v.ValueInt64Pointer(),
	}
}
