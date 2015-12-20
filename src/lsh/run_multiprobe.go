package lsh

import (
	"sort"
	"time"
)

func RunMultiprobe(data []DataPoint, queries []DataPoint,
	output string,
	k, nQuery, nWorker int,
	dim, m, l int, w float64, t int) {

	lsh := NewMultiprobeLsh(dim, l, m, w, t)
	for _, p := range data {
		lsh.Insert(p.Point, p.Id)
	}
	queryFunc := func(q DataPoint) QueryResult {
		start := time.Now()
		out := make(chan int)
		go func() {
			lsh.QueryK(q.Point, k, out)
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
	results := ParallelQueryIndex(queries, queryFunc, nWorker)
	DumpJson(output, results)
}
