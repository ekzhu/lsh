package main

import (
	"flag"
	"fmt"
	"log"
	"lsh"
	"path/filepath"
)

const (
	dim = 3072
)

var (
	datafile   string
	knnout     string
	varloutdir string
	vartoutdir string
	nWorker    int
	nQuery     int
	k          int
	m          int
	l          int
	w          float64
	t          int
	ls         []int
	ts         []int
)

func init() {
	flag.IntVar(&k, "k", 20, "Number of nearest neighbours")
	flag.StringVar(&datafile, "d", "./data/tiny_images_small.bin",
		"tiny image data file")
	flag.StringVar(&varloutdir, "varlout", "",
		"Output directory for experiment with different Ls")
	flag.StringVar(&vartoutdir, "vartout", "",
		"Output directory for experiment with different Ts")
	flag.IntVar(&nWorker, "t", 200, "Number of threads for query tests")
	flag.IntVar(&nQuery, "q", 10, "Number of queries")
	flag.IntVar(&m, "M", 4, "Size of combined hash function")
	flag.IntVar(&l, "L", 24, "Number of hash tables")
	flag.IntVar(&t, "T", 8, "Length of probing sequence in Multi-probe")
	flag.Float64Var(&w, "W", 3000.0, "projection slot size")
	knnout = "_knn_image"
	ls = []int{1, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40}
	ts = []int{4, 8, 16}
}

type AnalysisResult struct {
	Algorithm   string   `json:"algorithm"`
	ResultFiles []string `json:"result_files"`
}

type VarLMeta struct {
	AnalysisResults []AnalysisResult `json:"analysis_results"`
	M               int
	W               float64
	K               int `json:"k"`
	T               int
	Ls              []int
}

type VarTMeta struct {
	AnalysisResults []AnalysisResult `json:"analysis_results"`
	M               int
	W               float64
	K               int `json:"k"`
	L               int
	Ts              []int
}

func resultFileName(outdir, algorithm, paramName string, paramVal int) string {
	filename := fmt.Sprintf("%s_%s_%d", algorithm, paramName, paramVal)
	return filepath.Join(outdir, filename)
}

func analysisFileName(outdir, algorithm, paramName string, paramVal int) string {
	f := resultFileName(outdir, algorithm, paramName, paramVal)
	return fmt.Sprintf("%s_%s", f, "analysis")
}

func main() {
	flag.Parse()
	if vartoutdir == "" || varloutdir == "" {
		log.Fatal("No output directory given")
		return
	}
	parser := lsh.NewTinyImagePointParser()
	data := lsh.LoadData(datafile, parser)
	queries := lsh.SelectQueriesAsSubset(data, nQuery)

	// Run exact kNN
	log.Println("Running exact kNN")
	lsh.RunKnn(data, queries, knnout, k, nWorker)

	var analysisResults []string

	// Run Var L experiments
	varlmeta := VarLMeta{
		AnalysisResults: make([]AnalysisResult, 0),
		M:               m,
		W:               w,
		K:               k,
		T:               t,
		Ls:              ls,
	}
	// Basic LSH
	analysisResults = make([]string, 0)
	for _, l := range ls {
		log.Printf("Running Basic LSH: l = %d\n", l)
		result := resultFileName(varloutdir, "basic", "l", l)
		lsh.RunSimple(data, queries, result, k, nWorker, dim, m, l, w)
		analysis := analysisFileName(varloutdir, "basic", "l", l)
		lsh.RunAnalysis(result, knnout, analysis)
		analysisResults = append(analysisResults, analysis)
	}
	varlmeta.AnalysisResults = append(varlmeta.AnalysisResults,
		AnalysisResult{"Basic", analysisResults})
	// LSH Forest
	analysisResults = make([]string, 0)
	for _, l := range ls {
		log.Printf("Running LSH Forest: l = %d\n", l)
		result := resultFileName(varloutdir, "forest", "l", l)
		lsh.RunForest(data, queries, result, k, nWorker, dim, m, l, w)
		analysis := analysisFileName(varloutdir, "forest", "l", l)
		lsh.RunAnalysis(result, knnout, analysis)
		analysisResults = append(analysisResults, analysis)
	}
	varlmeta.AnalysisResults = append(varlmeta.AnalysisResults,
		AnalysisResult{"Forest", analysisResults})
	// Multi-probe
	analysisResults = make([]string, 0)
	for _, l := range ls {
		log.Printf("Running Multi-probe LSH: l = %d\n", l)
		result := resultFileName(varloutdir, "multiprobe", "l", l)
		lsh.RunMultiprobe(data, queries, result, k, nWorker, dim, m, l, w, t)
		analysis := analysisFileName(varloutdir, "multiprobe", "l", l)
		lsh.RunAnalysis(result, knnout, analysis)
		analysisResults = append(analysisResults, analysis)
	}
	varlmeta.AnalysisResults = append(varlmeta.AnalysisResults,
		AnalysisResult{"Multi-probe", analysisResults})
	lsh.DumpJson(filepath.Join(varloutdir, ".meta"), &varlmeta)

	// Run Var T experiments
	vartmeta := VarTMeta{
		AnalysisResults: make([]AnalysisResult, 0),
		M:               m,
		W:               w,
		K:               k,
		L:               l,
		Ts:              ts,
	}
	// Multi-probe
	analysisResults = make([]string, 0)
	for _, t := range ts {
		log.Printf("Running Multi-probe LSH: t = %d\n", t)
		result := resultFileName(vartoutdir, "multiprobe", "t", t)
		lsh.RunMultiprobe(data, queries, result, k, nWorker, dim, m, l, w, t)
		analysis := analysisFileName(vartoutdir, "multiprobe", "t", t)
		lsh.RunAnalysis(result, knnout, analysis)
		analysisResults = append(analysisResults, analysis)
	}
	vartmeta.AnalysisResults = append(vartmeta.AnalysisResults,
		AnalysisResult{"Multi-probe", analysisResults})
	lsh.DumpJson(filepath.Join(vartoutdir, ".meta"), &vartmeta)

}
