package table_test

import (
	"testing"

	"github.com/adamcolton/luce/util/table"
	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	tab := table.New[int]()
	tab.Size = table.Index{3, 3}
	i := tab.Iter()
	data := []int{3, 1, 4, 1, 5, 9, 2, 6, 5}
	for _, c := range data {
		i.Write(c)
	}

	count := 0
	i = tab.Iter()
	for _, c, done := i.Start(); !done; c, done = i.Next() {
		assert.Equal(t, c, data[i.Iter.Idx()])
		count++
	}
	idx, done := i.Iter.Next()
	assert.Equal(t, table.Index{-1, -1}, idx)
	assert.True(t, done)
	assert.Len(t, data, count)

	tab2 := table.New[int]()
	for i, v := range tab.Data {
		tab2.Add(i.Row, i.Col, v)
	}
	assert.Equal(t, tab.Size, tab2.Size)
}
