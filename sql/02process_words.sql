
-- count anagrams
-- SELECT count(DISTINCT anagram) FROM anagram.words WHERE excluded=FALSE;

-- exclude one letter words
-- SELECT * from anagram.words WHERE length(word) = 1;
UPDATE anagram.words SET excluded=TRUE WHERE length(word) = 1;

-- exclude words that contain only consonants
-- SELECT * FROM anagram.words WHERE word !~ '[aeiouy]'
UPDATE anagram.words SET excluded=TRUE WHERE word !~ '[aeiouy]';

