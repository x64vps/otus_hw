package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(s string) []string {
	words := getWordsWithCount(s)
	if len(words) == 0 {
		return []string{}
	}

	sorted := getSorted(words)

	if len(sorted) > 10 {
		return sorted[:10]
	}

	return sorted
}

const pattern = `[^A-zА-яеЁ]*([A-zА-я]+[A-zА-я-]*)[^A-zА-яеЁ]*`

func getWordsWithCount(s string) map[string]int {
	result := make(map[string]int)

	r := regexp.MustCompile(pattern)

	for _, w := range strings.Fields(s) {
		if !r.MatchString(w) {
			continue
		}

		w = r.ReplaceAllString(w, "$1")

		result[strings.ToLower(w)]++
	}

	return result
}

func getSorted(words map[string]int) []string {
	uniq := make([]string, len(words))
	i := 0

	for w := range words {
		uniq[i] = w
		i++
	}

	sort.Slice(uniq, func(i, j int) bool {
		if words[uniq[i]] == words[uniq[j]] {
			return uniq[i] < uniq[j]
		}

		return words[uniq[i]] > words[uniq[j]]
	})

	return uniq
}
