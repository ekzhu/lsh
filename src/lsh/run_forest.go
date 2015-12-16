package lsh

import (
	"sort"
	"time"
)

func RunForest(datafile, output string,
	k, nQuery, nWorker int,
	parser *PointParser,
	dim, m, l int, w float64) {

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

	// Build forest index
	forest := NewLshForest(dim, l, m, w)
	for i, p := range data {
		forest.Insert(p, ids[i])
	}
	// Forest query function wrapper
	queryFunc := func(q DataPoint) QueryResult {
		start := time.Now()
		out := make(chan int)
		go func() {
			forest.Query(q.Point, out)
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
