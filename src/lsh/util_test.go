package lsh

import (
	"testing"
)

const (
	path = "../data/tiny_images_small.bin"
)

func floatToIntPoint(p Point) []int {
	q := make([]int, len(p))
	for i := range p {
		q[i] = int(p[i])
	}
	return q
}

func Test_CountPoint(t *testing.T) {
	n := CountPoint(path, 3072)
	if n != 100 {
		t.Error("Should have 100 points in the test dataset")
	}
	t.Log(n)
}

func Test_PointIterator(t *testing.T) {
	parser := NewTinyImagePointParser()
	n := CountPoint(path, parser.ByteLen)
	it := NewDataPointIterator(path, parser)
	p, err := it.Next()
	for err == nil {
		t.Log(p)
		p, err = it.Next()
	}
	ids := SelectQueries(n, 10)
	it = NewQueryPointIterator(path, parser, ids)
	p, err = it.Next()
	for err == nil {
		t.Log(p)
		p, err = it.Next()
	}
}
