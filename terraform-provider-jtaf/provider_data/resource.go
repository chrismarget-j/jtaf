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
