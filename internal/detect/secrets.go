package detect

import (
	"encoding/base64"
	"encoding/json"
	"regexp"
	"strings"
)

var (
	awsRe   = regexp.MustCompile(`\bAKIA[0-9A-Z]{16}\b`)
	gcpRe   = regexp.MustCompile(`\bAIza[0-9A-Za-z\-_]{35}\b`)
	slackRe = regexp.MustCompile(`\bxox[baprs]-[A-Za-z0-9-]{10,}\b`)
	jwtRe   = regexp.MustCompile(`\beyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\b`)
)

type awsDetector struct{}
type gcpDetector struct{}
type slackDetector struct{}
type jwtDetector struct{}

func (awsDetector) Find(line string) []Match {
	return keyed(line, keySpec{awsRe, "aws-access-key", "AWS 액세스 키", 0})
}

func (gcpDetector) Find(line string) []Match {
	return keyed(line, keySpec{gcpRe, "gcp-api-key", "GCP API 키", 0})
}

func (slackDetector) Find(line string) []Match {
	return keyed(line, keySpec{slackRe, "slack-token", "슬랙 토큰", 0})
}

func (jwtDetector) Find(line string) []Match {
	var out []Match
	for _, loc := range jwtRe.FindAllStringIndex(line, -1) {
		secret := line[loc[0]:loc[1]]
		if !validJWT(secret) {
			continue
		}
		out = append(out, Match{
			RuleID:   "jwt",
			Severity: SeverityMedium,
			Column:   loc[0] + 1,
			Secret:   secret,
			Message:  "JWT 토큰",
		})
	}
	return out
}

func validJWT(s string) bool {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return false
	}
	header, herr := decodeB64URL(parts[0])
	payload, perr := decodeB64URL(parts[1])
	if herr != nil || perr != nil {
		return false
	}
	var h, p map[string]any
	if json.Unmarshal(header, &h) != nil || json.Unmarshal(payload, &p) != nil {
		return false
	}
	_, alg := h["alg"]
	_, typ := h["typ"]
	_, sub := p["sub"]
	_, iat := p["iat"]
	return alg || typ || sub || iat
}

func decodeB64URL(s string) ([]byte, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err == nil {
		return b, nil
	}
	return base64.URLEncoding.DecodeString(s)
}
