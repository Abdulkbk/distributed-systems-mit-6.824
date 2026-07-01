package mrapps

import (
	"mapreduce/mr"
	"strconv"
	"strings"
	"unicode"
)

// Map is the word-count map function. It is called once per input file:
//   - filename: the name of the file (unused here, but real jobs sometimes
//     want it, e.g. to build an inverted index of filename -> words).
//   - contents: the entire file as a single string.
//
// For every word in the file it emits the pair (word, "1"). Emitting "1" once
// per occurrence means the reduce step can get the total by simply counting
// how many values arrived for each word.
func Map(filename string, content string) []mr.KeyValue {
	// Split on anything that is not a letter, so "Hello, world!" -> Hello world.
	// FieldsFunc treats each run of separator runes as one delimiter and never
	// returns empty strings, which is exactly what we want.
	isSeperator := func(r rune) bool { return !unicode.IsLetter(r) }
	words := strings.FieldsFunc(content, isSeperator)

	kv := make([]mr.KeyValue, 0, len(words))
	for _, w := range words {
		// Lower-case so "The" and "the" count as the same word.
		kv = append(kv, mr.KeyValue{Key: strings.ToLower(w), Value: "1"})
	}

	return kv
}

// Reduce is the word-count reduce function. It is called once per distinct
// word, with values holding one "1" for every time that word appeared. The
// count is therefore just the length of the slice, returned as a string.
func Reduce(key string, values []string) string {
	return strconv.Itoa(len(values))
}
