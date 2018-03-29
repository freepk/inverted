package inverted

import (
	"sync"

	"github.com/freepk/radix"
)

const (
	smallIcebergSize = 1024
)

type iceberg struct {
	top  []int
	body []int
}

func newIceberg() *iceberg {
	body := make([]int, 0, smallIcebergSize)
	return &iceberg{
		top:  body,
		body: body}
}

func (i *iceberg) visible() bool {
	return len(i.top) == len(i.body)
}

// TODO: make sorting in backgound with separate schedule logic
func (i *iceberg) show() {
	if i.visible() {
		return
	}
	size := len(i.body)
	buff := make([]int, size)
	radix.Ints(i.body, buff, size)
	i.top = i.body
}

func (i *iceberg) append(value int) {
	i.body = append(i.body, value)
}

func (i *iceberg) items() []int {
	i.show()
	return i.top
}

type Index struct {
	tokens *sync.Map
}

func NewIndex() *Index {
	return &Index{tokens: &sync.Map{}}
}

func (x *Index) Append(key, token int) {
	list, ok := x.tokens.Load(token)
	if !ok {
		list = newIceberg()
		x.tokens.Store(token, list)
	}
	list.(*iceberg).append(key)
}

func (x *Index) Tokens() []int {
	tokens := make([]int, 0)
	x.tokens.Range(func(k, v interface{}) bool {
		tokens = append(tokens, k.(int))
		return true
	})
	return tokens
}

func (x *Index) ByToken(token int) []int {
	list, ok := x.tokens.Load(token)
	if !ok {
		return []int{}
	}
	return list.(*iceberg).items()
}
