package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/jeonjw85/Ksecret/internal/config"
	"github.com/jeonjw85/Ksecret/internal/report"
	"github.com/jeonjw85/Ksecret/internal/scan"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "파일 트리에서 시크릿/PII를 찾습니다",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runScan,
}

func init() {
	scanCmd.Flags().String("format", "console", "출력 형식 (console, json, sarif)")
	scanCmd.Flags().String("output", "", "결과 파일 경로 (비우면 표준출력)")
	scanCmd.Flags().Bool("staged", false, "커밋 대기 파일만 스캔")
}

func runScan(c *cobra.Command, args []string) error {
	root := "."
	if len(args) == 1 {
		root = args[0]
	}
	format, err := c.Flags().GetString("format")
	if err != nil {
		return err
	}
	output, err := c.Flags().GetString("output")
	if err != nil {
		return err
	}
	staged, err := c.Flags().GetBool("staged")
	if err != nil {
		return err
	}
	configPath, err := c.Flags().GetString("config")
	if err != nil {
		return err
	}
	cfg, err := config.LoadForRoot(root, configPath)
	if err != nil {
		return err
	}
	findings, err := scan.Run(c.Context(), scan.Options{Root: root, Staged: staged, Config: cfg})
	if err != nil {
		return err
	}
	out := io.Writer(os.Stdout)
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("결과 파일 생성: %w", err)
		}
		defer f.Close()
		out = f
	}
	if err := report.Write(out, format, findings); err != nil {
		return err
	}
	if len(findings) > 0 {
		return ErrHasFindings
	}
	return nil
}
