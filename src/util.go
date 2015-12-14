package lsh

import (
	"encoding/binary"
	"errors"
	"math/rand"
	"os"
	"sort"
)

const (
	randomSeed = 1
	dim        = 3072
	intsize    = 1
)

// SelectQueries returns ids of randomly selected queries
// n is the total number data points
// nq is the number of queries to select
func SelectQueries(n, nq int) []int {
	random := rand.New(rand.NewSource(randomSeed))
	return random.Perm(n)[:nq]
}

// CountPoint return the number of points stored in the serialized
// data file
func CountPoint(path string) int {
	f, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	fi, err := f.Stat()
	if err != nil {
		panic(err.Error())
	}
	filesize := fi.Size()
	unitByteSize := dim * intsize
	if int(filesize)%unitByteSize != 0 {
		panic("Corrupt data file")
	}
	return int(filesize) / unitByteSize
}

type PointIterator struct {
	dim     int   // number of dimensions of a point
	intsize int   // the size of an integer in the serialized format in bytes
	indices []int // indices of the points to visit
	curr    int   // the current index of the indices
	file    *os.File
	path    string
}

// NewQueryPointIterator returns an iterator for all the query points
// in the data file.
// indices are the indices of the queries in the data file
func NewQueryPointIterator(path string, indices []int) *PointIterator {
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	file.Seek(0, 0)
	sort.Ints(indices)
	return &PointIterator{
		dim:     dim,
		intsize: intsize,
		indices: indices,
		curr:    0,
		file:    file,
		path:    path,
	}
}

// NewDataPointIterator returns an iterator for all the points
// in the data file
func NewDataPointIterator(path string) *PointIterator {
	n := CountPoint(path)
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
		dim:     dim,
		intsize: intsize,
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
	unitByteSize := it.dim * it.intsize
	b := make([]byte, unitByteSize)
	_, err := it.file.ReadAt(b, int64(it.indices[it.curr]*unitByteSize))
	if err != nil {
		panic(err.Error())
	}
	// Parse the bytes into a Point
	p := make(Point, it.dim)
	for i := range p {
		switch it.intsize {
		case 1:
			p[i] = float64(int(b[i]))
			break
		case 2:
			p[i] = float64(int(binary.LittleEndian.Uint16(b[i*2 : (i+1)*2])))
			break
		case 4:
			p[i] = float64(int(binary.LittleEndian.Uint32(b[i*4 : (i+1)*4])))
			break
		default:
			panic("Do not support intsize")
		}
	}
	it.curr += 1
	return p, nil
}
