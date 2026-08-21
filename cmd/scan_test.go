package cmd

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_scanCmd_returns_ErrHasFindings_when_rrn_present(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.txt"), []byte("id=900101-1123459\n"), 0o600))
	out := filepath.Join(dir, "out.json")
	rootCmd.SetOut(&bytes.Buffer{})
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"scan", dir, "--format", "json", "--output", out})
	err := rootCmd.ExecuteContext(context.Background())
	require.ErrorIs(t, err, ErrHasFindings)
}
