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
	docTokens := i.newDocTokens
	i.newDocTokens = make(map[int][]int)
	remove := make(map[int][]int)
	insert := make(map[int][]int)
	for key, tokens := range docTokens {
		sort.Ints(tokens)
		tokens = arrays.Distinct(tokens)
		docTokens[key] = tokens
		temp := i.docTokens[key]
		i := 0
		r := 0
		for (i < len(tokens)) && (r < len(temp)) {
			switch {
			case temp[r] < tokens[i]:
				token := temp[r]
				remove[token] = append(remove[token], key)
				r++
			case temp[r] > tokens[i]:
				token := tokens[i]
				insert[token] = append(insert[token], key)
				i++
			default:
				r++
				i++
			}
		}
		for r < len(temp) {
			token := temp[r]
			remove[token] = append(remove[token], key)
			r++
		}
		for i < len(tokens) {
			token := tokens[i]
			insert[token] = append(insert[token], key)
			i++
		}
	}
	for key, tokens := range i.docTokens {
		if _, ok := docTokens[key]; !ok {
			docTokens[key] = tokens
		}
	}
	tokenDocs := make(map[int][]int)
	for token, docs := range remove {
		temp, ok := tokenDocs[token]
		if !ok {
			temp = append(temp, i.tokenDocs[token]...)
		}
		sort.Ints(docs)
		temp = arrays.Except(temp, docs)
		tokenDocs[token] = temp
	}
	for token, docs := range insert {
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
