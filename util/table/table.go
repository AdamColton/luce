package table

type Index struct {
	Row, Col int
}

func (i Index) UpdateSize(i2 Index) Index {
	if i2.Col >= i.Col {
		i.Col = i2.Col + 1
	}
	if i2.Row >= i.Row {
		i.Row = i2.Row + 1
	}
	return i
}

type Iter struct {
	Cur, Size *Index
}

func (ii *Iter) Next() (Index, bool) {
	if ii.Cur == nil {
		ii.Cur = &Index{}
		return *ii.Cur, false
	}
	if ii.Cur.Row >= ii.Size.Row {
		return Index{-1, -1}, true
	}
	ii.Cur.Col++
	if ii.Cur.Col >= ii.Size.Col {
		ii.Cur.Col = 0
		ii.Cur.Row++
		if ii.Cur.Row >= ii.Size.Row {
			return Index{-1, -1}, true
		}
	}
	return *(ii.Cur), false
}

func (ii *Iter) Idx() int {
	return ii.Cur.Col + ii.Cur.Row*ii.Size.Col
}

type Table[T any] struct {
	Data   map[Index]T
	Size   Index
	Labels struct {
		Rows, Cols []string
	}
}

func New[T any]() *Table[T] {
	return &Table[T]{
		Data: make(map[Index]T),
	}
}

func (t *Table[T]) Add(r, c int, cell T) {
	i := Index{r, c}
	t.Data[i] = cell
	t.Size = t.Size.UpdateSize(i)
}

func (t *Table[T]) Iter() *TableIter[T] {
	return &TableIter[T]{
		Table: t,
		Iter: &Iter{
			Size: &t.Size,
		},
	}
}

type TableIter[T any] struct {
	Table *Table[T]
	Iter  *Iter
}

func (ti *TableIter[T]) Start() (*TableIter[T], T, bool) {
	cell, done := ti.Next()
	return ti, cell, done
}

func (ti *TableIter[T]) Next() (T, bool) {
	i, done := ti.Iter.Next()
	return ti.Table.Data[i], done
}

func (ti *TableIter[T]) Write(cell T) bool {
	i, done := ti.Iter.Next()
	if !done {
		ti.Table.Add(i.Row, i.Col, cell)
	}
	return done
}
