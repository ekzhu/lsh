package main

import (
	"flag"
	"lsh"
)

var (
	datafile    string
	nWorker     int
	nSample     int
	k           int
	distOutput  string
	kDistOutput string
)

func init() {
	flag.IntVar(&k, "k", 1000, "K")
	flag.StringVar(&datafile, "d", "",
		"tiny image data file")
	flag.IntVar(&nWorker, "w", 200, "Number of threads")
	flag.IntVar(&nSample, "n", 1000, "sample size")
	distOutput = "_image_query_distance_sample"
	kDistOutput = "_image_all_pair_distance_sample"
}

func main() {
	flag.Parse()
	if datafile == "" {
		panic("No datafile given")
	}
	parser := lsh.NewTinyImagePointParser()

	// Query distance sample
	lsh.RunKnn(datafile, distOutput, k, nSample, nWorker, parser)

	// All pair distance sample
	lsh.RunKnnSampleAllPair(datafile, kDistOutput, nSample, nWorker, parser)
}
