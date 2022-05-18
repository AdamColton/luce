package linkedlist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedList(t *testing.T) {
	a := assert.New(t)

	n := Simple[int]()
	n.Set(123)
	a.Equal(123, n.Get())

	cp := n.Copy()
	cp.Prepend(3, 1, 4, 1)
	a.Equal(1, cp.Get())
	a.Equal(123, n.Get())

	cp.Next()
	a.Equal(4, cp.Get())

	cp.Next()
	a.Equal(1, cp.Get())

	n.Head()
	a.Equal(1, n.Get())

	cp.Tail()
	a.Equal(123, cp.Get())
	cp.Set(456)
	a.Equal(456, cp.Get())
	cp.Append()
	a.Equal(456, cp.Get())
	cp.Append(7, 8, 9)
	a.Equal(9, cp.Get())

	n.Next()
	n.Append(10, 11, 12)
	n.Prepend(13, 14)
	n.Prepend()
}

func TestLoop(t *testing.T) {
	n := Simple[int]()
	n.Prepend(1, 2, 3, 4, 5, 6)
	expected := []int{1, 2, 3, 4, 5, 6}
	i := 0
	for n.Tail(); !n.Nil(); n.Prev() {
		assert.Equal(t, expected[i], n.Get())
		i++
	}

	expected = []int{6, 5, 4, 3, 2, 1}
	i = 0
	for n.Head(); !n.Nil(); n.Next() {
		assert.Equal(t, expected[i], n.Get())
		i++
	}
}
