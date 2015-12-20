package lsh

import (
	"sort"
	"time"
)

func RunMultiprobe(datafile, output string,
	k, nQuery, nWorker int,
	parser *PointParser,
	dim, m, l int, w float64, t int) {

	// Load data
	nData := CountPoint(datafile, parser.ByteLen)
	iter := NewDataPointIterator(datafile, parser)
	data := make([]Point, nData)
	ids := make([]int, nData)
	for i := 0; i < nData; i++ {
		p, err := iter.Next()
		if err != nil {
			panic(err.Error())
		}
		data[i] = p.Point
		ids[i] = p.Id
	}

	lsh := NewMultiprobeLsh(dim, l, m, w, t)
	for i, p := range data {
		lsh.Insert(p, ids[i])
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
				Id:       r[i],
				Distance: q.Point.L2(data[i]),
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
	// Select queries
	queryIds := SelectQueries(nData, nQuery)
	iter = NewQueryPointIterator(datafile, parser, queryIds)
	// Run queries in parallel
	results := ParallelQueryIndex(iter, queryFunc, nWorker)
	DumpJson(output, results)
}
