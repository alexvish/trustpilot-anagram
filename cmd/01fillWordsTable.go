package main

import (
	"github.com/alexvish/trustpilot-anagram/pkgs/anagram"
	"io/ioutil"
	"strings"
	anagramData "github.com/alexvish/trustpilot-anagram/pkgs/data"
	"sync"
)


func LoadWords() ([]string, error) {
	content,err := ioutil.ReadFile("wordlist")
	if err != nil {
		return nil, err
	}
	return strings.Split(string(content),"\n"), nil;
}

func SendWordsToDB(ch chan <- anagramData.WordAnagramEntity) {
	defer close(ch)
	words,error := LoadWords();


	if error != nil {
		panic(error)
	}

	for _, word := range words {
		word = strings.TrimSpace(strings.ToLower(word))
		wordAnagram := anagram.NewAnagramRunes(word)
		_, err := anagram.PHRASE_ANAGRAM.Subtract(wordAnagram)
		if err != nil {
			continue;
		}
		ch <- anagramData.WordAnagramEntity{Anagram: wordAnagram.String(), Word: word}
	}
}


func main() {
	var wg sync.WaitGroup

	wg.Add(2)
	ch := make(chan anagramData.WordAnagramEntity)
	go func() {
		SendWordsToDB(ch)
		wg.Done()
	}()

	go func() {
		anagramData.WordsToDatabase(ch)
		wg.Done()
	}()
	wg.Wait()
}