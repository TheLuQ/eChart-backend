#!/bin/env python3
import json
import re
import sys
from difflib import SequenceMatcher
from dataclasses import dataclass

from drive import DriveFile


@dataclass
class RawInstrument:
    name_pol: str
    name_en : str
    def compare(self, rawinput: str):
        tocheck = [self.name_en, self.name_pol]
        words = rawinput.split(' ')
        possibles = words if len(self.name_en.split(' ')) == 1 or len(words) == 1 else [' '.join((words[i], words[i + 1])) for i in range(len(words) - 1)]
        thestrongest = max([SequenceMatcher(a=word, b=x).ratio() for x in tocheck for word in possibles])
        return thestrongest if thestrongest > 0.7 else 0.0
        

midi_instruments_pol_to_en = {
    "saksofon sopran": "soprano saxophone",
    "saksofon alt": "alto saxophone",
    "saksofon tenor": "tenor saxophone",
    "saksofon baryton": "baritone saxophone",
    "gitara basowa": "bass guitar",
    "klarnet basowy": "bass clarinet",
    "rożek angielski": "english horn",
    "puzon basowy": "bass trombone",
    "kontrabas": "double bass",
    "pianino": "piano",
    "celesta": "celesta",
    "wibrafon": "vibraphone",
    "marimba": "marimba",
    "ksylofon": "xylophone",
    "akordeon": "accordion",
    "harmonijka": "harmonica",
    "gitara": "guitar",
    "bas": "bass",
    "skrzypce": "violin",
    "altówka": "viola",
    "wiolonczela": "cello",
    "harfa": "harp",
    "kotły": "timpani",
    "trąbka": "trumpet",
    "puzon": "trombone",
    "tuba": "tuba",
    "obój": "oboe",
    "waltornia": "horn",
    "fagot": "bassoon",
    "klarnet": "clarinet",
    "pikolo": "piccolo",
    "flet": "flute",
    "tenor": "tenor",
    "perkusja": "drums",
    "partytura": "score",
    "dzwonki": "glockenspiel",
    "fortepian": "piano",
    "róg": "horn",
    "kornet": "cornet",
    "eufonium": "euphonium",
    "baryton": "euphonium",
}

raw_instruments = [
    RawInstrument(name_pol=pol, name_en=en)
    for pol, en in midi_instruments_pol_to_en.items()
]


possible_keys = ["c", "eb", "es", "bb", "b", "f", "a","g"]

@dataclass
class Instrument:
    file_name: str
    instrument_name_en: str
    instrument_name_pol: str
    voice: str | None
    key: str | None
    file_id: str
    
    def toJson(self):
        return json.dumps({k: v for k, v in self.__dict__.items() if v is not None}, ensure_ascii=False)

def createInstrument(input: str, file_id: str | None = None):
    cleanInput = replaceJunk(input)
    instrument = matchInstrument(cleanInput.lower())
    voice = matchVoice(cleanInput)
    key = matchKey(cleanInput)
    return Instrument(input, instrument.name_en, instrument.name_pol, voice, key, file_id) if instrument and file_id else None

def replaceJunk(input: str):
    cleaned = re.sub(r'[^\w\d]|_', ' ', input)
    withoutMultipleSpaces = re.sub(r'\s+', ' ', cleaned)
    return withoutMultipleSpaces

def matchInstrument(input: str):
    candidate = max(raw_instruments, key=lambda x: x.compare(input))
    return candidate if candidate.compare(input) > 0 else None

def matchVoice(input: str):
    candidates = [candidate for candidate in re.findall(r'\d+', input) if int(candidate) in [1,2,3,4]]
    return candidates[-1] if candidates else None

def matchKey(input: str):
    candidates = [candidate for candidate in input.split(' ') if candidate.lower() in possible_keys]
    return candidates[-1].title() if candidates else None

if __name__ == "__main__":
    raw = ' '.join(sys.argv[1:])
    result = createInstrument(raw)
    print(result)
