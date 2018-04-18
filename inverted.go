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

func (i *Index) Append(docId int, tokenIds []int) {
	i.newDocTokens[docId] = tokenIds
}

func (i *Index) Doc(docId int) []int {
	return i.docTokens[docId]
}

func (i *Index) Docs(tokenId int) []int {
	return i.tokenDocs[tokenId]
}

func copyNotExists(dst, src map[int][]int) {
	for k, v := range src {
		if _, ok := dst[k]; !ok {
			dst[k] = v
		}
	}
}

func (i *Index) Update() {
	docTokens := i.newDocTokens
	i.newDocTokens = make(map[int][]int)
	deleteIds := make(map[int][]int)
	insertIds := make(map[int][]int)
	for docId, newTokenIds := range docTokens {
		sort.Ints(newTokenIds)
		newTokenIds = arrays.Distinct(newTokenIds)
		docTokens[docId] = newTokenIds
		oldTokenIds := i.docTokens[docId]
		ins := 0
		del := 0
		for (ins < len(newTokenIds)) && (del < len(oldTokenIds)) {
			switch {
			case oldTokenIds[del] < newTokenIds[ins]:
				tokenId := oldTokenIds[del]
				deleteIds[tokenId] = append(deleteIds[tokenId], docId)
				del++
			case oldTokenIds[del] > newTokenIds[ins]:
				tokenId := newTokenIds[ins]
				insertIds[tokenId] = append(insertIds[tokenId], docId)
				ins++
			default:
				del++
				ins++
			}
		}
		for del < len(oldTokenIds) {
			tokenId := oldTokenIds[del]
			deleteIds[tokenId] = append(deleteIds[tokenId], docId)
			del++
		}
		for ins < len(newTokenIds) {
			tokenId := newTokenIds[ins]
			insertIds[tokenId] = append(insertIds[tokenId], docId)
			ins++
		}
	}
	copyNotExists(docTokens, i.docTokens)
	tokenDocs := make(map[int][]int)
	for tokenId, docIds := range deleteIds {
		temp, ok := tokenDocs[tokenId]
		if !ok {
			temp = append(temp, i.tokenDocs[tokenId]...)
		}
		sort.Ints(docIds)
		temp = arrays.Except(temp, docIds)
		tokenDocs[tokenId] = temp
	}
	for tokenId, docIds := range insertIds {
		temp, ok := tokenDocs[tokenId]
		if !ok {
			temp = append(temp, i.tokenDocs[tokenId]...)
		}
		temp = append(temp, docIds...)
		sort.Ints(temp)
		tokenDocs[tokenId] = temp
	}
	copyNotExists(tokenDocs, i.tokenDocs)
	i.docTokens = docTokens
	i.tokenDocs = tokenDocs
}
