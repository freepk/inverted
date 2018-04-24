package inverted

import (
	"sort"

	"github.com/freepk/arrays"
)

type Index struct {
	docs map[int][]int
	toks map[int][]int
	upds map[int][]int
}

func NewIndex() *Index {
	return &Index{
		docs: make(map[int][]int),
		toks: make(map[int][]int),
		upds: make(map[int][]int)}
}

func (i *Index) Set(doc int, toks []int) {
	i.upds[doc] = toks
}

func (i *Index) Toks(doc int) []int {
	return i.docs[doc]
}

func (i *Index) Docs(tok int) []int {
	return i.toks[tok]
}

func merge(d, s map[int][]int) {
	for k, v := range s {
		if _, ok := d[k]; !ok {
			d[k] = v
		}
	}
}

func (i *Index) Update() {
	docs := i.upds
	i.upds = make(map[int][]int)
	ptmp := make(map[int][]int)
	ntmp := make(map[int][]int)
	for doc, ntoks := range docs {
		sort.Ints(ntoks)
		ntoks = arrays.Distinct(ntoks)
		docs[doc] = ntoks
		ptoks := i.docs[doc]
		n := 0
		p := 0
		for (n < len(ntoks)) && (p < len(ptoks)) {
			ptok := ptoks[p]
			ntok := ntoks[n]
			if ptok < ntok {
				ptmp[ptok] = append(ptmp[ptok], doc)
				p++
				continue
			}
			if ptok > ntok {
				ntmp[ntok] = append(ntmp[ntok], doc)
				n++
				continue
			}
			p++
			n++
		}
		for p < len(ptoks) {
			ptok := ptoks[p]
			ptmp[ptok] = append(ptmp[ptok], doc)
			p++
		}
		for n < len(ntoks) {
			ntok := ntoks[n]
			ntmp[ntok] = append(ntmp[ntok], doc)
			n++
		}
	}
	merge(docs, i.docs)
	toks := make(map[int][]int)
	for tok, docs := range ptmp {
		tmp, ok := toks[tok]
		if !ok {
			tmp = append(tmp, i.toks[tok]...)
		}
		sort.Ints(docs)
		tmp = arrays.Except(tmp, docs)
		toks[tok] = tmp
	}
	for tok, docs := range ntmp {
		tmp, ok := toks[tok]
		if !ok {
			tmp = append(tmp, i.toks[tok]...)
		}
		tmp = append(tmp, docs...)
		sort.Ints(tmp)
		toks[tok] = tmp
	}
	merge(toks, i.toks)
	i.docs = docs
	i.toks = toks
}
