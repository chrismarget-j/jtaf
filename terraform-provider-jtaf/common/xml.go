package common

import (
	"strings"
)

// XmlWrappersFromPath returns header xml, a prefix to be used with nested xml, and suffix xml
func XmlWrappersFromPath(path Path, prefix, indent string) (string, string, string) {
	header := new(strings.Builder)
	for i, pe := range path {
		header.WriteString(prefix + strings.Repeat(indent, i) + pe.NoAttrString() + "\n")
		if name, ok := pe.Attributes["name"]; ok {
			header.WriteString(prefix + strings.Repeat(indent, i+1) + "<name>" + name + "</name>\n")
		}
	}

	footer := new(strings.Builder)
	for i := len(path) - 1; i >= 0; i-- {
		pe := path[i]
		penas := pe.NoAttrString()
		penas = penas[:1] + "/" + penas[1:]
		footer.WriteString(prefix + strings.Repeat(indent, i) + penas + "\n")
	}

	return header.String(), prefix + strings.Repeat(indent, len(path)), footer.String()
}
