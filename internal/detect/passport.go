package detect

import "regexp"

var (
	passportNewRe = regexp.MustCompile(`\b[MSRDG][0-9]{8}\b`)
	passportOldRe = regexp.MustCompile(`\b[MSRDG][A-Z][0-9]{7}\b`)
)

type passportDetector struct{}

func (passportDetector) Find(line string) []Match {
	out := findPassport(line, passportNewRe)
	return append(out, findPassport(line, passportOldRe)...)
}

func findPassport(line string, re *regexp.Regexp) []Match {
	var out []Match
	for _, loc := range re.FindAllStringIndex(line, -1) {
		out = append(out, Match{
			RuleID:   "kr-passport",
			Severity: SeverityHigh,
			Column:   loc[0] + 1,
			Secret:   line[loc[0]:loc[1]],
			Message:  "여권번호",
		})
	}
	return out
}
