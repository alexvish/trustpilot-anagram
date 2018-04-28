-- at first insert a single word "trustpilot"

INSERT INTO anagram.seed (word) VALUES ('trustpilot');
-- execute:
--      go run cmd/02FindAnagrams.go
--      go run cmd/03CheckMd5Sums.go
-- No anagrams, that match md5sums found



-- Now Clean up
DELETE FROM anagram.seed;

-- Fill in seed table with "long" words
INSERT INTO anagram.seed SELECT word FROM anagram.words where length(anagram) > 7;

-- exclude word trustpilot: we already checked all anagrams that contain this word
UPDATE anagram.words SET excluded = TRUE WHERE word = 'trustpilot';

-- execute:
--      go run cmd/02FindAnagrams.go
--      go run cmd/03CheckMd5Sums.go





















-- Check statement

SELECT
  word, anagram
FROM
  anagram.words
WHERE
    excluded = FALSE
  AND
    word IN (
      SELECT word from anagram.seed
    );

