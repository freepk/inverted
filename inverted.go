package inverted

import (
	"sort"

	"github.com/freepk/arrays"
)

type Index struct {
	docToks    map[int][]int
	tokDocs    map[int][]int
	newDocToks map[int][]int
}

func NewIndex() *Index {
	return &Index{
		docToks:    make(map[int][]int),
		tokDocs:    make(map[int][]int),
		newDocToks: make(map[int][]int)}
}

func (i *Index) Append(doc int, toks []int) {
	i.newDocToks[doc] = toks
}

func (i *Index) Doc(doc int) []int {
	return i.docToks[doc]
}

func (i *Index) Docs(tok int) []int {
	return i.tokDocs[tok]
}

func copyMap(d, s map[int][]int) {
	for k, v := range s {
		if _, ok := d[k]; !ok {
			d[k] = v
		}
	}
}

func (i *Index) Update() {
	docToks := i.newDocToks
	i.newDocToks = make(map[int][]int)
	delTokDocs := make(map[int][]int)
	insTokDocs := make(map[int][]int)
	for doc, newToks := range docToks {
		sort.Ints(newToks)
		newToks = arrays.Distinct(newToks)
		docToks[doc] = newToks
		oldToks := i.docToks[doc]
		ins := 0
		del := 0
		for (ins < len(newToks)) && (del < len(oldToks)) {
			switch {
			case oldToks[del] < newToks[ins]:
				tok := oldToks[del]
				delTokDocs[tok] = append(delTokDocs[tok], doc)
				del++
			case oldToks[del] > newToks[ins]:
				tok := newToks[ins]
				insTokDocs[tok] = append(insTokDocs[tok], doc)
				ins++
			default:
				del++
				ins++
			}
		}
		for del < len(oldToks) {
			tok := oldToks[del]
			delTokDocs[tok] = append(delTokDocs[tok], doc)
			del++
		}
		for ins < len(newToks) {
			tok := newToks[ins]
			insTokDocs[tok] = append(insTokDocs[tok], doc)
			ins++
		}
	}
	copyMap(docToks, i.docToks)
	tokDocs := make(map[int][]int)
	for tok, docs := range delTokDocs {
		temp, ok := tokDocs[tok]
		if !ok {
			temp = append(temp, i.tokDocs[tok]...)
		}
		sort.Ints(docs)
		temp = arrays.Except(temp, docs)
		tokDocs[tok] = temp
	}
	for tok, docs := range insTokDocs {
		temp, ok := tokDocs[tok]
		if !ok {
			temp = append(temp, i.tokDocs[tok]...)
		}
		temp = append(temp, docs...)
		sort.Ints(temp)
		tokDocs[tok] = temp
	}
	copyMap(tokDocs, i.tokDocs)
	i.docToks = docToks
	i.tokDocs = tokDocs
}
