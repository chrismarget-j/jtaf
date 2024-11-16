// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package common

import "regexp"

const (
	xPathRegex    = `^\/.*[^/]$`
	XPathRegexMsg = `value must begin with "/" and not end with "/"`
)

var XPathRegex = regexp.MustCompile(xPathRegex)
