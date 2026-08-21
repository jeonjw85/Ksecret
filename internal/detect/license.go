package detect

import (
	"regexp"
	"strings"
)

var (
	licenseHyphenRe  = regexp.MustCompile(`(?:1[1-9]|2[0-6]|28)-\d{2}-\d{6}-\d{2}`)
	licenseCompactRe = regexp.MustCompile(`(?:1[1-9]|2[0-6]|28)\d{10}`)
	licenseHintRe    = regexp.MustCompile(`(?i)면허|driver.?licen`)
)

type licenseDetector struct{}

func (licenseDetector) Find(line string) []Match {
	out := findLicense(line, licenseHyphenRe, true)
	if licenseHintRe.MatchString(line) {
		out = append(out, findLicense(line, licenseCompactRe, false)...)
	}
	return out
}

func findLicense(line string, re *regexp.Regexp, hyphenated bool) []Match {
	var out []Match
	for _, loc := range re.FindAllStringIndex(line, -1) {
		if !isolated(line, loc[0], loc[1]) {
			continue
		}
		secret := line[loc[0]:loc[1]]
		d := digitValues(secret)
		if len(d) != 12 {
			continue
		}
		region := d[0]*10 + d[1]
		if !((region >= 11 && region <= 26) || region == 28) {
			continue
		}
		if hyphenated && strings.Count(secret, "-") != 3 {
			continue
		}
		out = append(out, Match{
			RuleID:   "kr-driver-license",
			Severity: SeverityHigh,
			Column:   loc[0] + 1,
			Secret:   secret,
			Message:  "운전면허번호",
		})
	}
	return out
}
