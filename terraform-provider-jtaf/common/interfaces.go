package common

import "github.com/hashicorp/terraform-plugin-framework/attr"

type AttrTyper interface {
	AttrTypes() map[string]attr.Type
}

type XPathSetter interface {
	SetXPath(s string)
}
