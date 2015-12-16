package lsh

import (
	"bytes"
	"encoding/binary"
)

const (
	// TinyImage specific constants
	dim     = 3072
	gistDim = 384
)

// NewTinyImagePointParser returns a PointParser
// specific to the TinyImage dataset
// http://horatio.cs.nyu.edu/mit/tiny/data/index.html
func NewTinyImagePointParser() *PointParser {
	return &PointParser{
		ByteLen: dim,
		Parse:   ParseTinyImagePoint,
	}
}

// NewTinyImageGistParser returns a PointParser
// for the GIST descriptor of the TinyImage dataset
func NewTinyImageGistParser() *PointParser {
	return &PointParser{
		ByteLen: gistDim * 4,
		Parse:   ParseTinyImageGist,
	}
}

// ParseTinyImagePoint takes a serialized point vector b
// and returns the parsed Point
func ParseTinyImagePoint(b []byte) Point {
	if len(b) != dim {
		panic("Incorrect input for parsing serialized point")
	}
	p := make(Point, dim)
	for i := range p {
		p[i] = float64(int(b[i]))
	}
	return p
}

func ParseTinyImageGist(b []byte) Point {
	if len(b) != gistDim*4 {
		panic("Incorrect input for parsing serialized GIST")
	}
	buf := bytes.NewReader(b)
	p := make(Point, gistDim)
	for i := range p {
		var x float32
		err := binary.Read(buf, binary.LittleEndian, &x)
		if err != nil {
			panic(err.Error())
		}
		p[i] = float64(x)
	}
	return p
}
