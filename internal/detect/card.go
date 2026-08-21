package detect

import "regexp"

var cardRe = regexp.MustCompile(`\d(?:[ -]?\d){12,18}`)

type cardDetector struct{}

func (cardDetector) Find(line string) []Match {
	var out []Match
	for _, loc := range cardRe.FindAllStringIndex(line, -1) {
		if !isolated(line, loc[0], loc[1]) {
			continue
		}
		secret := line[loc[0]:loc[1]]
		d := digitValues(secret)
		if len(d) < 13 || len(d) > 19 || allSame(d) || !luhn(d) {
			continue
		}
		if len(d) == 13 && checksum13(d) {
			continue
		}
		out = append(out, Match{
			RuleID:   "luhn-card",
			Severity: SeverityCritical,
			Column:   loc[0] + 1,
			Secret:   secret,
			Message:  "카드번호",
		})
	}
	return out
}

func luhn(d []int) bool {
	sum := 0
	alt := false
	for i := len(d) - 1; i >= 0; i-- {
		n := d[i]
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}
