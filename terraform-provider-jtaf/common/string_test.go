package common_test

import (
	"testing"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/stretchr/testify/require"
)

func TestJoinNonEmptyPartsWithUnderscores(t *testing.T) {
	type testCase struct {
		d []string
		e string
	}

	testCases := map[string]testCase{
		"empty": {
			d: []string{""},
			e: "",
		},
		"single": {
			d: []string{"foo"},
			e: "foo",
		},
		"leading_empty": {
			d: []string{"", "foo"},
			e: "foo",
		},
		"trailing_empty": {
			d: []string{"foo", ""},
			e: "foo",
		},
		"single_empty_ends": {
			d: []string{"", "foo", ""},
			e: "foo",
		},
		"multiple_empty_ends": {
			d: []string{"", "foo", "bar", ""},
			e: "foo_bar",
		},
		"multiple": {
			d: []string{"foo", "bar"},
			e: "foo_bar",
		},
		"middle_empty": {
			d: []string{"foo", "", "bar"},
			e: "foo_bar",
		},
		"middle_multiple_empty": {
			d: []string{"foo", "", "", "bar"},
			e: "foo_bar",
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			a := common.JoinNonEmptyPartsWithUnderscores(tCase.d...)
			require.Equal(t, tCase.e, a)
		})
	}
}
