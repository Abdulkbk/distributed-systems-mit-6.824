package main

import (
	"fmt"
	"mapreduce/mr"
	"mapreduce/mrapps"
	"os"
)

// mrsequential is the command-line entry point for the sequential MapReduce.
//
// Usage:
//
//	go run ./main input1.txt input2.txt ...
//
// It wires the word-count application (mrapps.Map / mrapps.Reduce) into the
// generic engine (mr.RunSequential) and writes the result to mr-out-0 in the
// current directory.
func main() {
	// os.Args[0] is the program name; the rest are the input files.
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: mrsequential inputfile ...")
		os.Exit(1)
	}

	inputFiles := os.Args[1:]

	const outputFile = "mr-out-0"

	if err := mr.Runsequential(inputFiles, mrapps.Map, mrapps.Reduce, outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "mapreduce failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("done: results written to %s\n", outputFile)
}
