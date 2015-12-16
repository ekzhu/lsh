package lsh

import "sort"

// DataPoint is a wrapper for Point for experiments
// as it involes the id of the points
type DataPoint struct {
	Id    int
	Point Point
}

type Neighbour struct {
	Id       int     `json:"id"`
	Distance float64 `json:"distance"`
}

type Neighbours []Neighbour

func (r Neighbours) Len() int           { return len(r) }
func (r Neighbours) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r Neighbours) Less(i, j int) bool { return r[i].Distance < r[j].Distance }

type QueryResult struct {
	QueryId    int        `json:"id"`
	Neighbours Neighbours `json:"neighbours"`
	Time       float64    `json:"time"`
}

type QueryResults []QueryResult

func (r QueryResults) Len() int           { return len(r) }
func (r QueryResults) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r QueryResults) Less(i, j int) bool { return r[i].QueryId < r[j].QueryId }

// Recall provides the classic definition of overlap-based recall
func Recall(ann, groundTruth QueryResult) float64 {
	if len(groundTruth.Neighbours) == 0 {
		return 1.0
	}
	truth := make(map[int]bool)
	for _, n := range groundTruth.Neighbours {
		truth[n.Id] = true
	}
	overlap := 0
	for _, n := range ann.Neighbours {
		if _, found := truth[n.Id]; found {
			overlap += 1
		}
	}
	return float64(overlap) / float64(len(groundTruth.Neighbours))
}

// ErrorRatio provide a quality indicator for the approximate
// nearest neighbour search result
// http://www.vldb.org/conf/1999/P49.pdf
func ErrorRatio(ann, groundTruth QueryResult) float64 {
	k := len(groundTruth.Neighbours)
	sort.Sort(groundTruth.Neighbours)
	sort.Sort(ann.Neighbours)
	ratio := 0.0
	for i, n := range groundTruth.Neighbours {
		d := n.Distance
		dAnn := ann.Neighbours[i].Distance
		if d == 0.0 {
			ratio += 1.0
		} else {
			ratio += dAnn / d
		}
	}
	return ratio / float64(k)
}

type AnalysisResult struct {
	QueryIds    []int     `json:"ids"`
	Recalls     []float64 `json:"recalls"`
	ErrorRatios []float64 `json:"errorratios"`
	Times       []float64 `json:"times"`
}

func Analysis(qr, gt QueryResults) (ar *AnalysisResult) {
	sort.Sort(qr)
	sort.Sort(gt)
	ar = &AnalysisResult{
		QueryIds:    make([]int, len(qr)),
		Times:       make([]float64, len(qr)),
		ErrorRatios: make([]float64, len(qr)),
		Recalls:     make([]float64, len(qr)),
	}
	for i := range qr {
		if qr[i].QueryId != gt[i].QueryId {
			panic("Query result and ground truth did not use the same set of queries")
		}
		ar.Recalls[i] = Recall(qr[i], gt[i])
		ar.ErrorRatios[i] = ErrorRatio(qr[i], gt[i])
		ar.QueryIds[i] = qr[i].QueryId
		ar.Times[i] = qr[i].Time
	}
	return ar
}

func RunAnalysis(resultFile, groundTruthFile, output string) {
	var r, g QueryResults
	LoadJson(resultFile, &r)
	LoadJson(groundTruthFile, &g)
	a := Analysis(r, g)
	DumpJson(output, a)
}
