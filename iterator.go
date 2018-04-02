package inverted

type Iterator interface {
	Reset()
	Next() (int, bool)
}

type ArrayIterator struct {
	items  []int
	offset int
}

func NewArrayIterator(items []int) *ArrayIterator {
	return &ArrayIterator{items: items, offset: 0}
}

func (it *ArrayIterator) Reset() {
	it.offset = 0
}

func (it *ArrayIterator) Next() (int, bool) {
	if it.offset < len(it.items) {
		offset := it.offset
		it.offset++
		return it.items[offset], true
	}
	return 0, false
}

type IntersectIterator struct {
	iterators []Iterator
	values    []int
}

func NewIntersectIterator(iterators []Iterator) *IntersectIterator {
	return &IntersectIterator{
		iterators: iterators,
		values:    make([]int, len(iterators))}
}

func (it *IntersectIterator) Reset() {
	size := len(it.iterators)
	for i := 0; i < size; i++ {
		it.iterators[i].Reset()
	}
}

func (it *IntersectIterator) Next() (int, bool) {
	size := len(it.iterators)
	if size == 0 {
		return 0, false
	}
	if size == 1 {
		return it.iterators[0].Next()
	}
	ok := false
	advice := 0
	for i := 0; i < size; i++ {
		it.values[i], ok = it.iterators[i].Next()
		if !ok {
			return 0, false
		}
		if i == 0 || advice < it.values[i] {
			advice = it.values[i]
		}
	}
	for {
		for i := 0; i < size; i++ {
			if it.values[i] == advice {
				continue
			}
			for it.values[i] < advice {
				it.values[i], ok = it.iterators[i].Next()
				if !ok {
					return 0, false
				}
			}
		}
		equals := 0
		for i := 0; i < size; i++ {
			if it.values[i] == advice {
				equals++
			}
		}
		if equals == size {
			return advice, true
		}
	}
	return 0, false
}
