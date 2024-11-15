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
