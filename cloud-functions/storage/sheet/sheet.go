package sheet

import (
	"fmt"
	"regexp"
)

type Instrument struct {
	Name   string `firestore:"instrument"`
	Number string `firestore:"number,omitempty"`
	Key    string `firestore:"key,omitempty"`
}

func ParseFileName(name string) (Instrument, error) {
	filePattern := regexp.MustCompile(`(?i)^\s*([a-z]+(?:\s+[a-z]+)*)\s*([1-9]\d*)?(?:\s+in\s+([a-z]{1,2}))?\s*$`)

	instrument := filePattern.FindStringSubmatch(name)
	if instrument == nil || len(instrument) < 4 {
		return Instrument{}, fmt.Errorf("invalid instrument filename: %q", name)
	}
	return Instrument{instrument[1], instrument[2], instrument[3]}, nil
}
