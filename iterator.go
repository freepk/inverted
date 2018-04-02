package inverted

type Iterator interface {
	Next() (int, bool)
	Reset()
}

type ArrayIterator struct {
	array  []int
	offset int
}

func NewArrayIterator(array []int) *ArrayIterator {
	return &ArrayIterator{array: array, offset: 0}
}

func (it *ArrayIterator) Next() (int, bool) {
	if it.offset < len(it.array) {
		offset := it.offset
		it.offset++
		return it.array[offset], true
	}
	return 0, false
}

func (it *ArrayIterator) Reset() {
	it.offset = 0
}

type IntersectIterator struct {
	array []Iterator
}

func NewIntersectIterator(array []Iterator) *IntersectIterator {
	return &IntersectIterator{array: array}
}

func (it *IntersectIterator) Next() (int, bool) {
	var values [32]int
	size := len(it.array)
	if size == 0 {
		return 0, false
	}
	if size == 1 {
		return it.array[0].Next()
	}
	ok := false
	advice := 0
	for i := 0; i < size; i++ {
		values[i], ok = it.array[i].Next()
		if !ok {
			return 0, false
		}
		if i == 0 || advice < values[i] {
			advice = values[i]
		}
	}
	for {
		for i := 0; i < size; i++ {
			if values[i] == advice {
				continue
			}
			for values[i] < advice {
				values[i], ok = it.array[i].Next()
				if !ok {
					return 0, false
				}
			}
		}
		equals := 0
		for i := 0; i < size; i++ {
			if values[i] == advice {
				equals++
			}
		}
		if equals == size {
			return advice, true
		}
	}
	return 0, false
}

func (it *IntersectIterator) Reset() {
	size := len(it.array)
	for i := 0; i < size; i++ {
		it.array[i].Reset()
	}
}
