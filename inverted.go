package inverted

type Document struct {
	ID     int
	Parts  []int
	Fields [][]int
}

type DocumentIter interface {
	Reset()
	Next() (*Document, bool)
}

type Ref uint32

type Token struct {
	refs  []Ref
	count int
}

func (t *Token) reset() {
	t.refs = t.refs[:0]
	t.count = 0
}

func (t *Token) incr() {
	t.count++
}

func (t *Token) Len() int {
	return len(t.refs)
}

func (t *Token) Count() int {
	return t.count
}

func (t *Token) Iter() *TokenIter {
	return NewTokenIter(t.refs)
}

type Field struct {
	tokens []Token
}

func (f *Field) reset() {
	for i := range f.tokens {
		f.tokens[i].reset()
	}
}

func (f *Field) token(token int) *Token {
	if grow := token + 1 - len(f.tokens); grow > 0 {
		f.tokens = append(f.tokens, make([]Token, grow)...)
	}
	return &f.tokens[token]
}

func (f *Field) Token(token int) *Token {
	if token < len(f.tokens) {
		return &f.tokens[token]
	}
	return nil
}

func (f *Field) Len() int {
	return len(f.tokens)
}

type Part struct {
	fields []Field
}

func (p *Part) reset() {
	for i := range p.fields {
		p.fields[i].reset()
	}
}

func (p *Part) field(field int) *Field {
	if grow := field + 1 - len(p.fields); grow > 0 {
		p.fields = append(p.fields, make([]Field, grow)...)
	}
	return &p.fields[field]
}

func (p *Part) Field(field int) *Field {
	if field < len(p.fields) {
		return &p.fields[field]
	}
	return nil
}

func (p *Part) Len() int {
	return len(p.fields)
}

type Inverted struct {
	iter  DocumentIter
	parts []Part
}

func NewInverted(iter DocumentIter) *Inverted {
	return &Inverted{iter: iter}
}

func (inv *Inverted) reset() {
	for i := range inv.parts {
		inv.parts[i].reset()
	}
}

func (inv *Inverted) part(part int) *Part {
	if grow := part + 1 - len(inv.parts); grow > 0 {
		inv.parts = append(inv.parts, make([]Part, grow)...)
	}
	return &inv.parts[part]
}

func (inv *Inverted) Part(part int) *Part {
	if part < len(inv.parts) {
		return &inv.parts[part]
	}
	return nil
}

func (inv *Inverted) Len() int {
	return len(inv.parts)
}

func (inv *Inverted) walk(do func(*Token, Ref)) {
	inv.iter.Reset()
	for {
		doc, ok := inv.iter.Next()
		if !ok {
			break
		}
		for p := range doc.Parts {
			part := inv.part(doc.Parts[p])
			for f := range doc.Fields {
				field := part.field(f)
				for t := range doc.Fields[f] {
					token := field.token(doc.Fields[f][t])
					do(token, Ref(doc.ID))
				}
			}
		}
	}
}

func (inv *Inverted) Rebuild() {
	inv.reset()
	inv.walk(func(t *Token, r Ref) { t.incr() })
	inv.walk(func(t *Token, r Ref) {
		if t.count > cap(t.refs) {
			n := t.count * 105 / 100
			t.refs = make([]Ref, n)
		}
		t.refs = append(t.refs, r)
	})
}
