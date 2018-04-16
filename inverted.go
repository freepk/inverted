package inverted

import (
	"sort"

	"github.com/freepk/arrays"
)

type Index struct {
	itemTokens    map[int][]int
	updItemTokens map[int][]int
	tokenItems    map[int][]int
}

func NewIndex() *Index {
	return &Index{
		itemTokens:    make(map[int][]int),
		updItemTokens: make(map[int][]int),
		tokenItems:    make(map[int][]int)}
}

func (s *Index) Append(item int, tokens []int) {
	s.updItemTokens[item] = tokens
}

func (s *Index) Item(item int) []int {
	return s.itemTokens[item]
}

func (s *Index) Items(token int) []int {
	return s.tokenItems[token]
}

func (s *Index) update() {
	var ins, del, updItemTokens, updTokenItems map[int][]int
	var item, token, n, c int
	var newList, curList []int
	var ok bool

	del = make(map[int][]int)
	ins = make(map[int][]int)
	updItemTokens, s.updItemTokens = s.updItemTokens, make(map[int][]int)
	for item, newList = range updItemTokens {
		sort.Ints(newList)
		newList = arrays.Distinct(newList)
		updItemTokens[item] = newList
		curList = s.itemTokens[item]
		n = 0
		c = 0
		for (n < len(newList)) && (c < len(curList)) {
			switch {
			case curList[c] < newList[n]:
				token = curList[c]
				del[token] = append(del[token], item)
				c++
			case curList[c] > newList[n]:
				token = newList[n]
				ins[token] = append(ins[token], item)
				n++
			default:
				c++
				n++
			}
		}
		for c < len(curList) {
			token = curList[c]
			del[token] = append(del[token], item)
			c++
		}
		for n < len(newList) {
			token = newList[n]
			ins[token] = append(ins[token], item)
			n++
		}
	}
	for item, curList = range s.itemTokens {
		if _, ok = updItemTokens[item]; !ok {
			updItemTokens[item] = curList
		}
	}
	updTokenItems = make(map[int][]int)
	for token, newList = range del {
		curList, ok = updTokenItems[token]
		if !ok {
			curList = append(curList, s.tokenItems[token]...)
		}
		sort.Ints(newList)
		curList = arrays.Except(curList, newList)
		updTokenItems[token] = curList
	}
	for token, newList = range ins {
		curList, ok = updTokenItems[token]
		if !ok {
			curList = append(curList, s.tokenItems[token]...)
		}
		curList = append(curList, newList...)
		sort.Ints(curList)
		updTokenItems[token] = curList
	}
	for token, curList = range s.tokenItems {
		if _, ok = updTokenItems[token]; !ok {
			updTokenItems[token] = curList
		}
	}
	s.itemTokens = updItemTokens
	s.tokenItems = updTokenItems
}
