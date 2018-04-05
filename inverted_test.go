package inverted

import (
	"sort"
	"testing"

	"github.com/freepk/iterator"
)

func readAll(it iterator.Iterator) []int {
	a := make([]int, 0)
	for {
		v, ok := it.Next()
		if !ok {
			break
		}
		a = append(a, v)
	}
	return a
}

func iteratorsEqual(a, b iterator.Iterator) bool {
	for {
		av, aok := a.Next()
		bv, bok := b.Next()
		if aok != bok {
			return false
		}
		if av != bv {
			return false
		}
		if !aok {
			break
		}
	}
	return true
}

func indexesEqual(a, b *Index) bool {
	atokens := readAll(a.Tokens())
	btokens := readAll(b.Tokens())
	size := len(atokens)
	if size != len(btokens) {
		return false
	}
	sort.Ints(atokens)
	sort.Ints(btokens)
	for i := 0; i < size; i++ {
		atoken, btoken := atokens[i], btokens[i]
		if atoken != btoken {
			return false
		}
		if !iteratorsEqual(a.Items(atokens[i]), b.Items(btokens[i])) {
			return false
		}
	}
	return true
}

func TestIndex(t *testing.T) {
	x := NewIndex()
	x.Append(2000, []int{})
	x.Append(2002, []int{300, 555})
	x.Append(3001, []int{300})
	x.Append(2000, []int{})
	x.Append(4002, []int{100, 300})
	x.Append(1001, []int{300})
	x.Dump("test.dump")
	y := NewIndex()
	y.Restore("test.dump")
	if !indexesEqual(x, y) {
		t.Fail()
	}
}
