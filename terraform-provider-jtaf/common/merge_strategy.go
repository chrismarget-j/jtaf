// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/nemith/netconf"
)

var (
	RemoveConfig = ToPtr(netconf.RemoveConfig)
	MergeConfig  = ToPtr(netconf.MergeConfig)
)

func MergeStrategyFromValue(v attr.Value) *netconf.MergeStrategy {
	if v.IsNull() {
		return RemoveConfig
	}

	return MergeConfig
}
