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

func (i *Index) Item(key int) []int {
	return i.itemTokens[key]
}

func (i *Index) Items(token int) []int {
	return i.tokenItems[token]
}

func (i *Index) Update() {
	i.update()
}

func (i *Index) update() {
	deleted := make(map[int][]int)
	inserted := make(map[int][]int)
	itemTokens := i.updItemTokens
	i.updItemTokens = make(map[int][]int)
	for key, update := range itemTokens {
		sort.Ints(update)
		update = arrays.Distinct(update)
		itemTokens[key] = update
		current := i.itemTokens[key]
		u := 0
		c := 0
		for (u < len(update)) && (c < len(current)) {
			switch {
			case current[c] < update[u]:
				token := current[c]
				deleted[token] = append(deleted[token], key)
				c++
			case current[c] > update[u]:
				token := update[u]
				inserted[token] = append(inserted[token], key)
				u++
			default:
				c++
				u++
			}
		}
		for c < len(current) {
			token := current[c]
			deleted[token] = append(deleted[token], key)
			c++
		}
		for u < len(update) {
			token := update[u]
			inserted[token] = append(inserted[token], key)
			u++
		}
	}
	for key, current := range i.itemTokens {
		if _, ok := itemTokens[key]; !ok {
			itemTokens[key] = current
		}
	}
	tokenItems := make(map[int][]int)
	for token, update := range deleted {
		current, ok := tokenItems[token]
		if !ok {
			current = append(current, i.tokenItems[token]...)
		}
		sort.Ints(update)
		current = arrays.Except(current, update)
		tokenItems[token] = current
	}
	for token, update := range inserted {
		current, ok := tokenItems[token]
		if !ok {
			current = append(current, i.tokenItems[token]...)
		}
		current = append(current, update...)
		sort.Ints(current)
		tokenItems[token] = current
	}
	for token, current := range i.tokenItems {
		if _, ok := tokenItems[token]; !ok {
			tokenItems[token] = current
		}
	}
	i.itemTokens = itemTokens
	i.tokenItems = tokenItems
}
