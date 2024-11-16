// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package generated

import (
	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/generated/resources/interfaces_interface"
	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/generated/resources/interfaces_interface_unit"
	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/generated/resources/interfaces_interface_unit_family_inet_address"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var Resources = []func() resource.Resource{
	func() resource.Resource { return &interfacesinterface.Resource{} },
	func() resource.Resource { return &interfacesinterfaceunit.Resource{} },
	func() resource.Resource { return &interfacesinterfaceunitfamilyinetaddress.Resource{} },
}
