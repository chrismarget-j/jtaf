// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package interfacesinterface

import (
	"path"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
)

const (
	xPathPrefix = "interfaces"
	xPathBase   = "interface"
)

var (
	xPathParent = common.XPathRoot
	RawXPath    = path.Join(xPathParent, xPathPrefix, xPathBase)
)
