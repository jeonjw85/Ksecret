package report

import (
	"fmt"
	"io"

	"github.com/jeonjw85/Ksecret/internal/detect"
)

func writeConsole(w io.Writer, findings []detect.Finding) error {
	if len(findings) == 0 {
		_, err := fmt.Fprintln(w, "ksecret: 이슈 없음")
		return err
	}
	if _, err := fmt.Fprintf(w, "ksecret: 이슈 %d건\n\n", len(findings)); err != nil {
		return err
	}
	for _, f := range findings {
		if _, err := fmt.Fprintf(w, "%-20s %-8s %s:%d:%d\n  %s  %s\n\n",
			f.RuleID, koreanSeverity(f.Severity), f.File, f.Line, f.Column, f.Message, detect.Mask(f.Secret)); err != nil {
			return err
		}
	}
	return nil
}

func koreanSeverity(s detect.Severity) string {
	switch s {
	case detect.SeverityCritical:
		return "심각"
	case detect.SeverityHigh:
		return "높음"
	default:
		return "보통"
	}
}
