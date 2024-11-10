package common

import (
	"fmt"
	"strings"

	"github.com/ChrisTrenkamp/xsel"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// XmlWrappersFromPathGrammar returns header xml, a prefix to be used with nested xml, and suffix xml
func XmlWrappersFromPathGrammar(g xsel.Grammar, prefix, indent string, diags *diag.Diagnostics) (string, string, string) {
	type element struct {
		label      string
		attributes map[string]string
	}

	var elements []element
	var i int
loop:
	for {
		var m map[string]string
		var err error

		s := g.GetStringExtents(i, i+1)
		switch s {
		case "":
			break loop
		case "/":
			i++
			continue
		case "[":
			m, i, err = xPathGrammarAttributeAt(g, i)
			if err != nil {
				diags.AddError("while prepping XML header/footer - %w", err.Error())
			}
			if len(elements) == 0 {
				diags.AddError("found attribute block before label", fmt.Sprintf("query: %q", g.GetString()))
				return "", "", ""
			}
			elements[len(elements)-1].attributes = m
		default:
			i++
			elements = append(elements, element{label: s})
		}
	}

	header := new(strings.Builder)
	for i, pe := range elements {
		header.WriteString(prefix + strings.Repeat(indent, i) + "<" + pe.label + ">\n")
		if name, ok := pe.attributes["name"]; ok {
			header.WriteString(prefix + strings.Repeat(indent, i+1) + "<name>" + name + "</name>\n")
		}
	}

	footer := new(strings.Builder)
	for i := len(elements) - 1; i >= 0; i-- {
		footer.WriteString(prefix + strings.Repeat(indent, i) + "</" + elements[i].label + ">\n")
	}

	return header.String(), prefix + strings.Repeat(indent, len(elements)), footer.String()
}
