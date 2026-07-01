package mr

import (
	"fmt"
	"os"
	"sort"
)

// MapFunc and ReduceFunc are the two functions every MapReduce application
// must supply. We accept them as parameters instead of hard-coding word count,
// so the same engine can run any job (word count, inverted index, etc.).
//
//	MapFunc:    given a file's name and its full contents, return the pairs
//	            that file contributes.
//	ReduceFunc: given one key and every value emitted for it, return a single
//	            result string (e.g. the total count).
type MapFunc func(filename string, content string) []KeyValue
type ReduceFunc func(key string, values []string) string

// RunSequential executes an entire MapReduce job in a single goroutine.
//
// It runs the three classic phases in order:
//
//  1. MAP     – apply mapf to every input file, collecting all emitted pairs
//     into one big in-memory slice.
//  2. SHUFFLE – sort that slice by key so identical keys become adjacent.
//  3. REDUCE  – walk the sorted slice, and for each run of identical keys,
//     call reducef once and write "key result" to the output file.
func Runsequential(inputFiles []string, mapf MapFunc,
	reducef ReduceFunc, outputFile string) error {

	// ---------- 1. MAP PHASE ----------
	// intermediate accumulates all KeyValue entries from each file.
	intermediate := []KeyValue{}
	for _, filename := range inputFiles {
		content, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("cannot read input file %q: %w", filename, err)
		}

		// mapf turns one file into many key/value pairs; append them all.
		kv := mapf(filename, string(content))
		intermediate = append(intermediate, kv...)
	}

	// ---------- 2. SHUFFLE PHASE ----------
	sort.Sort(ByKey(intermediate))

	// ---------- 3. REDUCE PHASE ----------
	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("cannot create output file %q: %w", outputFile, err)
	}
	defer out.Close()

	i := 0
	for i < len(intermediate) {
		// intermediate[i:j] is the block of pairs that all share one key.
		// Advance j until the key changes (or we hit the end).
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}

		// Gather just the values for this key into a slice for reducef.
		values := make([]string, 0, j-1)
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}

		// One reduce call per distinct key. Write "key result" on its own line.
		result := reducef(intermediate[i].Key, values)
		fmt.Fprintf(out, "%v %v\n", intermediate[i].Key, result)

		// Jump past this block to the next distinct key.
		i = j
	}

	return nil
}
