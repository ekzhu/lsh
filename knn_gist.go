package main

import (
	"flag"
	"lsh"
)

var (
	datafile string
	nWorker  int
	nQuery   int
	output   string
	k        int
)

func init() {
	flag.IntVar(&k, "k", 20, "K")
	flag.StringVar(&datafile, "d", "./data/tiny_gist_small.bin",
		"tiny image gist data file")
	flag.StringVar(&output, "o", "_knn_gist.json", "output file for query results")
	flag.IntVar(&nWorker, "w", 200, "Number of threads for query tests")
	flag.IntVar(&nQuery, "q", 10, "Number of queries")
}

func main() {
	flag.Parse()
	parser := lsh.NewTinyImageGistParser()
	lsh.RunKnn(datafile, output, k, nQuery, nWorker, parser)
}
