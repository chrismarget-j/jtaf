// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package yang

import (
	"log"

	"github.com/openconfig/goyang/pkg/yang"
)

const (
	patternJunosRef = `<.*>|$.*`
	patternDefault  = `default`
	patternDerived  = `derived`
	patternAsterisk = `\*`
)

func TrimUnions(ye *yang.Entry) {
	trimUnions(ye, nil)
}

func trimUnions(ye *yang.Entry, xmlPath []string) {
	if ye.Type != nil && ye.Type.Kind == yang.Yunion {
		handleUnion(ye.Type, xmlPath)
	}

	for _, de := range ye.Dir {
		trimUnions(de, append(xmlPath, de.Name))
	}
}

func handleUnion(t *yang.YangType, xmlPath []string) {
	if t.Type == nil || t.Kind != yang.Yunion {
		return
	}

	for _, subType := range t.Type {
		if subType.Kind == yang.Yunion {
			handleUnion(subType, xmlPath) // this union contains a union!
		}
	}

	discardWellKnownStringConstraints(t, xmlPath)

	if t.Kind == yang.Yunion {
		s, err := TypeToString(t)
		if err != nil {
			log.Println("error - ", err.Error())
		} else {
			log.Println("unhandled Union: ", s)
		}
	}
}

func discardWellKnownStringConstraints(t *yang.YangType, xmlPath []string) bool {
	if t.Kind != yang.Yunion {
		return false
	}

	if len(t.Type) != 2 {
		return false
	}

	var wellKnownStringConstraintType *yang.YangType
	var otherType *yang.YangType
	for _, utm := range t.Type { // loop over union type members
		if utm.Kind == yang.Ystring && len(utm.Pattern) == 1 {
			switch utm.Pattern[0] {
			case patternAsterisk:
				wellKnownStringConstraintType = utm
				continue
			case patternDefault:
				wellKnownStringConstraintType = utm
				continue
			case patternDerived:
				wellKnownStringConstraintType = utm
				continue
			case patternJunosRef:
				wellKnownStringConstraintType = utm
				continue
			}
		}

		otherType = utm
	}

	if wellKnownStringConstraintType != nil && otherType != nil {
		*t = *otherType // replace the union with the "other" type found in the union
		return true
	}

	return false
}
