package detect

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_jwt_hits_when_payload_json_valid(t *testing.T) {
	got := ScanLine("t=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjMiLCJpYXQiOjE1MTYyMzkwMjJ9.dGVzdA")
	require.True(t, hasRule(got, "jwt"))
}

func Test_aws_hits_when_akia_prefix(t *testing.T) {
	got := ScanLine("AKIAIOSFODNN7EXAMPLE")
	require.True(t, hasRule(got, "aws-access-key"))
}
