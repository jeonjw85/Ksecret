package detect

import "regexp"

var brnRe = regexp.MustCompile(`\d{3}-?\d{2}-?\d{5}`)

type brnDetector struct{}

func (brnDetector) Find(line string) []Match {
	var out []Match
	for _, loc := range brnRe.FindAllStringIndex(line, -1) {
		if !isolated(line, loc[0], loc[1]) {
			continue
		}
		secret := line[loc[0]:loc[1]]
		d := digitValues(secret)
		if len(d) != 10 || !checksumBRN(d) {
			continue
		}
		out = append(out, Match{
			RuleID:   "kr-brn",
			Severity: SeverityHigh,
			Column:   loc[0] + 1,
			Secret:   secret,
			Message:  "사업자등록번호",
		})
	}
	return out
}

func checksumBRN(d []int) bool {
	w := [...]int{1, 3, 7, 1, 3, 7, 1, 3, 5}
	sum := 0
	for i := 0; i < 9; i++ {
		sum += d[i] * w[i]
	}
	sum += (d[8] * 5) / 10
	return d[9] == (10-sum%10)%10
}
