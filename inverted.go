package inverted

import (
	"github.com/freepk/radix"
)

const (
	maxTailSize = 4096
)

type sortedArray struct {
	body []int
	tail []int
}

func newSortedArray() *sortedArray {
	return &sortedArray{
		body: make([]int, 0, maxTailSize),
		tail: make([]int, 0, maxTailSize)}
}

func (a *sortedArray) union() {
	bs := len(a.body)
	ts := len(a.tail)
	if ts == 0 {
		return
	}
	buff := make([]int, (bs + ts))
	radix.Ints(a.tail, buff, ts)
	buff = buff[:0]
	prev, curr := 0, 0
	i, j := 0, 0
	for i < bs && j < ts {
		switch {
		case a.body[i] < a.tail[j]:
			curr = a.body[i]
			i++
		case a.body[i] > a.tail[j]:
			curr = a.tail[j]
			j++
		default:
			curr = a.body[i]
			i++
			j++
		}
		if prev != curr {
			prev = curr
			buff = append(buff, curr)
		}
	}
	for i < bs {
		curr = a.body[i]
		i++
		if prev != curr {
			prev = curr
			buff = append(buff, curr)
		}
	}
	for j < ts {
		curr = a.tail[j]
		j++
		if prev != curr {
			prev = curr
			buff = append(buff, curr)
		}
	}
	a.body = buff
	a.tail = a.tail[:0]
}

func (a *sortedArray) append(value int) {
	if len(a.tail) == maxTailSize {
		a.union()
	}
	a.tail = append(a.tail, value)
}

func (a *sortedArray) items() []int {
	a.union()
	return a.body
}

type Index struct {
	tags    map[int]*sortedArray
	trashed int
}

func NewIndex() *Index {
	return &Index{
		tags: make(map[int]*sortedArray)}
}

func (x *Index) Append(key, tag int) {
	if x.tags[tag] == nil {
		x.tags[tag] = newSortedArray()
	}
	x.tags[tag].append(key)
}

func (x *Index) Tags() []int {
	tags := make([]int, 0, len(x.tags))
	for tag := range x.tags {
		tags = append(tags, tag)
	}
	return tags
}

func (x *Index) ByTag(tag int) []int {
	return x.tags[tag].items()
}
