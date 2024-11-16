// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package interfacesinterfaceunit

import (
	"path"

	interfacesinterface "github.com/chrismarget-j/jtaf/terraform-provider-jtaf/generated/resources/interfaces_interface"
)

const (
	xPathPrefix = ""
	xPathBase   = "unit"
)

var (
	xPathParent = interfacesinterface.RawXPath
	RawXPath    = path.Join(xPathParent, xPathPrefix, xPathBase)
)
