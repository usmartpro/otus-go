package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	// карта слов (слово -> кол-во)
	words := make(map[string]int)
	for _, word := range strings.Fields(text) {
		words[word]++
	}

	// уникальные слова
	uniqWords := make([]string, 0, len(words))
	for word := range words {
		uniqWords = append(uniqWords, word)
	}

	// сортировка уникальных слов лексикографически
	sort.Strings(uniqWords)

	// сортировка НЕодинаковых слов по частоте
	sort.SliceStable(uniqWords, func(i, j int) bool {
		return words[uniqWords[i]] > words[uniqWords[j]]
	})

	var res []string
	countWords := len(uniqWords)

	// выборка топ-10 слов в засисимости от кол-ва слов
	switch {
	case countWords > 10:
		res = uniqWords[0:10]
	case countWords == 0:
		res = nil
	default:
		res = uniqWords
	}

	return res
}
