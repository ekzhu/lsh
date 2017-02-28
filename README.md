# LSH for Go

[![Build Status](https://travis-ci.org/ekzhu/lsh.svg?branch=master)](https://travis-ci.org/ekzhu/lsh)
[![GoDoc](https://godoc.org/github.com/ekzhu/lsh?status.svg)](https://godoc.org/github.com/ekzhu/lsh)
[![DOI](https://zenodo.org/badge/50131034.svg)](https://zenodo.org/badge/latestdoi/50131034)

[Documentation](https://godoc.org/github.com/ekzhu/lsh)

Install: `go get github.com/ekzhu/lsh`

This library includes various Locality Sensitive Hashing (LSH) algorithms
for the approximate nearest neighbour search problem in L2 metric space.
The family of LSH functions for L2 is the work of 
[Mayur Datar et.al.](http://www.cs.princeton.edu/courses/archive/spr05/cos598E/bib/p253-datar.pdf)

Currently includes:

* [Basic LSH](http://www.vldb.org/conf/1999/P49.pdf)
* [Multi-probe LSH](http://www.cs.princeton.edu/cass/papers/mplsh_vldb07.pdf)
* [LSH Forest](http://infolab.stanford.edu/~bawa/Pub/similarity.pdf)
