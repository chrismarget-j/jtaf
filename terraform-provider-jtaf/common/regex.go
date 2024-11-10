package common

import "regexp"

const (
	xPathRegex    = `^\/.*[^/]$`
	XPathRegexMsg = `value must begin with "/" and not end with "/"`
)

var XPathRegex = regexp.MustCompile(xPathRegex)
