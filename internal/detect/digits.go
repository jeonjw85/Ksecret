package detect

func isDigit(b byte) bool { return b >= '0' && b <= '9' }

func digitValues(s string) []int {
	out := make([]int, 0, len(s))
	for i := 0; i < len(s); i++ {
		if isDigit(s[i]) {
			out = append(out, int(s[i]-'0'))
		}
	}
	return out
}

func isolated(line string, start, end int) bool {
	if start > 0 && isDigit(line[start-1]) {
		return false
	}
	if end < len(line) && isDigit(line[end]) {
		return false
	}
	return true
}

func allSame(d []int) bool {
	if len(d) == 0 {
		return true
	}
	for i := 1; i < len(d); i++ {
		if d[i] != d[0] {
			return false
		}
	}
	return true
}
