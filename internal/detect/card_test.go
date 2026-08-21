package detect

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_card_hits_when_luhn_valid(t *testing.T) {
	got := ScanLine("pan=4111111111111111")
	require.True(t, hasRule(got, "luhn-card"))
}

func Test_card_misses_when_luhn_wrong(t *testing.T) {
	got := ScanLine("pan=4111111111111112")
	require.False(t, hasRule(got, "luhn-card"))
}
