package common

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/ChrisTrenkamp/xsel"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func XPathQuoteAttr(s string) (string, error) {
	if !strings.Contains(s, `"`) {
		return `"` + s + `"`, nil
	}

	if !strings.Contains(s, "'") {
		return `'` + s + `'`, nil
	}

	return "", fmt.Errorf("cowardly refusing to quote a string containing both single and double quotes")
}

func XPathUnquoteAttr(s string) (string, error) {
	if strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) {
		return s[1 : len(s)-1], nil
	}

	if strings.HasPrefix(s, `'`) && strings.HasSuffix(s, `'`) {
		return s[1 : len(s)-1], nil
	}

	return "", fmt.Errorf("don't know how to un-quote %q", s)
}

func AddPath(parent types.String, path string, name types.String, diags *diag.Diagnostics) types.String {
	encoded, err := XPathQuoteAttr(name.ValueString())
	if err != nil {
		diags.AddError(fmt.Sprintf("while encoding path element - %s", name.String()), err.Error())
	}

	return types.StringValue(
		fmt.Sprintf("%s/%s[name=%s]", parent.ValueString(), path, encoded),
	)
}

func xPathGrammarAttributeAt(g xsel.Grammar, left int) (map[string]string, int, error) {
	s := g.GetStringExtents(left, left+1)
	if s != "[" {
		return nil, 0, fmt.Errorf("expected token %q, got %q", "[", s)
	}

	result := make(map[string]string)

	var key *string
	for {
		s := g.GetStringExtents(left, left+1)
		switch {
		case s == "=" && key == nil:
			return nil, 0, fmt.Errorf("unexpected %q while parsing %q at token %d", s, g.GetString(), left)
		case s == "[" || s == "=" || s == "]":
			left++
		case s == "" && key != nil:
			return nil, 0, fmt.Errorf("EOF while parsing %q at token %d", g.GetString(), left)
		case s == "" || s == "/":
			return result, left, nil
		default:
			if key == nil {
				key = &s
				left++
			} else {
				unQuoted, err := XPathUnquoteAttr(s)
				if err != nil {
					return nil, 0, fmt.Errorf("while un-quoting %q - %w", s, err)
				}

				result[*key] = unQuoted
				key = nil
				left++
			}
		}
	}
}

func XPathHash(xpath types.String) types.String {
	h := sha256.New()
	h.Write([]byte(xpath.ValueString()))
	return types.StringValue(fmt.Sprintf("%x", h.Sum(nil))[:12])
}
