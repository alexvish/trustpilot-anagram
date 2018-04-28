# Trustpilot follow the white rabbit

A company named [Trustpilot](https://www.trustpilot.com/) hires developers. Candidates should submit solution
to thier code challenge. This is the solution to [anagram code challenge](http://followthewhiterabbit.trustpilot.com/cs/step3.html).


Solution to trustpilot pony challenge can be found [here](https://alexvish.github.io/ponychallenge-trustpilot/)

## Task

You need to find anagrams to phrase "poultry outwits ants" which md5 sums are specified in task description.
Here are the md5sums:

- e4820b45d2277f3844eac66c903e84be
- 23170acc097c24edb98fc5488ab033fe
- 665e5bcb0c20062fe8abaaf4628bb154


## Know what you are looking for.

As a first step I suggest to search google for md5 sums.
You will be directed to md5 reversal services, and you will find the answers.
I will try to keep them secret until the end of this document, but
generally they do not look like a normal English phrases.
They are not "beautiful" whatsoever, you cannot discern them from
incorrect answers in a "stream" of anagrams.

## Solution

### Anagram

Define anagram of the word of phrase as follows:
1. Strip spaces and special characters
2. Sort letters alphabetically

Anagram of phrase "poultry outwits ants" will be "ailnooprssttttuuwy"
Anagram of "trustpilot" will be: "iloprstttu"

Multiple words can produce the same anagram: 
```
[lony lyon only] => lnoy
```

### Subtract operation

Define `subtract` operation on two anagrams that produce another anagram:
```
A1 - A2 = R 
```
as:
1. For each letter in A2 remove the letter from A1.
2. If there is no corresponding letter in A1 then operation results in error.

Example:
```
In words: 
1. "poultry outwits ants" - "trustpilot"
2. "poultry outwits ants" - "only"

In anagrams:
1. "ailnooprssttttuuwy" - "iloprstttu" = "anostuwy"
2. "ailnooprssttttuuwy" - "lnoy" = "aioprssttttuuw"

As we see, both "trustpilot" and "only" can be part of anagram "poultry outwits ants".

But these words cannot be the part of anagram at the same time, because

"poultry outwits ants" - "trustpilot" - "only" ~ "ailnooprssttttuuwy" - "iloprstttu" - "lnoy"

results in error: there is not 'l' letter in the result of first subtraction.

```


Thus the subtraction operation answers the following two questions:

1. Can a word or phrase be a part of another anagram? - There is an error in subtraction if not.
2. If a word or phrase is a part of anagram, what is anagram that the rest of words must be a part - This is shown by succcessful subtration result.

### Anagram search

Now another aspect of anagram search: if anagram that is being tested
starts with some letter, it must contain at least one of those words
which anagrams starts with the same letter.

```
"poultry outwits ants" -> "ailnooprssttttuuwy"
 must contain the word which anagram starts with letter "a", such as
"wants" -> "anstw"

And if anagram contains "wants" the reminder:
"poultry outwits ants" - "wants" = "ilooprstttuuy"
must contain word, which anagram starts with "i", such as
"trustpilot" -> "iloprstttu"
```

So, we will find anagrams trying to subtract anagram, which starts with
the same letter as the test anagram, and use the reminder as a new
test anagram on the next recursive step. We will do this except
for very first step, in which we will use a pre-selected seed anagrams.

### Divide and Conquer

Word list contains less then 100K words, of which only about 2500
words can be part of anagram of original phrase. there are 1377 distinct
anagrams of such words.

But there are many millions of possible anagrams, that can be created
from these words, so we need a way to save intermediate results in order
 to systematically search for anagrams.

So, scripts in this repository work with local postgres database,
defined in a vagrant box. Definition of a vagrant box is in file
[vagrant/Vagrantfile](https://github.com/alexvish/trustpilot-anagram/tree/master/vagrant/Vagrantfile)

Database schema is defined in [sql/01anagram_schema.sql](https://github.com/alexvish/trustpilot-anagram/tree/master/sql/01anagram_schema.sql)
Schema contains tables:
 -  **words** - contains words, that can be part of anagram, their anagrams
    and also contains flag to exclude words
 - **anagrams** - contains calculated anagrams, which are not anagrams
   yet, but anagrams can be generated from data in this table. Basically,
   the key of this table is a list of anagrams of words, in alphabetical
   order separated by space. To generate anagrams from these entries,
   we need to look up words by anagrams in a **words** table and generate
   all possible combinations. Boolean field `md5CheckDone` indicates
   that this anagram was already processed, field md5CheckAnagramCount
   shows how many md5 checks was done for the anagram.
 - **seed** contain seed words (as words) that will be used as an initial
   seed for creating data in **anagrams** table from **words** table.
   For each such calculation, anagrams that are added to **anagrams**
   table will contain at least one anagram, that matches a word
   from the **seed** table.
 - **search_result** table contains final phrases, md5sums, and reference
   to anagram from anagrams table


Scripts:
 - [cmd/01fillWordsTable.go](https://github.com/alexvish/trustpilot-anagram/tree/master/cmd/01fillWordsTable.go) - populate table
   **words** from wordlist file
 - [cmd/02FindAnagrams.go](https://github.com/alexvish/trustpilot-anagram/tree/master/cmd/02FindAnagrams.go) - calculate anagrams
   and populate **anagrams** table. Table **words** is used as word
   anagram source, table **seed** is used as inital seed. Algorithm
   is described in "Anagram Search" section.
 - [cmd/03CheckMd5Sums.go](https://github.com/alexvish/trustpilot-anagram/tree/master/cmd/03CheckMd5Sums.go) - takes anagrams
   from **anagrams** table, then using **words** table calculate all
   possible combinations of words for particular anagram, calculates
   md5sums and if matches found store results to **search_result** table.
   Also, marks anagrams as processed and reports the number of times
   an md5sum was calculated for particular anagram.

## Lets do the search

### Start database
So, first let us start a vagrant box with postgres.
Change dir to [vagrant](https://github.com/alexvish/trustpilot-anagram/tree/master/vagrant) and execute
```
vagrant up
```
A vagrant box will be provisioned with postgres database on it.
That database can be accessed as:
```
Host: localhost
Port: 5432
Database: postgres
User: postgres
Password: postgres
sslMode: off
```


### Create database schema
Now connect to database and populate anagram schema. Just execute file
[sql/01anagram_schema.sql](https://github.com/alexvish/trustpilot-anagram/tree/master/sql/01anagram_schema.sql)
using your favorite database tool, against the database.

### Fill words table
Execute
```
 go run cmd/01fillWordsTable.go
```
Now table **words** contain 2497 words, that have 1377 distinct
anagrams (one anagram can match multiple words)

Now let us make a little grooming on words. We will exclude words
that contain a single letter and words that consist only from consonants

Execute these statements from file [sql/02process_words.sql](https://github.com/alexvish/trustpilot-anagram/tree/master/sql/02process_words.sql):
```
UPDATE anagram.words SET excluded=TRUE WHERE length(word) = 1;
UPDATE anagram.words SET excluded=TRUE WHERE word !~ '[aeiouy]';
```
Do not worry, we can unexclude any word later.
Now we have 2450 and 1338 anagrams to process.

### Start with a word "trustpilot" as seed

For the start let us find all anagrams that contains word "trustpilot"
and check whether any of them match predefined checksum.

Use this statement to populate **seed** table:
```
INSERT INTO anagram.seed (word) VALUES ('trustpilot');
```

Let us find all anagrams:
```
go run cmd/02FindAnagrams.go
```

Now table **anagrams** contain 92 anagrams. Now let us check md5 sums:
```
go run cmd/03CheckMd5Sums.go
```

No hashsum is found: table *search_result* remains empty.
So far we checked 92 anagrams, that give us 11262 strings, from which we
calculated md5 sums.

### Choose different seed
Ok, let us first clean up a bit. Execute the following sql statements
against database:
```
-- clean up seed table
DELETE FROM anagram.seed;

-- Now we already processed all anagrams that contain word "trustpilot"
-- Let us exclude this word:
UPDATE anagram.words SET excluded = TRUE WHERE word = 'trustpilot';

```

So, now we have to choose a different seed. Phrase "poultry outwit ants"
consists of 18 non-space letters, so there is a chance that at least
one of the words in anagram that we are trying to find contains at least
8 letters. Let us use all words, that are longer then 7 letters as a seed:

```
-- Fill in seed table with "long" words
INSERT INTO anagram.seed SELECT word FROM anagram.words where length(anagram) > 7;

```
Now seed contains 171 words.

Execute:
```
$ time go run cmd/02FindAnagrams.go

real    1m13,171s
user    0m0,000s
sys     0m0,015s

$ time go run cmd/03CheckMd5Sums.go
Anagram found:
ty outlaws printouts - 23170acc097c24edb98fc5488ab033fe
Anagram found:
wu lisp not statutory - 665e5bcb0c20062fe8abaaf4628bb154
Anagram found:
printout stout yawls - e4820b45d2277f3844eac66c903e84be

real    0m12,306s
user    0m0,000s
sys     0m0,015s

```
And we found all three!

our **search_result** table now contains:

|       anagram_key     |         phrase           |               md5sum             |
|-----------------------|--------------------------|----------------------------------|
| alostuw inoprsttu ty  | ty outlaws printouts     | 23170acc097c24edb98fc5488ab033fe |
| aorstttuy ilps not uw | wu lisp not statutory    | 665e5bcb0c20062fe8abaaf4628bb154 |
| alswy inoprttu osttu  | printout stout yawls     | e4820b45d2277f3844eac66c903e84be |

In total we have found 10646 distinct "anagram sources", and calculated
hash sum for 4095516 anagrams.