// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package common_test

import (
	"testing"

	"github.com/ChrisTrenkamp/xsel"
	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/stretchr/testify/require"
)

func TestEncodeExpression(t *testing.T) {
	type testCase struct {
		d string
		e string
	}

	testCases := map[string]testCase{
		"empty": {
			d: "",
			e: `""`,
		},
		"a": {
			d: "a",
			e: `"a"`,
		},
		"ge-0/0/0": {
			d: "ge-0/0/0",
			e: `"ge-0/0/0"`,
		},
		//"rqe": {
		//	d: "Ralph's \"Quote Emporium\"",
		//	e: `concat('Ralph',"'",'s "Quote Emporium"')`,
		//},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			a, err := common.XPathQuoteAttr(tCase.d)
			require.NoError(t, err)
			require.Equal(t, tCase.e, a)
		})
	}
}

func TestXPathGrammarAttributeAt(t *testing.T) {
	type testCase struct {
		xpath   string
		left    int
		newLeft int
		expMap  map[string]string
		expErr  string
	}

	testCases := map[string]testCase{
		"a": {
			xpath:  `/`,
			expErr: `expected token "["`,
		},
		"b": {
			xpath:   `/foo[name="bar"]`,
			left:    2,
			newLeft: 7,
			expMap:  map[string]string{"name": "bar"},
		},
		"c": {
			xpath:   `/foo[a="A"][b="B"]`,
			left:    2,
			newLeft: 12,
			expMap:  map[string]string{"a": "A", "b": "B"},
		},
		"d": {
			xpath:   `/foo[a="A"][b="B"]/bar[c="C"]`,
			left:    2,
			newLeft: 12,
			expMap:  map[string]string{"a": "A", "b": "B"},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			g, err := xsel.BuildExpr(tCase.xpath)
			require.NoError(t, err)

			actual, newLeft, err := common.XPathGrammarAttributeAt(g, tCase.left)
			if tCase.expErr == "" {
				require.NoError(t, err)
				if tCase.newLeft > 0 {
					require.Equal(t, tCase.newLeft, newLeft)
				}
				require.Equal(t, tCase.expMap, actual)
			} else {
				require.ErrorContains(t, err, tCase.expErr)
			}
		})
	}
}
