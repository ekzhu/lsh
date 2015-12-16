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
	flag.StringVar(&datafile, "d", "./data/tiny_images_small.bin",
		"tiny image data file")
	flag.StringVar(&output, "o", "_knn_image.json", "output file for query results")
	flag.IntVar(&nWorker, "w", 200, "Number of threads for query tests")
	flag.IntVar(&nQuery, "q", 10, "Number of queries")
}

func main() {
	flag.Parse()
	parser := lsh.NewTinyImagePointParser()
	lsh.RunKnn(datafile, output, k, nQuery, nWorker, parser)
}
