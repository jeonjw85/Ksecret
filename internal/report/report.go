package report

import (
	"fmt"
	"io"

	"github.com/jeonjw85/Ksecret/internal/detect"
)

func Write(w io.Writer, format string, findings []detect.Finding) error {
	switch format {
	case "", "console":
		return writeConsole(w, findings)
	case "json":
		return writeJSON(w, findings)
	case "sarif":
		return writeSARIF(w, findings)
	default:
		return fmt.Errorf("알 수 없는 형식 %q", format)
	}
}
