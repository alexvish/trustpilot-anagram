

CREATE SCHEMA IF NOT EXISTS anagram AUTHORIZATION postgres;

SET search_path TO anagram;

CREATE TABLE IF NOT EXISTS words (
  word VARCHAR(100),
  anagram VARCHAR(100),
  excluded BOOLEAN DEFAULT FALSE,
  PRIMARY KEY(word)
);

CREATE INDEX IF NOT EXISTS words_anagram ON words (anagram);



CREATE TABLE IF NOT EXISTS seed (
  word VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS anagrams (
  key VARCHAR(100),
  md5CheckDone BOOLEAN DEFAULT FALSE,
  md5CheckAnagramCount INTEGER,
  PRIMARY KEY (key)
);


CREATE TABLE IF NOT EXISTS anagram_words (
  anagram_key VARCHAR(100),
  word_anagram VARCHAR(100),
  count INT,
  PRIMARY KEY (anagram_key, word_anagram),
  FOREIGN KEY (anagram_key) REFERENCES anagrams(key)
);

CREATE INDEX IF NOT EXISTS anagram_words_word_anagram_index ON anagram_words(word_anagram);

CREATE TABLE IF NOT EXISTS search_result (
  anagram_key VARCHAR(100),
  phrase VARCHAR(100),
  md5sum VARCHAR(100),
  PRIMARY KEY (md5sum)
);
