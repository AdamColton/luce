package table

// Index represents a position in a table.
type Index struct {
	Row, Col int
}

// UpdateSize assumes i is a size and i2 is a position. If i2 would imply a
// larger size, i is updated.
func (i Index) UpdateSize(i2 Index) Index {
	if i2.Col >= i.Col {
		i.Col = i2.Col + 1
	}
	if i2.Row >= i.Row {
		i.Row = i2.Row + 1
	}
	return i
}

// Iter loops over all the positions inside Size.
type Iter struct {
	Cur, Size *Index
}

// Next Iter value.
func (ii *Iter) Next() (idx Index, done bool) {
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

// Idx of the Iter.
func (ii *Iter) Idx() int {
	return ii.Cur.Col + ii.Cur.Row*ii.Size.Col
}

// Table of data.
type Table[T any] struct {
	Data   map[Index]T
	Size   Index
	Labels struct {
		Rows, Cols []string
	}
}

// New creates a Table
func New[T any]() *Table[T] {
	return &Table[T]{
		Data: make(map[Index]T),
	}
}

// Add (or overwrite) a cell value.
func (t *Table[T]) Add(r, c int, cell T) {
	i := Index{r, c}
	t.Data[i] = cell
	t.Size = t.Size.UpdateSize(i)
}

// Iter for iterating over the table.
func (t *Table[T]) Iter() *TableIter[T] {
	return &TableIter[T]{
		Table: t,
		Iter: &Iter{
			Size: &t.Size,
		},
	}
}

// TableIter is used to iterate over a table.
type TableIter[T any] struct {
	Table *Table[T]
	Iter  *Iter
}

// Start the iterator.
func (ti *TableIter[T]) Start() (*TableIter[T], T, bool) {
	ti.Iter.Cur = nil
	cell, done := ti.Next()
	return ti, cell, done
}

// Next cell in the iterator.
func (ti *TableIter[T]) Next() (T, bool) {
	i, done := ti.Iter.Next()
	return ti.Table.Data[i], done
}

// Write a value to the cell and move to the next iterator position.
func (ti *TableIter[T]) Write(cell T) bool {
	i, done := ti.Iter.Next()
	if !done {
		ti.Table.Add(i.Row, i.Col, cell)
	}
	return done
}
