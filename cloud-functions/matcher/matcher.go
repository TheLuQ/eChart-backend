package matcher

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"

	"golang.org/x/text/unicode/norm"
)

const MAX_INSTRUMENT_WORDS = 2

type Instrument struct {
	NamePol string
	NameEng string
}

func NewInstrument(rawName string) (error, *Instrument) {
	for i := range instruments {
		switch rawName {
		case instruments[i].NamePol:
			return nil, &instruments[i]
		case instruments[i].NameEng:
			return nil, &instruments[i]
		}
	}

	return fmt.Errorf("instrument not found: %s", rawName), nil
}

func (i *Instrument) GetSimilarityRatio(other string) float64 {
	if i.NamePol == other || i.NameEng == other {
		return 1.0
	}

	comobinations := getAllPairs(other)
	currentMax := -1.0
	for _, v := range comobinations {
		strongerResult := max(similarityRatioStrUtil(i.NameEng, v), similarityRatioStrUtil(i.NamePol, v))
		if strongerResult > currentMax {
			currentMax = strongerResult
		}
	}
	return currentMax
}

func GetMostSimilarInstrument(rawName string) (error, *Instrument) {
	var bestInstrument *Instrument
	bestRatio := -1.0

	for i := range instruments {
		ratio := instruments[i].GetSimilarityRatio(rawName)
		if ratio > bestRatio {
			bestRatio = ratio
			bestInstrument = &instruments[i]
		}
	}

	if bestInstrument == nil {
		return fmt.Errorf("no similar instrument found for: %s", rawName), nil
	}

	return nil, bestInstrument
}

func getAllPairs(rawName string) []string {
	parts := strings.Fields(rawName)
	if len(parts) < 2 {
		return []string{rawName}
	}
	shiftedParts := parts[1:]
	partsWithoutLast := parts[:len(parts)-1]
	pairs := make([]string, 0)
	for i, e := range shiftedParts {
		pairs = append(pairs, partsWithoutLast[i]+" "+e)
	}
	return pairs
}

var instruments = []Instrument{
	{"saksofon sopran", "soprano saxophone"},
	{"saksofon alt", "alto saxophone"},
	{"saksofon tenor", "tenor saxophone"},
	{"saksofon baryton", "baritone saxophone"},
	{"tenorhorn", "tenorhorn"},
	{"gitara basowa", "bass guitar"},
	{"klarnet basowy", "bass clarinet"},
	{"rożek angielski", "english horn"},
	{"puzon basowy", "bass trombone"},
	{"kontrabas", "double bass"},
	{"pianino", "piano"},
	{"celesta", "celesta"},
	{"wibrafon", "vibraphone"},
	{"marimba", "marimba"},
	{"ksylofon", "xylophone"},
	{"akordeon", "accordion"},
	{"harmonijka", "harmonica"},
	{"gitara", "guitar"},
	{"bas", "bass"},
	{"skrzypce", "violin"},
	{"altówka", "viola"},
	{"wiolonczela", "cello"},
	{"harfa", "harp"},
	{"kotły", "timpani"},
	{"trąbka", "trumpet"},
	{"puzon", "trombone"},
	{"tuba", "tuba"},
	{"obój", "oboe"},
	{"waltornia", "horn"},
	{"fagot", "bassoon"},
	{"klarnet", "clarinet"},
	{"pikolo", "piccolo"},
	{"flet", "flute"},
	{"perkusja", "drums"},
	{"partytura", "score"},
	{"dzwonki", "glockenspiel"},
	{"fortepian", "piano"},
	{"róg", "horn"},
	{"kornet", "cornet"},
	{"eufonium", "euphonium"},
}

func normalizeName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	decomposed := norm.NFD.String(s)

	out := make([]rune, 0, len(decomposed))
	for _, r := range decomposed {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		out = append(out, r)
	}
	return string(out)
}

func similarityRatioStrUtil(a, b string) float64 {
	return strutil.Similarity(a, b, metrics.NewJaroWinkler())
}
