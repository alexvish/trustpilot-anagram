package anagram

import (
	"testing"
)

func TestAnagramStringCreation(t *testing.T) {
	t.Run("empty", func (t *testing.T) {
		if len(NewAnagramRunes("")) != 0 {
			t.Error("length of ampty AnagramRunes is not 0\n")
		}
	})

	t.Run("abc", func(t *testing.T) {
		if ToAnagram("bac") != "abc" {
			t.Errorf("abc test failed %s should be abc\n", ToAnagram("bac"));
		}
	})

	t.Run("Same string", func(t *testing.T) {
		if ToAnagram("abc") != "abc" {
			t.Errorf("ToAnagram(\"abc\") should be abc but is %s\n", ToAnagram("abc"));
		}
	})
	t.Run("Same string siwth punktuation", func(t *testing.T) {
		if ToAnagram(" ',c,'b a'', ") != "abc" {
			t.Errorf("ToAnagram(\" ',c,'b a'', \") should be abc but is %s\n", ToAnagram(" ',c,'b a'', "));
		}
	})

	t.Run("poultry outwits ants", func(t *testing.T) {
		if ToAnagram("poultry outwits ants") != "ailnooprssttttuuwy" {
			t.Errorf("ToAnagram(\"poultry outwits ants\") should be ailnooprssttttuuwy but is %s\n", ToAnagram("poultry outwits ants"));
		}
	})


}

func TestSubtract(t *testing.T) {
	t.Run("Subtract empty from empty", func(t *testing.T){
		minuend := NewAnagramRunes("")
		subtrahend := NewAnagramRunes("")

		difference,err := minuend.Subtract(subtrahend)

		if err != nil {
			t.Errorf("'%s' - '%s' :Unexpected error %v ", minuend, subtrahend, err)
			return
		}

		if difference.String() != "" {
			t.Errorf("'%s' - '%s' = '%s', shoud be ''", minuend, subtrahend, difference)
			return
		}

	})

	t.Run("Subtract something from empty", func(t *testing.T){
		minuend := NewAnagramRunes("")
		subtrahend := NewAnagramRunes("a")

		_, err := minuend.Subtract(subtrahend)

		if err == nil {
			t.Errorf("'%s' - '%s' :Expected error %v ", minuend, subtrahend, err)
			return
		}

	})


	t.Run("Subtract trustpilot from phrase", func(t *testing.T){
		minuend := NewAnagramRunes("poultry outwits ants")
		subtrahend := NewAnagramRunes("trustpilot")

		difference,err := minuend.Subtract(subtrahend)

		if err != nil {
			t.Errorf("'%s' - '%s' :Unexpected error %v ", minuend, subtrahend, err)
			return
		}

		if difference.String() != "anostuwy" {
			t.Errorf("'%s' - '%s' = '%s', shoud be 'anostuwy'", minuend, subtrahend, difference)
			return
		}

	})

	t.Run("Subtract ad from abc", func(t *testing.T){
		minuend := NewAnagramRunes("abc")
		subtrahend := NewAnagramRunes("ad")

		difference,err := minuend.Subtract(subtrahend)

		if err != SubtractionError {
			t.Errorf("'%s' - '%s' :Expected error %v , but result is '%s'",
				minuend, subtrahend, SubtractionError, difference)
		}
	})

	t.Run("Subtract ab from aabcd", func(t *testing.T){
		minuend := NewAnagramRunes("aabcd")
		subtrahend := NewAnagramRunes("ab")

		difference,err := minuend.Subtract(subtrahend)

		if err != nil {
			t.Errorf("'%s' - '%s' :Unexpected error %v ", minuend, subtrahend, err)
			return
		}

		if difference.String() != "acd" {
			t.Errorf("'%s' - '%s' = '%s', shoud be 'anostuwy'", minuend, subtrahend, difference)
			return
		}

	})


}