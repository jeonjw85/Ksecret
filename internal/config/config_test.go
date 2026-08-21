package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_MatchPath_matches_double_star_prefix(t *testing.T) {
	require.True(t, MatchPath("testdata/**", "testdata/hits.env"))
	require.False(t, MatchPath("testdata/**", "internal/scan/scan.go"))
}

func Test_Allowed_excludes_when_pattern_in_line(t *testing.T) {
	cfg := Config{Allowlist: Allowlist{Patterns: []string{"EXAMPLE"}}}
	require.True(t, cfg.Allowed("a.env", "aws-access-key", "AKIAIOSFODNN7EXAMPLE"))
	require.False(t, cfg.Allowed("a.env", "aws-access-key", "AKIAAREALKEYNOTHERE00"))
}

func Test_Allowed_excludes_when_rule_listed(t *testing.T) {
	cfg := Config{Allowlist: Allowlist{Rules: []string{"jwt"}}}
	require.True(t, cfg.Allowed("a.go", "jwt", "eyJ"))
	require.False(t, cfg.Allowed("a.go", "kr-rrn", "x"))
}
