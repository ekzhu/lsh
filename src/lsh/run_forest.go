package lsh

import (
	"sort"
	"time"
)

func RunForest(data []DataPoint, queries []DataPoint,
	output string,
	k, nWorker int,
	dim, m, l int, w float64) {

	// Build forest index
	forest := NewLshForest(dim, l, m, w)
	for _, p := range data {
		forest.Insert(p.Point, p.Id)
	}
	// Forest query function wrapper
	queryFunc := func(q DataPoint) QueryResult {
		start := time.Now()
		out := make(chan int)
		go func() {
			forest.QueryK(q.Point, k, out)
			close(out)
		}()
		r := make([]int, 0)
		for i := range out {
			r = append(r, i)
		}
		dur := time.Since(start)
		ns := make(Neighbours, len(r))
		for i := range r {
			ns[i] = Neighbour{
				Id: r[i],
				// We assume the id is equal to the index
				// of the data point in the input data
				Distance: q.Point.L2(data[r[i]].Point),
			}
		}
		sort.Sort(ns)
		if len(ns) > k {
			ns = ns[:k]
		}
		return QueryResult{
			QueryId:    q.Id,
			Neighbours: ns,
			Time:       float64(dur) / float64(time.Millisecond),
		}
	}
	// Run queries in parallel
	results := ParallelQueryIndex(queries, queryFunc, nWorker)
	DumpJson(output, results)
}
