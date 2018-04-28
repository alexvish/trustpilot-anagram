package anagram

type SortableAnagrams []AnagramRunes

func (sa SortableAnagrams) Len() int {
	return len(sa)
}

func (sa SortableAnagrams) Less(i, j int) bool {
	return sa[i].String() < sa[j].String()
}

func (sa SortableAnagrams) Swap(i, j int) {
	sa[i], sa[j] = sa[j], sa[i]
}

