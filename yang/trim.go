package yang

import (
	"github.com/openconfig/goyang/pkg/yang"
)

// TrimToConfig trims the *yang.Entry tree down to only those elements required by
// configPaths. Each element of configPaths is a slice of config XML path elements:
//
//	[][]string{
//	  {"configuration", "interfaces", "interface", "unit", "family", "inet"}
//	  {"configuration", "protocols", "ospf", "area", "interace", "passive"}
//	}
func TrimToConfig(ye *yang.Entry, configPaths [][]string) {
	var nextPaths [][]string // used in a recursive call to this function

	// trim unnecessary top level items from this entry
	for name, entry := range ye.Dir {
		var required bool
		for _, cp := range configPaths {
			if len(cp) > 0 && cp[0] == name {
				required = true

				// chop the head from this path and use the remainder on recursion
				if len(cp) > 1 {
					nextPaths = append(nextPaths, cp[1:])
				}
			}
		}

		if !required { // did we find a reason to let this entry remain?
			delete(ye.Dir, name)
			continue
		}

		// recurse if any path elements remain
		if len(nextPaths) > 0 {
			TrimToConfig(entry, nextPaths)
		}
	}
}
