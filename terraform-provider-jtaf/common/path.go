package common

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const upperHex = "0123456789ABCDEF"

var pointyBracketRegex = regexp.MustCompile("<[^<>]*>")

type PathElement struct {
	Name       string
	Attributes map[string]string
}

func (p PathElement) String() string {
	sb := new(strings.Builder)
	sb.WriteString("<")
	sb.WriteString(escape(p.Name))

	keys := make([]string, len(p.Attributes))
	var i int
	for key := range p.Attributes {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	for _, k := range keys {
		sb.WriteString(fmt.Sprintf(" %s=%s", escape(k), escape(p.Attributes[k])))
	}
	sb.WriteString(">")
	return sb.String()
}

func (p PathElement) NoAttrString() string {
	return "<" + escape(p.Name) + ">"
}

func NewPathElement(name string, attributes map[string]string) PathElement {
	return PathElement{
		Name:       name,
		Attributes: attributes,
	}
}

func NewPathElementFromString(s string) (PathElement, error) {
	var result PathElement

	if len(s) < 2 || s[0] != '<' || s[len(s)-1] != '>' {
		return result, fmt.Errorf("unparseable path element string: %q", s)
	}

	parts := strings.Fields(s[1 : len(s)-1])

	if len(parts) > 0 {
		var err error
		result.Name, err = url.QueryUnescape(parts[0])
		if err != nil {
			return result, fmt.Errorf("failed to parse name from %q - %w", s, err)
		}
		parts = parts[1:]
	}

	if len(parts) > 0 {
		result.Attributes = make(map[string]string)
	}

	for i, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return result, fmt.Errorf("unparseable attribute at index %d: %q", i, kv)
		}
		k, err := url.QueryUnescape(kv[0])
		if err != nil {
			return result, fmt.Errorf("failed to parse attribute key at index %d %q - %w", i, s, err)
		}
		v, err := url.QueryUnescape(kv[1])
		if err != nil {
			return result, fmt.Errorf("failed to parse attribute value at index %d %q - %w", i, s, err)
		}
		result.Attributes[k] = v
	}

	return result, nil
}

type Path []PathElement

func (p Path) String() string {
	sb := new(strings.Builder)
	for _, pe := range p {
		sb.WriteString(pe.String())
	}
	return sb.String()
}

func (p Path) NoAttrString() string {
	sb := new(strings.Builder)
	for _, pe := range p {
		sb.WriteString(pe.NoAttrString())
	}
	return sb.String()
}

func NewPathFromString(s string) (Path, error) {
	reIndexes := pointyBracketRegex.FindAllStringIndex(s, -1)

	var err error

	// split string into components
	result := make(Path, len(reIndexes))
	for i, component := range reIndexes {
		result[i], err = NewPathElementFromString(s[component[0]:component[1]])
		if err != nil {
			return nil, fmt.Errorf("failed while parsing path %q component at index %d - %w", s, i, err)
		}
	}

	return result, nil
}

func ParseParentPath(parentPath types.String, diags *diag.Diagnostics) Path {
	result, err := NewPathFromString(parentPath.ValueString())
	if err != nil {
		diags.AddAttributeError(
			path.Root("parent_path"),
			fmt.Sprintf("cannot parse parent path %s", parentPath),
			err.Error(),
		)
	}

	return result
}

// shouldEscape is based on the net/url function of the same name
func shouldEscape(c byte) bool {
	if c == '%' || //            //        0x25 must be escaped
		c == '+' || //           //        0x2b must be escaped
		('<' <= c && c <= '>') { // 0x3c - 0x3e must be escaped
		return true
	}

	if '!' <= c && c <= '~' { // remaining values from 0x21 - 0x7e are okay
		return false
	}

	return true // everything else must be escaped
}

// escape is based on the net/url function of the same name
func escape(s string) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			if c == ' ' {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	var buf [64]byte
	var t []byte

	required := len(s) + 2*hexCount
	if required <= len(buf) {
		t = buf[:required]
	} else {
		t = make([]byte, required)
	}

	if hexCount == 0 {
		copy(t, s)
		for i := 0; i < len(s); i++ {
			if s[i] == ' ' {
				t[i] = '+'
			}
		}
		return string(t)
	}

	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ':
			t[j] = '+'
			j++
		case shouldEscape(c):
			t[j] = '%'
			t[j+1] = upperHex[c>>4]
			t[j+2] = upperHex[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}

	return string(t)
}
