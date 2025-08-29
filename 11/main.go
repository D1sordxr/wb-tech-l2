package main

import (
	"fmt"
	"slices"
	"strings"
)

func main() {
	sl := []string{
		"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол",
	}
	fmt.Println(searchForAnagrams(sl))
}

func searchForAnagrams(sl []string) map[string][]string {
	if len(sl) == 0 {
		return nil
	}

	anagramStore := make(map[string][]string)
	getCleanStore := func() map[string][]string {
		newStore := make(map[string][]string, len(anagramStore))
		for _, v := range anagramStore {
			if len(v) > 1 {
				newStore[v[0]] = v
			}
		}

		return newStore
	}

	sortWord := func(word string) string {
		lowerWord := strings.ToLower(word)
		runes := []rune(lowerWord)
		slices.Sort(runes)
		return string(runes)
	}

	for i := 0; i < len(sl); i++ {
		sortedWord := sortWord(sl[i])
		v, ok := anagramStore[sortedWord]
		if !ok {
			anagramStore[sortedWord] = []string{sl[i]}
		} else {
			anagramStore[sortedWord] = append(v, sl[i])
		}
	}

	return getCleanStore()
}
