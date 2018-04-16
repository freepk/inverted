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

func (i *Index) Append(item int, tokens []int) {
	i.updItemTokens[item] = tokens
}

func (i *Index) Item(item int) []int {
	return i.itemTokens[item]
}

func (i *Index) Items(token int) []int {
	return i.tokenItems[token]
}

func (i *Index) update() {
	var inserted, deleted, itemTokens, tokenItems map[int][]int
	var item, token, u, c int
	var update, current []int
	var ok bool

	deleted = make(map[int][]int)
	inserted = make(map[int][]int)
	itemTokens, i.updItemTokens = i.updItemTokens, make(map[int][]int)
	for item, update = range itemTokens {
		sort.Ints(update)
		update = arrays.Distinct(update)
		itemTokens[item] = update
		current = i.itemTokens[item]
		u = 0
		c = 0
		for (u < len(update)) && (c < len(current)) {
			switch {
			case current[c] < update[u]:
				token = current[c]
				deleted[token] = append(deleted[token], item)
				c++
			case current[c] > update[u]:
				token = update[u]
				inserted[token] = append(inserted[token], item)
				u++
			default:
				c++
				u++
			}
		}
		for c < len(current) {
			token = current[c]
			deleted[token] = append(deleted[token], item)
			c++
		}
		for u < len(update) {
			token = update[u]
			inserted[token] = append(inserted[token], item)
			u++
		}
	}
	for item, current = range i.itemTokens {
		if _, ok = itemTokens[item]; !ok {
			itemTokens[item] = current
		}
	}
	tokenItems = make(map[int][]int)
	for token, update = range deleted {
		current, ok = tokenItems[token]
		if !ok {
			current = append(current, i.tokenItems[token]...)
		}
		sort.Ints(update)
		current = arrays.Except(current, update)
		tokenItems[token] = current
	}
	for token, update = range inserted {
		current, ok = tokenItems[token]
		if !ok {
			current = append(current, i.tokenItems[token]...)
		}
		current = append(current, update...)
		sort.Ints(current)
		tokenItems[token] = current
	}
	for token, current = range i.tokenItems {
		if _, ok = tokenItems[token]; !ok {
			tokenItems[token] = current
		}
	}
	i.itemTokens = itemTokens
	i.tokenItems = tokenItems
}
