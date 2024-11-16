// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package common_test

import (
	"testing"

	"github.com/ChrisTrenkamp/xsel"
	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/require"
)

func TestXmlWrappersFromPathGrammar(t *testing.T) {
	type testCase struct {
		pathString     string
		prefix         string
		indent         string
		expectedHeader string
		expectedPrefix string
		expectedFooter string
	}

	testCases := map[string]testCase{
		"a": {
			pathString:     "/foo/bar",
			expectedHeader: "<foo>\n<bar>\n",
			expectedFooter: "</bar>\n</foo>\n",
		},
		"b": {
			pathString:     "/foo/bar/baz",
			prefix:         "P",
			indent:         "I",
			expectedHeader: "P<foo>\nPI<bar>\nPII<baz>\n",
			expectedPrefix: "PIII",
			expectedFooter: "PII</baz>\nPI</bar>\nP</foo>\n",
		},
		"c": {
			pathString:     "/foo[name='F']/bar[x='X']/baz[name='B']",
			prefix:         "P",
			indent:         "I",
			expectedHeader: "P<foo>\nPI<name>F</name>\nPI<bar>\nPII<baz>\nPIII<name>B</name>\n",
			expectedPrefix: "PIII",
			expectedFooter: "PII</baz>\nPI</bar>\nP</foo>\n",
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			g, err := xsel.BuildExpr(tCase.pathString)
			require.NoError(t, err)

			var diags *diag.Diagnostics

			header, prefix, footer := common.XmlWrappersFromPathGrammar(g, tCase.prefix, tCase.indent, diags)
			require.Nil(t, diags)
			require.Equal(t, tCase.expectedHeader, header)
			require.Equal(t, tCase.expectedPrefix, prefix)
			require.Equal(t, tCase.expectedFooter, footer)
		})
	}
}
