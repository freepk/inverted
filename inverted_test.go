package inverted

import (
	"math/rand"
	"sort"
	"testing"
)

func testIndex(keys, maxKey, tags, maxTag int) *Index {
	index := NewIndex()
	for i := 0; i < keys; i++ {
		key := rand.Intn(maxKey)
		for j := 0; j < tags; j++ {
			tag := rand.Intn(maxTag)
			index.Append(key, tag)
		}
	}
	return index
}

func TestIndex(t *testing.T) {
	index := testIndex(10, 200000, 10, 10)
	tags := index.Tags()
	for _, tag := range tags {
		keys := index.ByTag(tag)
		if !sort.IntsAreSorted(keys) {
			t.Fail()
		}
	}
}
