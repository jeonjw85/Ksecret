package cmd

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var ErrHasFindings = errors.New("시크릿 발견")

var rootCmd = &cobra.Command{
	Use:               "ksecret",
	Short:             "한국 특화 시크릿/PII 스캐너",
	SilenceUsage:      true,
	SilenceErrors:     true,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	PersistentPreRunE: func(c *cobra.Command, args []string) error {
		return setupLogger(c)
	},
}

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}

func setupLogger(c *cobra.Command) error {
	verbose, err := c.Flags().GetBool("verbose")
	if err != nil {
		return err
	}
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
	return nil
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "상세 로그")
	rootCmd.PersistentFlags().StringP("config", "c", "", "설정 파일 경로")
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(installCmd)
}
