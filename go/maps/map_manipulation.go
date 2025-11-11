package maps

import (
	"log/slog"
	"strings"
)

func main() {
}

// this paticualr problem was not easy to get right off the cuff, maps have always been a bit tricky for me
// as I really don't havfe to work with them a whole log, I couldn't get it done in 35mins
// I didn't realy need to track previousk/first, etc, I over comlicted that and there is a better solution
// than this even, see addToMap2
func addToMap() {
	s := "cat bat cat mat bat cat"
	wordsSlice := strings.Split(s, " ")

	// "cat": 1,
	wordCounts := make(map[string]int)

	for i, v := range wordsSlice {
		slog.Info("--->", "v", v, "i", i)
		// if word isn't in map add it
		if wordCounts[v] == 0 {
			wordCounts[v] = 1
		} else { // update it
			wordCounts[v] = wordCounts[v] + 1
		}

	}

	slog.Info("")
	slog.Info("Map", "wordCounts", wordCounts)
}

func addToMap2() {
	s := "cat bat cat mat bat cat"
	wordsSlice := strings.Split(s, " ")

	// "cat": 1,
	wordCounts := make(map[string]int)

	for _, v := range wordsSlice {
		// maps in go have 0 value for non-existent
		// so no need to check if its there or not, this just adds or updates it
		wordCounts[v]++
	}
	slog.Info("Map", "wordCounts", wordCounts)
}
