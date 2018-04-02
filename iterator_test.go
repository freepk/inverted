package inverted

import "testing"

func TestArrayIterator(t *testing.T) {
	iter := ArrayIterator{array: []int{1, 2, 3}}
	v, ok := iter.Next()
	if !ok || v != 1 {
		t.Fail()
	}
	v, ok = iter.Next()
	if !ok || v != 2 {
		t.Fail()
	}
	v, ok = iter.Next()
	if !ok || v != 3 {
		t.Fail()
	}
	v, ok = iter.Next()
	if ok || v != 0 {
		t.Fail()
	}
}

func TestIntersectIterator(t *testing.T) {
	a := &ArrayIterator{array: []int{0, 1, 2, 3}}
	b := &ArrayIterator{array: []int{1, 2, 3, 4}}
	iter := IntersectIterator{array: []Iterator{a, b}}
	v, ok := iter.Next()
	if !ok || v != 1 {
		t.Fail()
	}
	v, ok = iter.Next()
	if !ok || v != 2 {
		t.Fail()
	}
	v, ok = iter.Next()
	if !ok || v != 3 {
		t.Fail()
	}
	v, ok = iter.Next()
	if ok || v != 0 {
		t.Fail()
	}
}
