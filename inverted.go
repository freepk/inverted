package inverted

import (
	"encoding/gob"
	"io"
	"os"

	"github.com/freepk/radix"
)

const (
	smallIcebergSize = 1024
)

type iceberg struct {
	top  []int
	body []int
}

func newIceberg(body []int) *iceberg {
	if body == nil {
		body = make([]int, 0, smallIcebergSize)
	}
	return &iceberg{
		top:  body,
		body: body}
}

func (i *iceberg) visible() bool {
	return len(i.top) == len(i.body)
}

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
	tokens map[int]*iceberg
}

func NewIndex() *Index {
	return &Index{tokens: make(map[int]*iceberg)}
}

func (x *Index) Append(key int, tokens []int) {
	size := len(tokens)
	for i := 0; i < size; i++ {
		items, ok := x.tokens[tokens[i]]
		if !ok {
			items = newIceberg([]int{})
			x.tokens[tokens[i]] = items
		}
		items.append(key)
	}
}

func (x *Index) Items(token int) []int {
	items, ok := x.tokens[token]
	if !ok {
		return []int{}
	}
	return items.items()
}

func (x *Index) Tokens() []int {
	tokens := make([]int, 0, len(x.tokens))
	for token := range x.tokens {
		tokens = append(tokens, token)
	}
	return tokens
}

func (x *Index) Dump(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	for token, items := range x.tokens {
		err = enc.Encode(token)
		if err != nil {
			return err
		}
		err = enc.Encode(items.items())
		if err != nil {
			return err
		}
	}
	return nil
}

func (x *Index) Restore(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	dec := gob.NewDecoder(file)
	tokens := make(map[int]*iceberg)
	for {
		token := 0
		items := []int{}
		err = dec.Decode(&token)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = dec.Decode(&items)
		if err != nil {
			return err
		}
		tokens[token] = newIceberg(items)
	}
	x.tokens = tokens
	return nil
}
