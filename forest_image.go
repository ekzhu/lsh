package main

import (
	"flag"
	"lsh"
)

const (
	dim = 3072
)

var (
	datafile string
	nWorker  int
	nQuery   int
	output   string
	k        int
	m        int
	l        int
	w        float64
)

func init() {
	flag.IntVar(&k, "k", 20, "K")
	flag.StringVar(&datafile, "d", "./data/tiny_images_small.bin",
		"tiny image data file")
	flag.StringVar(&output, "o", "_forest_image.json",
		"output file for query results")
	flag.IntVar(&nWorker, "t", 200, "Number of threads for query tests")
	flag.IntVar(&nQuery, "q", 10, "Number of queries")
	flag.IntVar(&m, "m", 4, "Size of combined hash function")
	flag.IntVar(&l, "l", 25, "Number of hash tables")
	flag.Float64Var(&w, "w", 1000.0, "projection slot size")
}

func main() {
	flag.Parse()
	parser := lsh.NewTinyImagePointParser()
	lsh.RunForest(datafile, output, k, nQuery, nWorker, parser,
		dim, m, l, w)
}
