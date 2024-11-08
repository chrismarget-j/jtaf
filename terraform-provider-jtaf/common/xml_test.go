package common_test

import (
	"testing"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/stretchr/testify/require"
)

//func TestXmlNamePath(t *testing.T) {
//	type testCase struct {
//		data     string
//		expected string
//	}
//
//	testCases := map[string]testCase{
//		"interface": {
//			data:     "<interfaces><interface><name>ge-0/0/0</name></interface></interfaces></configuration>",
//			expected: "<interfaces><interface name=ge-0/0/0>",
//		},
//	}
//
//	for tName, tCase := range testCases {
//		t.Run(tName, func(t *testing.T) {
//			t.Parallel()
//
//			var diags diag.Diagnostics
//
//			actual := common.XmlNamePath([]byte(tCase.data), &diags)
//			require.Nilf(t, diags, "diags not nil - %d errors and %d warnings", len(diags.Errors()), len(diags.Warnings()))
//			require.Equal(t, tCase.expected, actual)
//		})
//	}
//}

//func TestAddRemoveBrackets(t *testing.T) {
//	type testCase struct {
//		data             string
//		dataWithBrackets string
//	}
//
//	testCases := map[string]testCase{
//		"foo": {
//			data:             "foo",
//			dataWithBrackets: "<foo>",
//		},
//		"empty": {
//			data:             "",
//			dataWithBrackets: "<>",
//		},
//	}
//
//	for tName, tCase := range testCases {
//		t.Run(tName, func(t *testing.T) {
//			t.Parallel()
//
//			data := tCase.data
//			common.AddBrackets(&data)
//			require.Equal(t, tCase.dataWithBrackets, data)
//			common.RemoveBrackets(&data)
//			require.Equal(t, tCase.data, data)
//		})
//	}
//}

//func TestXmlWrappersFromNamePath(t *testing.T) {
//	type result struct {
//		prefix string
//		indent string
//		suffix string
//	}
//
//	type testCase struct {
//		data     string
//		prefix   string
//		indent   string
//		expected result
//	}
//
//	testCases := map[string]testCase{
//		"a": {
//			data:   "<interfaces><interface name=ge-0/0/0>",
//			prefix: "",
//			indent: "  ",
//			expected: result{
//				prefix: "<interfaces>\n  <interface>\n    <name>ge-0/0/0</name>\n",
//				indent: "    ",
//				suffix: "  </interface>\n</interfaces>\n",
//			},
//		},
//		"b": {
//			data:   "<interfaces><interface name=ge-0/0/0 bar=baz>",
//			prefix: "P",
//			indent: "I",
//			expected: result{
//				prefix: "P<interfaces>\nPI<interface bar=baz>\nPII<name>ge-0/0/0</name>\n",
//				indent: "PII",
//				suffix: "PI</interface>\nP</interfaces>\n",
//			},
//		},
//	}
//
//	for tName, tCase := range testCases {
//		t.Run(tName, func(t *testing.T) {
//			t.Parallel()
//
//			prefix, indent, suffix := common.XmlWrappersFromNamePath(tCase.data, tCase.prefix, tCase.indent)
//			require.Equal(t, tCase.expected.prefix, prefix)
//			require.Equal(t, tCase.expected.indent, indent)
//			require.Equal(t, tCase.expected.suffix, suffix)
//		})
//	}
//}

func TestXmlWrappersFromPath(t *testing.T) {
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
			pathString:     "<foo><bar>",
			expectedHeader: "<foo>\n<bar>\n",
			expectedFooter: "</bar>\n</foo>\n",
		},
		"b": {
			pathString:     "<foo><bar><baz>",
			prefix:         "P",
			indent:         "I",
			expectedHeader: "P<foo>\nPI<bar>\nPII<baz>\n",
			expectedPrefix: "PIII",
			expectedFooter: "PII</baz>\nPI</bar>\nP</foo>\n",
		},
		"c": {
			pathString:     "<foo name=F><bar X=X><baz name=B>",
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

			path, err := common.NewPathFromString(tCase.pathString)
			require.NoError(t, err)
			header, prefix, footer := common.XmlWrappersFromPath(path, tCase.prefix, tCase.indent)
			require.Equal(t, tCase.expectedHeader, header)
			require.Equal(t, tCase.expectedPrefix, prefix)
			require.Equal(t, tCase.expectedFooter, footer)
		})
	}
}
