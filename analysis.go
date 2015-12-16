package main

import (
	"flag"
	"lsh"
)

var (
	resultFile      string
	groundTruthFile string
	output          string
)

func init() {
	flag.StringVar(&resultFile, "r", "", "Query result file")
	flag.StringVar(&groundTruthFile, "g", "", "Ground truth result file")
	flag.StringVar(&output, "o", "_analysis.json", "Output analysis file")
}

func main() {
	flag.Parse()
	if resultFile == "" {
		panic("Missing query result file")
	}
	if groundTruthFile == "" {
		panic("Missing ground truth result file")
	}
	lsh.RunAnalysis(resultFile, groundTruthFile, output)
}
