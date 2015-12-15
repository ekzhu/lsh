package lsh

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
)

const (
	randomSeed = 1
)

type PointParser struct {
	ByteLen int
	Parse   func([]byte) Point
}

// SelectQueries returns ids of randomly selected queries
// n is the total number data points
// nq is the number of queries to select
func SelectQueries(n, nq int) []int {
	random := rand.New(rand.NewSource(randomSeed))
	return random.Perm(n)[:nq]
}

// CountPoint return the number of points stored in the serialized
// data file
func CountPoint(path string, byteLen int) int {
	f, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	fi, err := f.Stat()
	if err != nil {
		panic(err.Error())
	}
	filesize := fi.Size()
	if int(filesize)%byteLen != 0 {
		panic("Corrupt data file")
	}
	err = f.Close()
	if err != nil {
		panic(err.Error())
	}
	return int(filesize) / byteLen
}

type PointIterator struct {
	parser  *PointParser
	indices []int // indices of the points to visit
	curr    int   // the current index of the indices
	file    *os.File
	path    string
}

// NewQueryPointIterator returns an iterator for all the query points
// in the data file.
// indices are the indices of the queries in the data file
func NewQueryPointIterator(path string, parser *PointParser, indices []int) *PointIterator {
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	file.Seek(0, 0)
	sort.Ints(indices)
	return &PointIterator{
		parser:  parser,
		indices: indices,
		curr:    0,
		file:    file,
		path:    path,
	}
}

// NewDataPointIterator returns an iterator for all the points
// in the data file
func NewDataPointIterator(path string, parser *PointParser) *PointIterator {
	n := CountPoint(path, parser.ByteLen)
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	file.Seek(0, 0)
	sort.Ints(indices)
	return &PointIterator{
		parser:  parser,
		indices: indices,
		curr:    0,
		file:    file,
		path:    path,
	}
}

// Next returns the next point in the data file
func (it *PointIterator) Next() (Point, error) {
	if len(it.indices) == it.curr {
		return nil, errors.New("Empty result")
	}
	b := make([]byte, it.parser.ByteLen)
	_, err := it.file.ReadAt(b, int64(it.indices[it.curr]*it.parser.ByteLen))
	if err != nil {
		panic(err.Error())
	}
	// Parse the bytes into a Point
	p := it.parser.Parse(b)
	it.curr += 1
	return p, nil
}

// Close releases resources used by the iterator
func (it *PointIterator) Close() {
	err := it.file.Close()
	if err != nil {
		panic(err.Error())
	}
	it.indices = nil
}

func LoadJson(file string, v interface{}) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(buffer, v)
	if err != nil {
		panic(err.Error())
	}
}

func DumpJson(file string, v interface{}) {
	buffer, err := json.Marshal(v)
	if err != nil {
		panic(err.Error())
	}
	err = ioutil.WriteFile(file, buffer, 0777)
	if err != nil {
		panic(err.Error())
	}
}
