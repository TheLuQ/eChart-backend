package matcher

import (
	"regexp"
	"strconv"
)

type Match struct {
	Instrument *Instrument
	Voice      []int
	Key        string
}

func (m *Match) ToString() string {
	var voice string
	if len(m.Voice) > 0 {
		voice = strconv.Itoa(m.Voice[0])
	}
	return m.Instrument.NameEng + " " + m.Instrument.NamePol + m.Key + voice
}

func (m *Match) GetVoice() string {
	if len(m.Voice) == 0 {
		return ""
	}
	return strconv.Itoa(m.Voice[0])
}

func Analyze(rawName string) (*Match, error) {
	normalizeName := normalizeName(rawName)
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	result := re.ReplaceAllString(normalizeName, " ")
	err, instrument := GetMostSimilarInstrument(result)
	if err != nil {
		return nil, err
	}

	return &Match{
		Instrument: instrument,
		Voice:      ParseVoice(normalizeName),
		Key:        ParseKey(normalizeName),
	}, nil
}
