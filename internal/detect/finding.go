package detect

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
)

type Finding struct {
	RuleID   string   `json:"rule_id"`
	Severity Severity `json:"severity"`
	File     string   `json:"file"`
	Line     int      `json:"line"`
	Column   int      `json:"column"`
	Secret   string   `json:"secret"`
	Message  string   `json:"message"`
}

type Match struct {
	RuleID   string
	Severity Severity
	Column   int
	Secret   string
	Message  string
}

func Mask(secret string) string {
	if len(secret) <= 8 {
		return secret[:1] + "..."
	}
	return secret[:4] + "..." + secret[len(secret)-4:]
}
