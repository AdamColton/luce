package prefix_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/prefix"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

var lt = slice.LT[string]()

func TestBasic(t *testing.T) {
	p := prefix.New()

	n, insert := p.Upsert("")
	assert.Nil(t, n)
	assert.False(t, insert)

	n, insert = p.Upsert("abc")
	assert.True(t, n.IsWord())
	assert.True(t, insert)
	assert.Equal(t, "abc", n.Gram())
	assert.Equal(t, 0, n.ChildrenCount())
	_, insert = p.Upsert("abc")
	assert.False(t, insert)

	n = p.Find("a")
	assert.False(t, n.IsWord())
	assert.Equal(t, "a", n.Gram())
	assert.Equal(t, 1, n.ChildrenCount())
	assert.Equal(t, []rune{'b'}, n.Children())

	n = n.Child('b')
	assert.False(t, n.IsWord())
	assert.Equal(t, "ab", n.Gram())
	assert.Equal(t, 1, n.ChildrenCount())
	assert.Equal(t, []rune{'c'}, n.Children())

	n = p.Find("aaaa")
	assert.Nil(t, n)
}

func TestAllWords(t *testing.T) {
	words := []string{"abcd", "abce", "abcf", "bbcd"}
	p := prefix.New()
	for _, w := range words {
		p.Upsert(w)
	}

	ls := p.Find("").AllWords().Strings()
	got := ls.ToSlice(nil)
	lt.Sort(got)
	assert.Equal(t, words, got)
}
