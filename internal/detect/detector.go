package detect

type Detector interface {
	Find(line string) []Match
}

func ScanLine(line string) []Match {
	out := make([]Match, 0, 2)
	for _, d := range detectors {
		out = append(out, d.Find(line)...)
	}
	return out
}

var detectors = []Detector{
	rrnDetector{},
	frnDetector{},
	licenseDetector{},
	passportDetector{},
	brnDetector{},
	cardDetector{},
	kakaoDetector{},
	naverDetector{},
	tossDetector{},
	coolsmsDetector{},
	portoneDetector{},
	ncpDetector{},
	awsDetector{},
	gcpDetector{},
	slackDetector{},
	jwtDetector{},
}
