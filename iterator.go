package inverted

type Iterator interface {
	Next() (int, bool)
}

type ArrayIterator struct {
	array  []int
	offset int
}

func (it *ArrayIterator) Next() (int, bool) {
	if it.offset < len(it.array) {
		offset := it.offset
		it.offset++
		return it.array[offset], true
	}
	return 0, false
}

type IntersectIterator struct {
	array []Iterator
}

func (it *IntersectIterator) Next() (int, bool) {
	size := len(it.array)
	if size == 0 {
		return 0, false
	}
	if size == 1 {
		return it.array[0].Next()
	}
	ok := false
	advice := 0
	values := make([]int, size)
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
