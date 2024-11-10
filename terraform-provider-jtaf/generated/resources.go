package generated

import (
	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/generated/resources/interface"
	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/generated/resources/interface_unit"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var Resources = []func() resource.Resource{
	func() resource.Resource { return &resourceinterface.Resource{} },
	func() resource.Resource { return &resourceinterfaceunit.Resource{} },
}
