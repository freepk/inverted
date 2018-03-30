package inverted

/*
type Iterator interface {
	func Next() int, bool
	func Reset()
}

Plain(int) Iterator
Union([]Iterator) Iterator
Inter([]Iterarot) Iterator

// (1 || 2 || 3 || 4) && (5 || 6 || 7 || 8)

it := Inter([]Iterator{
	Union( []Iterator{Plain(1), Plain(1), Plain(1), Plain(1)} ),
	Union( []Iterator{Plain(5), Plain(6), Plain(7), Plain(8)} ),
})

it.Reset()
for {
	key, ok := it.Next()
	if !ok {
		break
	}
	...
}
*/
