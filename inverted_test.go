package inverted

import (
	"testing"

	"github.com/freepk/arrays"
)

func TestUpdate(t *testing.T) {
	index := NewIndex()
	index.Set(300, []int{1000, 3000, 2000})
	index.Set(200, []int{1000, 2000, 3000})
	index.Set(400, []int{1000, 3000, 2000})
	index.Set(100, []int{2000, 3000, 1000})
	index.Update()
	if !arrays.IsEqual(index.Toks(100), []int{1000, 2000, 3000}) {
		t.Fail()
	}
	if !arrays.IsEqual(index.Docs(1000), []int{100, 200, 300, 400}) {
		t.Fail()
	}
	index.Set(300, []int{1000, 2000, 2000, 2000, 3000, 4000})
	index.Set(200, []int{4000, 1000, 2000, 1000, 2000, 3000, 6000})
	index.Set(200, []int{1000, 6000, 2000, 3000})
	index.Update()
	if !arrays.IsEqual(index.Toks(200), []int{1000, 2000, 3000, 6000}) {
		t.Fail()
	}
	if !arrays.IsEqual(index.Toks(300), []int{1000, 2000, 3000, 4000}) {
		t.Fail()
	}
	if !arrays.IsEqual(index.Docs(4000), []int{300}) {
		t.Fail()
	}
	if !arrays.IsEqual(index.Docs(6000), []int{200}) {
		t.Fail()
	}
	index.Set(300, []int{4000, 1000, 1000, 1000, 5000})
	index.Update()
	if !arrays.IsEqual(index.Toks(300), []int{1000, 4000, 5000}) {
		t.Fail()
	}
	if !arrays.IsEqual(index.Docs(2000), []int{100, 200, 400}) {
		t.Fail()
	}
	if !arrays.IsEqual(index.Docs(5000), []int{300}) {
		t.Fail()
	}
}
