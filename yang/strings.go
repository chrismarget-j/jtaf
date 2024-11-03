package yang

import (
	"fmt"
	"strings"

	"github.com/openconfig/goyang/pkg/yang"
)

const baseIndent = "  "

func TypeToString(yt *yang.YangType) (string, error) {
	return typeToString(yt, baseIndent)
}

func typeToString(yt *yang.YangType, indent string) (string, error) {
	switch yt.Kind {
	case yang.Yempty:
		return emptyKindToString(yt, indent)
	case yang.Yenum:
		return enumKindToString(yt, indent)
	case yang.Yuint8, yang.Yuint16, yang.Yuint32, yang.Yuint64, yang.Yint8, yang.Yint16, yang.Yint32, yang.Yint64:
		return intKindToString(yt, indent)
	case yang.Ystring:
		return stringKindToString(yt, indent)
	case yang.Yunion:
		return unionKindToString(yt, indent)
	default:
		return "", fmt.Errorf("unhandled type kind %s (%d)", yang.TypeKindToName[yt.Kind], yt.Kind)
	}
}

//func decimal64KindToString(yt *yang.YangType, indent string) (string, error) {
//	switch yt.Kind {
//	case yang.Ydecimal64:
//	default:
//		return "", fmt.Errorf("expected an decimal64 kind, got %s (%d)", yang.TypeKindToName[yt.Kind], yt.Kind)
//	}
//
//	var rnge string
//	if len(yt.Range) > 0 {
//		rnge = "\n" + indent + "  Range: " + yt.Range.String()
//	}
//
//	var unit string
//	if yt.Units != "" {
//		unit = fmt.Sprintf("\n"+indent+"  Units: %s", yt.Units)
//	}
//
//	return indent + yang.TypeKindToName[yt.Kind] + unit + rnge, nil
//}

func emptyKindToString(yt *yang.YangType, indent string) (string, error) {
	switch yt.Kind {
	case yang.Yempty:
	default:
		return "", fmt.Errorf("expected an empty kind, got %s (%d)", yang.TypeKindToName[yt.Kind], yt.Kind)
	}

	return indent + "<" + yang.TypeKindToName[yt.Kind] + ">", nil
}

func enumKindToString(yt *yang.YangType, indent string) (string, error) {
	switch yt.Kind {
	case yang.Yenum:
	default:
		return "", fmt.Errorf("expected an enum kind, got %s (%d)", yang.TypeKindToName[yt.Kind], yt.Kind)
	}

	// var vals string
	// if yt.E

	return indent + "<" + yang.TypeKindToName[yt.Kind] + ">", nil
}

func intKindToString(yt *yang.YangType, indent string) (string, error) {
	switch yt.Kind {
	case yang.Yuint8, yang.Yuint16, yang.Yuint32, yang.Yuint64:
	case yang.Yint8, yang.Yint16, yang.Yint32, yang.Yint64:
	default:
		return "", fmt.Errorf("expected an integer kind, got %s (%d)", yang.TypeKindToName[yt.Kind], yt.Kind)
	}

	var rnge string
	if len(yt.Range) > 0 {
		rnge = "  Range: " + yt.Range.String()
	}

	var unit string
	if yt.Units != "" {
		unit = "  Units: %s" + yt.Units
	}

	return indent + "<" + yang.TypeKindToName[yt.Kind] + ">" + unit + rnge, nil
}

func stringKindToString(yt *yang.YangType, indent string) (string, error) {
	switch yt.Kind {
	case yang.Ystring:
	default:
		return "", fmt.Errorf("expected a string kind, got %s (%d)", yang.TypeKindToName[yt.Kind], yt.Kind)
	}

	var pattern string
	if len(yt.Pattern) > 0 {
		pattern = fmt.Sprintf("  Pattern: %s", yt.Pattern)
	}

	return indent + "<" + yang.TypeKindToName[yt.Kind] + ">" + pattern, nil
}

func unionKindToString(yt *yang.YangType, indent string) (string, error) {
	switch yt.Kind {
	case yang.Yunion:
	default:
		return "", fmt.Errorf("expected an union kind, got %s (%d)", yang.TypeKindToName[yt.Kind], yt.Kind)
	}

	sb := new(strings.Builder)
	if len(yt.Type) > 0 {
		for i, t := range yt.Type {
			s, err := typeToString(t, indent+baseIndent)
			if err != nil {
				return "", fmt.Errorf("while string-ing union element %d - %w", i, err)
			}

			sb.WriteString("\n" + indent + s)
		}
	}
	return indent + "<" + yang.TypeKindToName[yt.Kind] + ">" + sb.String(), nil
}
