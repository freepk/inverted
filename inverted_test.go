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
	index := testIndex(1000000, 20000000, 40, 44000)
	tokens := index.Tokens()
	for _, token := range tokens {
		keys := index.ByToken(token)
		if !sort.IntsAreSorted(keys) {
			t.Fail()
		}
	}
}
