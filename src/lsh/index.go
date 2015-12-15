package lsh

import "sync"

type QueryFunc func(DataPoint) QueryResult

func ParallelQueryIndex(queryIter *PointIterator, queryFunc QueryFunc,
	nWorker int) QueryResults {

	// Input Thread
	queries := make(chan DataPoint)
	go func() {
		p, err := queryIter.Next()
		for err == nil {
			queries <- p
			p, err = queryIter.Next()
		}
		queryIter.Close()
		close(queries)
	}()

	var wg sync.WaitGroup
	// Worker threads will write results to this channel
	// before they exit
	queryResults := make(chan QueryResult)
	for i := 0; i < nWorker; i++ {
		wg.Add(1)
		go func(workerId int) {
			for q := range queries {
				r := queryFunc(q)
				queryResults <- r
			}
			wg.Done()
		}(i)
	}
	// Waiting thread for the workers, close the output channel when
	// all workers exit
	go func() {
		wg.Wait()
		close(queryResults)
	}()
	// Merge all query results from workers
	completeQueryResults := make(QueryResults, 0)
	for r := range queryResults {
		completeQueryResults = append(completeQueryResults, r)
	}
	return completeQueryResults
}
