package scan

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jeonjw85/Ksecret/internal/config"
	"github.com/stretchr/testify/require"
)

func Test_Run_hits_expected_rules_when_scanning_testdata(t *testing.T) {
	findings, err := Run(context.Background(), Options{
		Root:   testdataDir(t),
		Config: config.Default(),
	})
	require.NoError(t, err)
	got := map[string]bool{}
	for _, f := range findings {
		got[f.RuleID] = true
	}
	for _, id := range []string{"kr-rrn", "kr-brn", "luhn-card", "kr-kakao", "kr-toss"} {
		require.True(t, got[id], "missing %s in %#v", id, got)
	}
	for _, f := range findings {
		require.NotEqual(t, "clean.txt", f.File)
	}
}

func Test_Run_skips_finding_when_allowlist_pattern_matches(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.env"), []byte("AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n"), 0o600))
	cfg := config.Default()
	cfg.Allowlist.Patterns = []string{"EXAMPLE"}
	findings, err := Run(context.Background(), Options{Root: dir, Config: cfg})
	require.NoError(t, err)
	require.Empty(t, findings)
}

func Test_Run_skips_file_when_allowlist_path_matches(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "secret.env"), []byte("id=900101-1123459\n"), 0o600))
	cfg := config.Default()
	cfg.Allowlist.Paths = []string{"secret.env"}
	findings, err := Run(context.Background(), Options{Root: dir, Config: cfg})
	require.NoError(t, err)
	require.Empty(t, findings)
}

func testdataDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "testdata"))
}
