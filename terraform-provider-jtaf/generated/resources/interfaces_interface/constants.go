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
