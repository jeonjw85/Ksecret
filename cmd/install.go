package cmd

import (
	"fmt"
	"os"

	"github.com/jeonjw85/Ksecret/internal/hook"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "로컬 git 훅을 설치합니다",
	RunE:  runInstall,
}

func init() {
	installCmd.Flags().Bool("pre-commit", false, ".git/hooks/pre-commit 설치")
}

func runInstall(c *cobra.Command, args []string) error {
	pre, err := c.Flags().GetBool("pre-commit")
	if err != nil {
		return err
	}
	if !pre {
		return fmt.Errorf("--pre-commit 이 필요합니다")
	}
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("현재 디렉토리: %w", err)
	}
	root, err := hook.GitRoot(cwd)
	if err != nil {
		return err
	}
	if err := hook.Install(root); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "pre-commit 훅을 설치했습니다")
	return nil
}
