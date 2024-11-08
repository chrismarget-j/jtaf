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

var pointyBracketRegex = regexp.MustCompile("<[^<>]*>")

type PathElement struct {
	Name       string
	Attributes map[string]string
}

func (p PathElement) String() string {
	sb := new(strings.Builder)
	sb.WriteString("<")
	sb.WriteString(url.QueryEscape(p.Name))

	keys := make([]string, len(p.Attributes))
	var i int
	for key := range p.Attributes {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	for _, k := range keys {
		sb.WriteString(fmt.Sprintf(" %s=%s", url.QueryEscape(k), url.QueryEscape(p.Attributes[k])))
	}
	sb.WriteString(">")
	return sb.String()
}

func (p PathElement) NoAttrString() string {
	return "<" + url.QueryEscape(p.Name) + ">"
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
