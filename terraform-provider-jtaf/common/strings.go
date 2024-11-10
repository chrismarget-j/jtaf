package common

import "strings"

func JoinNonEmptyPartsWithUnderscores(parts ...string) string {
	sb := new(strings.Builder)
	var notFirst bool
	for _, part := range parts {
		if part != "" {
			if notFirst {
				sb.WriteString("_" + part)
			} else {
				sb.WriteString(part)
			}
			notFirst = true
		}
	}

	return sb.String()
}
