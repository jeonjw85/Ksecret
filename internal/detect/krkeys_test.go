package detect

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_kakao_hits_when_key_assigned(t *testing.T) {
	got := ScanLine("KAKAO_REST_API_KEY=0123456789abcdef0123456789abcdef")
	require.True(t, hasRule(got, "kr-kakao"))
}

func Test_toss_hits_when_secret_prefix(t *testing.T) {
	got := ScanLine("key=test_sk_a1b2c3d4e5f6g7h8i9j0")
	require.True(t, hasRule(got, "kr-toss"))
}

func Test_kakao_misses_when_hex_without_keyword(t *testing.T) {
	got := ScanLine("token=0123456789abcdef0123456789abcdef")
	require.False(t, hasRule(got, "kr-kakao"))
}
