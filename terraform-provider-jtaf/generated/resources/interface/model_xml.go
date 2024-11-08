package resourceinterface

import "encoding/xml"

type xmlModel struct {
	XMLName     xml.Name `xml:"interface"`
	Name        string   `xml:"name"`
	Description *string  `xml:"description,omitempty"`
	Mtu         *int64   `xml:"mtu,omitempty"`
}
