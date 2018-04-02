package inverted

import (
	"hash/fnv"
	"testing"
)

var (
	testIndex     *Index
	testIndexPath = "index.dump"
	fnvHash64     = fnv.New64()
)

func init() {
	testIndex = NewIndex()
	testIndex.Restore(testIndexPath)
}

func TestArrayIterator(t *testing.T) {
	iter := NewArrayIterator([]int{1, 2, 3})
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
	iter := NewIntersectIterator([]Iterator{
		NewArrayIterator([]int{0, 1, 2, 3}),
		NewArrayIterator([]int{1, 2, 3, 4}),
	})
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

func stringToInt(s string) int {
	fnvHash64.Reset()
	fnvHash64.Write([]byte(s))
	return int(fnvHash64.Sum64())
}

func BenchmarkIntersect(b *testing.B) {
	b.StopTimer()
	brand := stringToInt("B14426")
	group := stringToInt("G850")
	iterator := NewIntersectIterator([]Iterator{
		testIndex.Vector(brand),
		testIndex.Vector(group)})
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		iterator.Reset()
		for {
			value, ok := iterator.Next()
			if !ok {
				break
			}
			_ = value
		}
	}
}
