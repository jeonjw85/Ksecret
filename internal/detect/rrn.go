package detect

import (
	"regexp"
	"time"
)

var ident13Re = regexp.MustCompile(`\d{6}[-\s]?\d{7}`)

type rrnDetector struct{}
type frnDetector struct{}

func (rrnDetector) Find(line string) []Match {
	return findIdent13(line, identSpec{"kr-rrn", "주민등록번호", koreanGender})
}

func (frnDetector) Find(line string) []Match {
	return findIdent13(line, identSpec{"kr-frn", "외국인등록번호", foreignerGender})
}

type identSpec struct {
	rule     string
	msg      string
	genderOK func(int) bool
}

func koreanGender(g int) bool    { return g >= 1 && g <= 4 }
func foreignerGender(g int) bool { return g >= 5 && g <= 8 }

func findIdent13(line string, spec identSpec) []Match {
	var out []Match
	for _, loc := range ident13Re.FindAllStringIndex(line, -1) {
		if !isolated(line, loc[0], loc[1]) {
			continue
		}
		secret := line[loc[0]:loc[1]]
		d := digitValues(secret)
		if len(d) != 13 || !spec.genderOK(d[6]) || !validCivilDate(d) || !checksum13(d) {
			continue
		}
		out = append(out, Match{
			RuleID:   spec.rule,
			Severity: SeverityCritical,
			Column:   loc[0] + 1,
			Secret:   secret,
			Message:  spec.msg,
		})
	}
	return out
}

func checksum13(d []int) bool {
	w := [...]int{2, 3, 4, 5, 6, 7, 8, 9, 2, 3, 4, 5}
	sum := 0
	for i := 0; i < 12; i++ {
		sum += d[i] * w[i]
	}
	return d[12] == (11-sum%11)%10
}

func validCivilDate(d []int) bool {
	yy := d[0]*10 + d[1]
	mm := d[2]*10 + d[3]
	dd := d[4]*10 + d[5]
	year := 1900 + yy
	switch d[6] {
	case 3, 4, 7, 8:
		year = 2000 + yy
	}
	t := time.Date(year, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
	return t.Year() == year && int(t.Month()) == mm && t.Day() == dd
}
