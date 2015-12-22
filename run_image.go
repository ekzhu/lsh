package main

import (
	"flag"
	"fmt"
	"log"
	"lsh"
	"os"
	"path/filepath"
)

const (
	dim = 3072
)

var (
	datafile   string
	knnresult  string
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
	ms         []int
	ws         []float64
)

func init() {
	flag.IntVar(&k, "k", 50, "Number of nearest neighbours")
	flag.StringVar(&datafile, "d", "./data/tiny_images_10k.bin",
		"tiny image data file")
	flag.StringVar(&varloutdir, "varlout", "",
		"Output directory for experiment with different Ls")
	flag.StringVar(&vartoutdir, "vartout", "",
		"Output directory for experiment with different Ts")
	flag.StringVar(&knnresult, "knnresult", "_knn_image_10k_k_50",
		"Exact k-NN result file, will re-run exact k-NN if not exist")
	flag.IntVar(&nWorker, "t", 200, "Number of threads for query tests")
	flag.IntVar(&nQuery, "q", 1000, "Number of queries")
	flag.IntVar(&t, "T", 64, "Length of probing sequence in Multi-probe")
	flag.IntVar(&m, "M", 9, "Size of combined hash function")
	flag.Float64Var(&w, "W", 8000.0, "projection slot size")
	flag.IntVar(&l, "L", 4, "Number of hash tables")
	ls = []int{2, 4, 8, 16, 32, 64}
	//ms = []int{9, 9, 9, 9, 9, 9}
	//ws = []float64{8000.0, 8000.0, 8000.0, 8000.0, 8000.0, 8000.0}
	ms = []int{5, 7, 9, 11, 11, 11}
	ws = []float64{12398.0, 11683.0, 11153.0, 10778.0, 9093.0, 7889.0}
	ts = []int{2, 4, 8, 16, 32, 64, 128}
}

type AnalysisResult struct {
	Algorithm   string   `json:"algorithm"`
	ResultFiles []string `json:"result_files"`
}

type VarLMeta struct {
	AnalysisResults []AnalysisResult `json:"analysis_results"`
	Ms              []int
	Ws              []float64
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

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
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

	exist, err := exists(knnresult)
	if err != nil {
		panic(err.Error())
	}
	if !exist {
		// Run exact kNN
		log.Println("Running exact kNN")
		lsh.RunKnn(data, queries, knnresult, k, nWorker)
	}

	var analysisResults []string

	// Run Var L experiments
	varlmeta := VarLMeta{
		AnalysisResults: make([]AnalysisResult, 0),
		Ls:              ls,
		Ms:              ms,
		Ws:              ws,
		K:               k,
		T:               t,
	}
	// Basic LSH
	analysisResults = make([]string, 0)
	for i, l := range ls {
		log.Printf("Running Basic LSH: l = %d\n", l)
		result := resultFileName(varloutdir, "basic", "l", l)
		lsh.RunSimple(data, queries, result, k, nWorker, dim, ms[i], l, ws[i])
		analysis := analysisFileName(varloutdir, "basic", "l", l)
		lsh.RunAnalysis(result, knnresult, analysis)
		analysisResults = append(analysisResults, analysis)
	}
	varlmeta.AnalysisResults = append(varlmeta.AnalysisResults,
		AnalysisResult{"Basic", analysisResults})
	// LSH Forest
	analysisResults = make([]string, 0)
	for i, l := range ls {
		log.Printf("Running LSH Forest: l = %d\n", l)
		result := resultFileName(varloutdir, "forest", "l", l)
		lsh.RunForest(data, queries, result, k, nWorker, dim, ms[i], l, ws[i])
		analysis := analysisFileName(varloutdir, "forest", "l", l)
		lsh.RunAnalysis(result, knnresult, analysis)
		analysisResults = append(analysisResults, analysis)
	}
	varlmeta.AnalysisResults = append(varlmeta.AnalysisResults,
		AnalysisResult{"Forest", analysisResults})
	// Multi-probe
	analysisResults = make([]string, 0)
	for i, l := range ls {
		log.Printf("Running Multi-probe LSH: l = %d\n", l)
		result := resultFileName(varloutdir, "multiprobe", "l", l)
		lsh.RunMultiprobe(data, queries, result, k, nWorker, dim, ms[i], l, ws[i], t)
		analysis := analysisFileName(varloutdir, "multiprobe", "l", l)
		lsh.RunAnalysis(result, knnresult, analysis)
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
		lsh.RunAnalysis(result, knnresult, analysis)
		analysisResults = append(analysisResults, analysis)
	}
	vartmeta.AnalysisResults = append(vartmeta.AnalysisResults,
		AnalysisResult{"Multi-probe", analysisResults})
	lsh.DumpJson(filepath.Join(vartoutdir, ".meta"), &vartmeta)

}
