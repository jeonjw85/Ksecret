package detect

import "regexp"

var (
	kakaoAssignRe   = regexp.MustCompile(`(?i)(?:kakao[_-]?ak|kakao(?:[_-]?(?:rest[_-]?)?(?:api|admin|native))?[_-]?key)\s*[:=]\s*([0-9a-f]{32})`)
	kakaoHeaderRe   = regexp.MustCompile(`KakaoAK\s+([A-Za-z0-9]{32})`)
	naverSecretRe   = regexp.MustCompile(`(?i)(?:NAVER_CLIENT_SECRET|X-Naver-Client-Secret)\s*[:=]\s*(\S{8,})`)
	tossKeyRe       = regexp.MustCompile(`\b(?:live|test)_(?:sk|ck|gsk)_[A-Za-z0-9]{16,}\b`)
	coolsmsRe       = regexp.MustCompile(`(?i)(?:COOLSMS|SOLAPI)[_-]?(?:API[_-]?)?(?:KEY|SECRET)\s*[:=]\s*(\S{8,})`)
	portoneAssignRe = regexp.MustCompile(`(?i)(?:IMP_(?:API)?(?:KEY|SECRET)|PORTONE_(?:API[_-]?)?(?:KEY|SECRET))\s*[:=]\s*(\S{8,})`)
	portoneImpRe    = regexp.MustCompile(`\bimp_[A-Za-z0-9]{8,}\b`)
	ncpIamRe        = regexp.MustCompile(`\bncp_iam_[A-Za-z0-9]{8,}\b`)
	ncpAssignRe     = regexp.MustCompile(`(?i)(?:X-NCP-IAM-ACCESS-KEY|NCP_(?:SECRET|ACCESS)_KEY)\s*[:=]\s*(\S{8,})`)
)

type kakaoDetector struct{}
type naverDetector struct{}
type tossDetector struct{}
type coolsmsDetector struct{}
type portoneDetector struct{}
type ncpDetector struct{}

func (kakaoDetector) Find(line string) []Match {
	return append(keyed(line, keySpec{kakaoAssignRe, "kr-kakao", "카카오 API 키", 1}),
		keyed(line, keySpec{kakaoHeaderRe, "kr-kakao", "카카오 API 키", 1})...)
}

func (naverDetector) Find(line string) []Match {
	return keyed(line, keySpec{naverSecretRe, "kr-naver", "네이버 클라이언트 시크릿", 1})
}

func (tossDetector) Find(line string) []Match {
	return keyed(line, keySpec{tossKeyRe, "kr-toss", "토스페이먼츠 키", 0})
}

func (coolsmsDetector) Find(line string) []Match {
	return keyed(line, keySpec{coolsmsRe, "kr-coolsms", "쿨에스엠에스/솔라피 키", 1})
}

func (portoneDetector) Find(line string) []Match {
	out := keyed(line, keySpec{portoneAssignRe, "kr-portone", "포트원 키", 1})
	return append(out, keyed(line, keySpec{portoneImpRe, "kr-portone", "포트원 가맹점 아이디", 0})...)
}

func (ncpDetector) Find(line string) []Match {
	out := keyed(line, keySpec{ncpIamRe, "kr-ncp", "네이버 클라우드 IAM 키", 0})
	return append(out, keyed(line, keySpec{ncpAssignRe, "kr-ncp", "네이버 클라우드 IAM 키", 1})...)
}

type keySpec struct {
	re    *regexp.Regexp
	rule  string
	msg   string
	group int
}

func keyed(line string, spec keySpec) []Match {
	var out []Match
	for _, sub := range spec.re.FindAllStringSubmatchIndex(line, -1) {
		if len(sub) < (spec.group+1)*2 {
			continue
		}
		lo, hi := sub[spec.group*2], sub[spec.group*2+1]
		if lo < 0 {
			continue
		}
		out = append(out, Match{
			RuleID:   spec.rule,
			Severity: SeverityCritical,
			Column:   lo + 1,
			Secret:   line[lo:hi],
			Message:  spec.msg,
		})
	}
	return out
}
