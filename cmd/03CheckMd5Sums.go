package main

import (
	"github.com/alexvish/trustpilot-anagram/pkgs/data"
	"fmt"
	"strings"
	"sync"
	"crypto/md5"
	"encoding/hex"
	"sync/atomic"
)

var wordByAnagram map[string][]string
var anagramsToCheck []string




type AnagramMsg struct {
	anagram string
	wg *sync.WaitGroup
	counter *int32
}

func _anagramToWordsArray(anagram string, wg *sync.WaitGroup, counter *int32, out chan <- WordsMsg) {
	defer wg.Done()
	anagrams := strings.Split(anagram, " ")
	length := len(anagrams)
	wordIndexes := make([]int, length)
	for {
		words := make([]string, 0, length)
		for i, a := range anagrams {
			words = append(words, wordByAnagram[a][wordIndexes[i]])
		}
		wg.Add(1)
		out <- WordsMsg{
			anagram: anagram,
			words: words,
			wg: wg,
			counter: counter,
		}

		for k := 0; ; k++ {
			if k >= length {
				return;
			}
			wordIndexes[k]++
			if wordIndexes[k] >= len(wordByAnagram[anagrams[k]]) {
				wordIndexes[k] = 0
				continue
			}
			break;
		}
	}
}
func anagramToWordsArray(in  <-chan AnagramMsg, out chan <- WordsMsg) {
	for anagramMsg := range in {
		_anagramToWordsArray(anagramMsg.anagram, anagramMsg.wg, anagramMsg.counter, out)
	}
}

type WordsMsg struct {
	anagram string
	words []string
	wg *sync.WaitGroup
	counter *int32
}

func _permutations(anagram string, words []string, wg *sync.WaitGroup, counter *int32, out chan <- PermutationMsg) {
	defer wg.Done()
	nextShifts := func (p []int) {
		for i := len(p) - 1; i >= 0; i-- {
			if i == 0 || p[i] < len(p)-i-1 {
				p[i]++
				return
			}
			p[i] = 0
		}
	}

	for p := make([]int, len(words)); p[0] < len(p); nextShifts(p) {
		//copy words array
		permWords := make([]string, 0, len(words))
		permWords = append(permWords, words...)

		//permutate
		for i, shift := range p {
			permWords[i], permWords[i + shift] = permWords[i + shift], permWords[i]
		}

		wg.Add(1)
		out <- PermutationMsg{
			anagram: anagram,
			permutation: strings.Join(permWords, " "),
			wg: wg,
			counter: counter,
		}
	}
}

func permutations(in <- chan WordsMsg, out chan <- PermutationMsg) {
	for msg := range in {
		_permutations(msg.anagram, msg.words, msg.wg, msg.counter, out)
	}
}

type PermutationMsg struct {
	anagram string
	permutation string
	wg *sync.WaitGroup
	counter *int32
}

func _digest(anagram, permutation string, wg *sync.WaitGroup, counter *int32, out chan <- FoundAnagramMsg) {
	defer wg.Done()
	hasher := md5.New()
	hasher.Write([]byte(permutation))
	hash := hex.EncodeToString(hasher.Sum(nil))

	atomic.AddInt32(counter, 1)

	if hash == "e4820b45d2277f3844eac66c903e84be" ||
		hash == "23170acc097c24edb98fc5488ab033fe" ||
		hash == "665e5bcb0c20062fe8abaaf4628bb154" {

			wg.Add(1)
			out <- FoundAnagramMsg{
				anagram: anagram,
				result: permutation,
				md5sum: hash,
				wg: wg,
			}
	}

}

func digest( in <- chan PermutationMsg, out chan <- FoundAnagramMsg) {
	for msg := range in {
		_digest(msg.anagram, msg.permutation, msg.wg, msg.counter, out)
	}
}


type FoundAnagramMsg struct {
	anagram string
	result string
	md5sum string
	wg *sync.WaitGroup
}

func storeFound (in <- chan FoundAnagramMsg) {
	for msg := range in {
		fmt.Printf("Anagram found:\n%s - %s\n", msg.result, msg.md5sum)
		data.StoreSearchResult(msg.anagram, msg.result, msg.md5sum)
		msg.wg.Done()
	}
}



func main() {
	anagramCh := make(chan AnagramMsg, 50)
	wordsCh := make(chan WordsMsg, 50)
	permCh  := make(chan PermutationMsg, 50)
	foundCh := make(chan FoundAnagramMsg, 50)

	defer close(foundCh)
	defer close(permCh)
	defer close(wordsCh)
	defer close(anagramCh)

	for i := 0; i < 4; i++ {
		go anagramToWordsArray(anagramCh, wordsCh)
		go permutations(wordsCh, permCh)
		go digest(permCh, foundCh)
		go storeFound(foundCh)
	}


	var markCompleteWg sync.WaitGroup
	markCompleteCh := make(chan data.AnagramProcessedMsg,50)

	markCompleteWg.Add(1)
	go data.MarkAnagramProcessed(markCompleteCh, &markCompleteWg)


	wordByAnagram = data.LoadWordsByAnagram()
	anagramsToCheck = data.LoadAnagramsToCheck()

	for _, anagram := range anagramsToCheck {
		var wg sync.WaitGroup
		var anagramCount int32

		wg.Add(1)
		anagramCh <- AnagramMsg{
			anagram: anagram,
			wg: &wg,
			counter: &anagramCount,
		}

		wg.Wait()

		markCompleteCh <- data.AnagramProcessedMsg{
			AnagramKey: anagram,
			Count: anagramCount,
		}
	}
	close(markCompleteCh)

	markCompleteWg.Wait()

}
