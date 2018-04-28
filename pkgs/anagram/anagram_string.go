package anagram

import (
	"sort"
	"unicode"
	"unicode/utf8"
	"errors"
)

const PHRASE = "poultry outwits ants"
var PHRASE_ANAGRAM = NewAnagramRunes(PHRASE)

var MinuendLessThenSubtrahend = errors.New("Minuend length is less than subtrahend length");
var SubtractionError = errors.New("Subtraction Error")


type AnagramRunes []rune


func (r AnagramRunes) Len() int {
	return utf8.RuneCountInString(string(r))
}

func (r AnagramRunes) Less(i, j int) bool {
	return r[i] < r[j]
}

func (r AnagramRunes) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r AnagramRunes) Subtract(s AnagramRunes) (AnagramRunes, error) {
	if r.Len() < s.Len() {
		return nil, MinuendLessThenSubtrahend
	}

	diffSlice := make([]rune, 0, r.Len() - s.Len())

	for {
		if s.Len() == 0 {
			diffSlice = append(diffSlice, r...)
			return AnagramRunes(diffSlice), nil
		}

		var i int;
		for i = 0; i < r.Len(); i++ {
			if r[i] >= s[0] {
				break;
			}
		}
		diffSlice = append(diffSlice,r[:i]...)
		r = r[i:]
		if r.Len() == 0 {
			// No more runes in minuend, still some runes in subtrahend
			return nil, SubtractionError
		}
		if r[0] > s[0] {
			return nil, SubtractionError
		} else {
			for r.Len() > 0 && s.Len() > 0 && r[0] == s[0] {
				r, s = r[1:], s[1:]
			}
		}
	}
}

func (r AnagramRunes) FirstChar() string {
	if r.Len() == 0 {
		return "<empty>"
	}
	return string(r[0:1])
}

func (r AnagramRunes) String() string {
	return string(r)
}



func RemoveSpaceAndPunkt(sourceRunes []rune) []rune {
	var filteredRunes []rune

	for i,r := range sourceRunes {
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			filteredRunes = make([]rune,i,len(sourceRunes))
			copy(filteredRunes,sourceRunes[:i])
			sourceRunes = sourceRunes[i + 1:]
			break;
		}
	}
	if filteredRunes == nil {
		return sourceRunes;
	}
	for _, r := range sourceRunes {
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			continue;
		}
		filteredRunes = append(filteredRunes, r);
	}
	return filteredRunes
}


func NewAnagramRunes(source string) AnagramRunes {
	result := AnagramRunes(RemoveSpaceAndPunkt([]rune(source)))
	sort.Sort(result)
	return result;
}

func ToAnagram(s string) string {
	return NewAnagramRunes(s).String()
}
