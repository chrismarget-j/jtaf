// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package providerdata

import (
	"github.com/nemith/netconf"
)

type ResourceData struct {
	Client Client
}

func NewResourceData(ns *netconf.Session) *ResourceData {
	return &ResourceData{
		Client: newClient(ns),
	}
}
