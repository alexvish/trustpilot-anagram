package anagram

import (
	"sync"
)

func findAnagrams(
	test AnagramRunes,
	seed []AnagramRunes,
	stack []AnagramRunes,
	found chan <- []AnagramRunes,
	anagramByFirstChar map[string][]AnagramRunes,
	seedFirstChar string,
	wg *sync.WaitGroup) {

	defer wg.Done()

	for i, next := range seed {
		diff, err := test.Subtract(next)
		if err != nil {
			continue;
		}

		lenStack := len(stack)
		newStack := make([]AnagramRunes, lenStack, lenStack + 1)
		copy(newStack, stack)
		newStack = append(newStack,next)

		if diff.Len() == 0 {
			found <- newStack
			continue;
		}
		//need to go deeper
		diffFirstChar := diff.FirstChar();
		var newSeed []AnagramRunes;
		if diffFirstChar == seedFirstChar {
			// optimization: do not repeat check on the whole seed array if
			// seed is the same as previous seed
			newSeed = seed[i:]
		} else {
			newSeed = anagramByFirstChar[diffFirstChar]
		}

		wg.Add(1);
		go findAnagrams(diff, newSeed, newStack, found, anagramByFirstChar, diffFirstChar, wg);
	}
}

func FindAnagrams(seed []AnagramRunes, words []AnagramRunes, found chan <- []AnagramRunes) {
	defer close(found)
	var wg sync.WaitGroup

	anagramToFirstChar := make(map[string][]AnagramRunes);
	for _,w := range words {
		ch := w.FirstChar()

		anagramToFirstChar[ch] = append(anagramToFirstChar[ch], w)
	}

	stack := make([]AnagramRunes, 0)
	wg.Add(1)
	findAnagrams(PHRASE_ANAGRAM,seed, stack, found, anagramToFirstChar, "init",  &wg)
	wg.Wait();
}
