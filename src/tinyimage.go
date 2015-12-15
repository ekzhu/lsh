package lsh

const (
	// TinyImage specific constants
	dim = 3072
)

// NewTinyImagePointParser returns a PointParser
// specific to the TinyImage dataset
// http://horatio.cs.nyu.edu/mit/tiny/data/index.html
func NewTinyImagePointParser() *PointParser {
	return &PointParser{
		byteLen: dim,
		parse:   ParseTinyImagePoint,
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
