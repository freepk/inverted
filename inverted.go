package inverted

import (
        "encoding/gob"
        "io"
        "os"
        "sync"

        "github.com/freepk/iterator"
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
        tokens *sync.Map
}

func NewIndex() *Index {
        return &Index{tokens: &sync.Map{}}
}

func (x *Index) Append(key, token int) {
        list, ok := x.tokens.Load(token)
        if !ok {
                list = newIceberg(nil)
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

func (x *Index) Vector(token int) *iterator.ArrayIterator {
        list, ok := x.tokens.Load(token)
        if !ok {
                return iterator.NewArrayIterator([]int{})
        }
        return iterator.NewArrayIterator(list.(*iceberg).items())
}

func (x *Index) Dump(path string) error {
        file, err := os.Create(path)
        if err != nil {
                return err
        }
        defer file.Close()
        enc := gob.NewEncoder(file)
        x.tokens.Range(func(k, v interface{}) bool {
                token := k.(int)
                items := v.(*iceberg).items()
                enc.Encode(token)
                enc.Encode(items)
                return true
        })
        return nil
}

func (x *Index) Restore(path string) error {
        file, err := os.Open(path)
        if err != nil {
                return err
        }
        defer file.Close()
        tokens := &sync.Map{}
        token := 0
        items := make([]int, 0)
        dec := gob.NewDecoder(file)
        for {
                err = dec.Decode(&token)
                if err != nil {
                        if err == io.EOF {
                                break
                        }
                        return err
                }
                err = dec.Decode(&items)
                if err != nil {
                        return err
                }
                list := newIceberg(items)
                tokens.Store(token, list)
        }
        x.tokens = tokens
        return nil
}
