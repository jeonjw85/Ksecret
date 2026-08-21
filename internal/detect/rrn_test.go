package detect

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_rrn_hits_when_checksum_and_date_valid(t *testing.T) {
	got := ScanLine("id=900101-1123459")
	require.True(t, hasRule(got, "kr-rrn"))
}

func Test_rrn_misses_when_checksum_wrong(t *testing.T) {
	got := ScanLine("id=900101-1123458")
	require.False(t, hasRule(got, "kr-rrn"))
}

func Test_frn_hits_when_gender_5_to_8(t *testing.T) {
	got := ScanLine("id=900101-5123450")
	require.True(t, hasRule(got, "kr-frn"))
	require.False(t, hasRule(got, "kr-rrn"))
}

func hasRule(ms []Match, id string) bool {
	for _, m := range ms {
		if m.RuleID == id {
			return true
		}
	}
	return false
}
