package resourceinterfaceunit

import "encoding/xml"

type xmlModel struct {
	XMLName           xml.Name        `xml:"unit"`
	Name              string          `xml:"name"`
	Description       *string         `xml:"description,omitempty"`
	NativeInnerVlanId *int64          `xml:"native-inner-vlan-id,omitempty"`
	Family            *xmlModelFamily `xml:"family,omitempty"`
}

type xmlModelFamily struct {
	Inet *xmlModelFamilyInet `xml:"inet,omitempty"`
}

type xmlModelFamilyInet struct {
	ArpMaxCache *int64 `xml:"arp-max-cache,omitempty"`
}
