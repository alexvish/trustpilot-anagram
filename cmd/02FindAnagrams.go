package main

import (
	"github.com/alexvish/trustpilot-anagram/pkgs/anagram"
	anagramData "github.com/alexvish/trustpilot-anagram/pkgs/data"
	"sync"
)


func FindAnagrams(found chan <- []anagram.AnagramRunes) {
	wordAnagrams := anagramData.LoadWordAnagrams()
	seed := anagramData.LoadSeedAnagrams()

	anagram.FindAnagrams(seed, wordAnagrams, found)
}

func main() {
	found := make(chan []anagram.AnagramRunes, 100)
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		FindAnagrams(found)
	} ()

	go func() {
		defer wg.Done()
		anagramData.StoreAnagramResults(found)
	} ()

	wg.Wait()

}
