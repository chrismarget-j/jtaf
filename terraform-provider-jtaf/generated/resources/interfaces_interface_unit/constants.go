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
