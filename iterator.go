package inverted

type TokenIter struct {
	a []Ref
	i int
}

func NewTokenIter(a []Ref) *TokenIter {
	return &TokenIter{a: a, i: 0}
}

func (it *TokenIter) Reset() {
	it.i = 0
}

func (it *TokenIter) Next() (int, bool) {
	i := it.i
	if i < len(it.a) {
		it.i++
		return int(it.a[i]), true
	}
	return 0, false
}
