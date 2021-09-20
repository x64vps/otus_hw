package hw03frequencyanalysis

import (
	"regexp"
	"strings"
)

type TopWord struct {
	count int
	word  string
}

func (c *TopWord) Initialize() {
	c.count = 0
	c.word = ""
}

func (c *TopWord) Check(word string, count int) {
	if count > c.count {
		c.count = count
		c.word = word
	}

	if count == c.count && strings.Compare(c.word, word) == 1 {
		c.word = word
	}
}

func (c *TopWord) String() string {
	return c.word
}

func Top10(s string) []string {
	var result []string

	words := getWordsWithCount(s)
	if len(words) == 0 {
		return result
	}

	var curr TopWord

	for len(result) < 10 && len(words) > 0 {
		curr.Initialize()

		for word, count := range words {
			curr.Check(word, count)
		}

		result = append(result, curr.String())
		delete(words, curr.String())
	}

	return result
}

func getWordsWithCount(s string) map[string]int {
	result := make(map[string]int)

	r := regexp.MustCompile(`[^A-zА-яеЁ]*([A-zА-я]+[A-zА-я-]*)[^A-zА-яеЁ]*`)

	for _, w := range strings.Fields(s) {
		if !r.MatchString(w) {
			continue
		}

		w = r.ReplaceAllString(w, "$1")

		result[strings.ToLower(w)]++
	}

	return result
}
