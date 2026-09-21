package main

import (
	"os"
	"strconv"

	"github.com/TheLuQ/eChart-backend/matcher"
)

func main() {
	promptInput := os.Args[1]
	match, err := matcher.Analyze(promptInput)

	if err != nil {
		println("Error: " + err.Error())
	} else {
		possibleVoices := ""
		if len(match.Voice) > 0 {
			possibleVoices = " and voice: " + strconv.Itoa(match.Voice[0])
		}
		println("Most similar instrument for '" + promptInput + "' is: " + match.Instrument.NameEng + " (" + match.Instrument.NamePol + ")" + " with key: " + match.Key + possibleVoices)
	}
}
