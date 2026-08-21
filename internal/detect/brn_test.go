package detect

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_brn_hits_when_checksum_valid(t *testing.T) {
	got := ScanLine("biz=120-81-47521")
	require.True(t, hasRule(got, "kr-brn"))
}

func Test_brn_misses_when_checksum_wrong(t *testing.T) {
	got := ScanLine("biz=120-81-47522")
	require.False(t, hasRule(got, "kr-brn"))
}
