// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package interfacesinterfaceunitfamilyinetaddress

import (
	"path"

	interfacesinterfaceunit "github.com/chrismarget-j/jtaf/terraform-provider-jtaf/generated/resources/interfaces_interface_unit"
)

const (
	xPathPrefix = "family/inet"
	xPathBase   = "address"
)

var (
	xPathParent = interfacesinterfaceunit.RawXPath
	RawXPath    = path.Join(xPathParent, xPathPrefix, xPathBase)
)
