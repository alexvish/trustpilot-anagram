package data

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/alexvish/trustpilot-anagram/pkgs/anagram"
	"sort"
	"strings"
	"log"
	"sync"
)

const DATABASE_CONNECT_STRING = "postgres://postgres:postgres@localhost/postgres?search_path=anagram&sslmode=disable"

type WordAnagramEntity struct{
	Word string
	Anagram string
}

func WordsToDatabase(ch <-chan(WordAnagramEntity)) {
	db, err := sql.Open("postgres", DATABASE_CONNECT_STRING)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	stmt, err:= tx.Prepare("INSERT INTO anagram.words (word, anagram) VALUES ($1, $2) ON CONFLICT DO NOTHING")
	if err != nil {
		panic(err)
	}
	for e := range ch {
		stmt.Exec(e.Word, e.Anagram)
	}
	stmt.Close()
	tx.Commit()
}

func LoadWordAnagrams() []anagram.AnagramRunes {
	db, err := sql.Open("postgres", DATABASE_CONNECT_STRING)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT DISTINCT anagram FROM anagram.words WHERE excluded = FALSE")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	result := make([]anagram.AnagramRunes,0,1500)
	for rows.Next() {
		var anagramStr string;
		rows.Scan(&anagramStr)
		result = append(result, anagram.NewAnagramRunes(anagramStr))
	}

	return result;
}

func LoadSeedAnagrams() []anagram.AnagramRunes {
	db, err := sql.Open("postgres", DATABASE_CONNECT_STRING)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT DISTINCT anagram FROM anagram.words WHERE excluded = FALSE AND word in (SELECT word from seed)")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	result := make([]anagram.AnagramRunes,0,1500)
	for rows.Next() {
		var anagramStr string;
		rows.Scan(&anagramStr)
		result = append(result, anagram.NewAnagramRunes(anagramStr))
	}
	return result;
}

func StoreAnagramResults(foundAnagrams <- chan []anagram.AnagramRunes) {
	db, err := sql.Open("postgres", DATABASE_CONNECT_STRING)
	if err != nil {
		panic(err)
	}
	defer db.Close()


	for found := range foundAnagrams {
		sort.Sort(anagram.SortableAnagrams(found))

		keyComponents := make([]string,0,len(found));
		anagramMap := make(map[string]int)
		for _, component := range found {
			compStr := component.String();
			keyComponents = append(keyComponents, compStr)
			anagramMap[compStr]++
		}
		key := strings.Join(keyComponents," ");

		tx, err := db.Begin()
		if err != nil {
			panic(err)
		}
		anagramInsertStmt, err := tx.Prepare("INSERT INTO anagram.anagrams (key) VALUES ($1)")
		if err != nil {
			panic(err)
		}
		anagramWordsInsertStmt, err := tx.Prepare("INSERT INTO anagram.anagram_words (anagram_key, word_anagram, count) VALUES ($1, $2, $3)")
		if  err != nil {
			panic(err)
		}
		_, err = anagramInsertStmt.Exec(key)
		if err != nil {
			log.Println(err)
			anagramInsertStmt.Close()
			anagramWordsInsertStmt.Close()
			tx.Rollback()
			continue
		}
		isRollback := false
		for anagramWord, count := range anagramMap {
			_, err := anagramWordsInsertStmt.Exec(key, anagramWord, count)
			if err != nil {
				log.Println(err)
				anagramWordsInsertStmt.Close()
				anagramInsertStmt.Close()
				tx.Rollback()
				isRollback = true
				break
			}
		}

		if isRollback {
			continue
		}

		anagramInsertStmt.Close()
		anagramWordsInsertStmt.Close()
		tx.Commit()
	}
}

func LoadWordsByAnagram() map[string][]string {
	db, err := sql.Open("postgres", DATABASE_CONNECT_STRING)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT anagram, word FROM anagram.words WHERE excluded = FALSE")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	result := make(map[string][]string)
	for rows.Next() {
		var anagramStr, wordStr string;
		rows.Scan(&anagramStr, &wordStr)
		result[anagramStr] = append(result[anagramStr], wordStr)
	}
	return result;
}

func LoadAnagramsToCheck() []string {
	db, err := sql.Open("postgres", DATABASE_CONNECT_STRING)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT key FROM anagram.anagrams WHERE md5checkdone = FALSE")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	result := make([]string,0,1024)
	for rows.Next() {
		var anagramStr string;
		rows.Scan(&anagramStr)
		result = append(result, anagramStr)
	}
	return result;
}



func StoreSearchResult(anagramKey, phrase, md5sum string) {
	db, err := sql.Open("postgres", DATABASE_CONNECT_STRING)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(
		"INSERT INTO anagram.search_result (anagram_key, phrase, md5sum) VALUES ($1, $2, $3)",
			anagramKey, phrase, md5sum)

	if err != nil {
		log.Println(err)
	}
}

type AnagramProcessedMsg struct {
	AnagramKey string
	Count int32
}

func MarkAnagramProcessed(in <- chan  AnagramProcessedMsg, wg *sync.WaitGroup) {
	db, err := sql.Open("postgres", DATABASE_CONNECT_STRING)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	updateStmt, err := tx.Prepare("UPDATE anagram.anagrams SET md5CheckDone = TRUE, md5CheckAnagramCount = $1 WHERE key = $2")
	if err != nil {
		panic(err)
	}
	for msg := range in {
		_, err := updateStmt.Exec(msg.Count, msg.AnagramKey)
		if err != nil {
			log.Println(err)
		}
	}
	updateStmt.Close()
	tx.Commit()
	wg.Done()
}