package report

import (
	"encoding/json"
	"io"

	"github.com/jeonjw85/Ksecret/internal/detect"
)

func writeSARIF(w io.Writer, findings []detect.Finding) error {
	rules := map[string]struct{}{}
	var results []sarifResult
	for _, f := range findings {
		rules[f.RuleID] = struct{}{}
		results = append(results, sarifResult{
			RuleID: f.RuleID,
			Level:  sarifLevel(f.Severity),
			Message: sarifText{
				Text: f.Message + " " + detect.Mask(f.Secret),
			},
			Locations: []sarifLocation{{
				Physical: sarifPhysical{
					Artifact: sarifArtifact{URI: f.File},
					Region:   sarifRegion{StartLine: f.Line, StartColumn: f.Column},
				},
			}},
		})
	}
	var ruleList []sarifRule
	for id := range rules {
		ruleList = append(ruleList, sarifRule{ID: id, Name: id})
	}
	doc := sarifDoc{
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Version: "2.1.0",
		Runs: []sarifRun{{
			Tool: sarifTool{Driver: sarifDriver{
				Name:           "ksecret",
				InformationURI: "https://github.com/jeonjw85/Ksecret",
				Rules:          ruleList,
			}},
			Results: results,
		}},
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(doc)
}

func sarifLevel(s detect.Severity) string {
	switch s {
	case detect.SeverityCritical, detect.SeverityHigh:
		return "error"
	default:
		return "warning"
	}
}

type sarifDoc struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string      `json:"name"`
	InformationURI string      `json:"informationUri"`
	Rules          []sarifRule `json:"rules"`
}

type sarifRule struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type sarifResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   sarifText       `json:"message"`
	Locations []sarifLocation `json:"locations"`
}

type sarifText struct {
	Text string `json:"text"`
}

type sarifLocation struct {
	Physical sarifPhysical `json:"physicalLocation"`
}

type sarifPhysical struct {
	Artifact sarifArtifact `json:"artifactLocation"`
	Region   sarifRegion   `json:"region"`
}

type sarifArtifact struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
}
