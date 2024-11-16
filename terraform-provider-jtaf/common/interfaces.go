// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package common

import "github.com/hashicorp/terraform-plugin-framework/attr"

type AttrTyper interface {
	AttrTypes() map[string]attr.Type
}

type XPathSetter interface {
	SetXPath(s string)
}
