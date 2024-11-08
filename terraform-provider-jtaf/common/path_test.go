package common_test

import (
	"testing"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/stretchr/testify/require"
)

func TestPathElementStrings(t *testing.T) {
	type testCase struct {
		source               string
		expectedObj          common.PathElement
		expectedNoAttrString string
	}

	testCases := map[string]testCase{
		"a": {
			source:               "<foo>",
			expectedObj:          common.PathElement{Name: "foo", Attributes: nil},
			expectedNoAttrString: "<foo>",
		},
		"b": {
			source:               "<foo%2Fbar>",
			expectedObj:          common.PathElement{Name: "foo/bar", Attributes: nil},
			expectedNoAttrString: "<foo%2Fbar>",
		},
		"c": {
			source:               "<bar%2Fbaz attr%2F2=1%2F2 attr1=A1>",
			expectedObj:          common.PathElement{Name: "bar/baz", Attributes: map[string]string{"attr/2": "1/2", "attr1": "A1"}},
			expectedNoAttrString: "<bar%2Fbaz>",
		},
		"d": {
			source:               "<f%3Eoo>",
			expectedObj:          common.PathElement{Name: "f>oo", Attributes: nil},
			expectedNoAttrString: "<f%3Eoo>",
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			actualObj, err := common.NewPathElementFromString(tCase.source)
			require.NoError(t, err)
			require.Equal(t, tCase.expectedObj, actualObj)

			actualString := actualObj.String()
			require.Equal(t, tCase.source, actualString)

			actualNoAttrString := actualObj.NoAttrString()
			require.Equal(t, tCase.expectedNoAttrString, actualNoAttrString)
		})
	}
}

func TestPathStrings(t *testing.T) {
	type testCase struct {
		source          string
		expectedPath    common.Path
		expectedNoAttrs string
	}

	testCases := map[string]testCase{
		"a": {
			source: "<foo><bar>",
			expectedPath: common.Path{
				common.PathElement{Name: "foo"},
				common.PathElement{Name: "bar"},
			},
			expectedNoAttrs: "<foo><bar>",
		},
		"b": {
			source: "<foo><b%2Fr><baz>",
			expectedPath: common.Path{
				common.PathElement{Name: "foo"},
				common.PathElement{Name: "b/r"},
				common.PathElement{Name: "baz"},
			},
			expectedNoAttrs: "<foo><b%2Fr><baz>",
		},
		"c": {
			source: "<foo><b%2Fr a1=attr-one percent=%25><baz>",
			expectedPath: common.Path{
				common.PathElement{Name: "foo"},
				common.PathElement{Name: "b/r", Attributes: map[string]string{
					"a1":      "attr-one",
					"percent": "%",
				}},
				common.PathElement{Name: "baz"},
			},
			expectedNoAttrs: "<foo><b%2Fr><baz>",
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			path, err := common.NewPathFromString(tCase.source)
			require.NoError(t, err)
			require.Equal(t, tCase.expectedPath, path)
			require.Equal(t, tCase.source, path.String())
			require.Equal(t, tCase.expectedNoAttrs, path.NoAttrString())
		})
	}
}
