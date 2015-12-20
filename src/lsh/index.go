package lsh

import "sync"

type QueryFunc func(DataPoint) QueryResult

func ParallelQueryIndex(input []DataPoint, queryFunc QueryFunc, nWorker int) QueryResults {

	// Input Thread
	queries := make(chan DataPoint)
	go func() {
		for _, q := range input {
			queries <- q
		}
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
