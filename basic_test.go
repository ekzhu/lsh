package lsh

import (
	"strconv"
	"testing"
)

func Test_NewBasicLsh(t *testing.T) {
	lsh := NewBasicLsh(5, 5, 100, 5.0)
	if len(lsh.tables) != 5 {
		t.Error("Lsh init fail")
	}
}

func Test_Insert(t *testing.T) {
	lsh := NewBasicLsh(100, 5, 5, 5.0)
	points := randomPoints(10, 100, 32.0)
	for i, p := range points {
		lsh.Insert(p, strconv.Itoa(i))
	}
	for _, table := range lsh.tables {
		if len(table) == 0 {
			t.Error("Insert fail")
		}
	}
}

func Test_Query(t *testing.T) {
	lsh := NewBasicLsh(100, 5, 5, 5.0)
	points := randomPoints(10, 100, 32.0)
	insertedKeys := make([]string, 10)
	for i, p := range points {
		lsh.Insert(p, strconv.Itoa(i))
		insertedKeys[i] = strconv.Itoa(i)
	}
	// Use the inserted points as queries, and
	// verify that we can get back each query itself
	for i, key := range insertedKeys {
		found := false
		for _, foundKey := range lsh.Query(points[i]) {
			if foundKey == key {
				found = true
			}
		}
		if !found {
			t.Error("Query fail")
		}
	}
}

func Test_Delete(t *testing.T) {
	lsh := NewBasicLsh(100, 5, 5, 5.0)
	points := randomPoints(10, 100, 32.0)
	for i, p := range points {
		lsh.Insert(p, strconv.Itoa(i))
	}
	for i, p := range points {
		lsh.Delete(strconv.Itoa(i))
		for _, table := range lsh.tables {
			if len(table) != len(points)-(i+1) {
				t.Errorf("Failed to delete point %v. Expected to have %v points, found %v.", i, len(points)-(i+1), len(table))
			}
		}
		found := false
		for _, foundKey := range lsh.Query(p) {
			if foundKey == strconv.Itoa(i) {
				found = true
			}
		}
		if found {
			t.Errorf("Failed to delete point %v.", i)
		}
	}
	Test_Insert(t)
}
