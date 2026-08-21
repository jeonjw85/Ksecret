package report

import (
	"encoding/json"
	"io"

	"github.com/jeonjw85/Ksecret/internal/detect"
)

type jsonResult struct {
	Findings []detect.Finding `json:"findings"`
	Counts   map[string]int   `json:"counts"`
}

func writeJSON(w io.Writer, findings []detect.Finding) error {
	masked := make([]detect.Finding, len(findings))
	counts := map[string]int{}
	for i, f := range findings {
		f.Secret = detect.Mask(f.Secret)
		masked[i] = f
		counts[f.RuleID]++
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(jsonResult{Findings: masked, Counts: counts})
}
