package inverted

import (
	"sort"

	"github.com/freepk/arrays"
)

type Index struct {
	docTokens    map[int][]int
	tokenDocs    map[int][]int
	newDocTokens map[int][]int
}

func NewIndex() *Index {
	return &Index{
		docTokens:    make(map[int][]int),
		tokenDocs:    make(map[int][]int),
		newDocTokens: make(map[int][]int)}
}

func (i *Index) Append(key int, tokens []int) {
	i.newDocTokens[key] = tokens
}

func (i *Index) Doc(key int) []int {
	return i.docTokens[key]
}

func (i *Index) Docs(token int) []int {
	return i.tokenDocs[token]
}

func (i *Index) Update() {
	del := make(map[int][]int)
	ins := make(map[int][]int)
	docTokens := i.newDocTokens
	i.newDocTokens = make(map[int][]int)
	for key, tokens := range docTokens {
		sort.Ints(tokens)
		tokens = arrays.Distinct(tokens)
		docTokens[key] = tokens
		temp := i.docTokens[key]
		i := 0
		j := 0
		for (i < len(tokens)) && (j < len(temp)) {
			switch {
			case temp[j] < tokens[i]:
				token := temp[j]
				del[token] = append(del[token], key)
				j++
			case temp[j] > tokens[i]:
				token := tokens[i]
				ins[token] = append(ins[token], key)
				i++
			default:
				j++
				i++
			}
		}
		for j < len(temp) {
			token := temp[j]
			del[token] = append(del[token], key)
			j++
		}
		for i < len(tokens) {
			token := tokens[i]
			ins[token] = append(ins[token], key)
			i++
		}
	}
	for key, tokens := range i.docTokens {
		if _, ok := docTokens[key]; !ok {
			docTokens[key] = tokens
		}
	}
	tokenDocs := make(map[int][]int)
	for token, docs := range del {
		temp, ok := tokenDocs[token]
		if !ok {
			temp = append(temp, i.tokenDocs[token]...)
		}
		sort.Ints(docs)
		temp = arrays.Except(temp, docs)
		tokenDocs[token] = temp
	}
	for token, docs := range ins {
		temp, ok := tokenDocs[token]
		if !ok {
			temp = append(temp, i.tokenDocs[token]...)
		}
		temp = append(temp, docs...)
		sort.Ints(temp)
		tokenDocs[token] = temp
	}
	for token, docs := range i.tokenDocs {
		if _, ok := tokenDocs[token]; !ok {
			tokenDocs[token] = docs
		}
	}
	i.docTokens = docTokens
	i.tokenDocs = tokenDocs
}
