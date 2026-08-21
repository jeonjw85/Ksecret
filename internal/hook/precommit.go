package hook

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const marker = "# ksecret pre-commit"

func Install(repoRoot string) error {
	gitDir := filepath.Join(repoRoot, ".git")
	info, err := os.Stat(gitDir)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("git 저장소가 아닙니다! :  %s", repoRoot)
	}
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("실행 파일 경로 : %w", err)
	}
	hookPath := filepath.Join(gitDir, "hooks", "pre-commit")
	if existing, err := os.ReadFile(hookPath); err == nil && !strings.Contains(string(existing), marker) {
		return fmt.Errorf("pre-commit 훅이 이미 있습니다 : %s", hookPath)
	}
	if err := os.MkdirAll(filepath.Dir(hookPath), 0o755); err != nil {
		return fmt.Errorf("hooks 디렉터리 : %w", err)
	}
	body := fmt.Sprintf("#!/bin/sh\n%s\nexec %q scan --staged\n", marker, exe)
	if err := os.WriteFile(hookPath, []byte(body), 0o755); err != nil {
		return fmt.Errorf("훅 파일 쓰기 : %w", err)
	}
	return nil
}

func GitRoot(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git 루트를 찾지 못했습니다 : %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
