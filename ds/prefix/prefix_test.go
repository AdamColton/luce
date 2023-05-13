package prefix_test

import (
	"sort"
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
	words := []string{"abcdc", "abce", "abcf", "bbcd"}
	p := prefix.New()
	for _, w := range words {
		p.Upsert(w)
	}

	ls := p.Find("").AllWords().Strings()
	got := ls.Slice(nil)
	lt.Sort(got)
	assert.Equal(t, words, got)

	got = p.Containing("c").AllWords().Strings().Slice(nil)
	lt.Sort(got)
	assert.Equal(t, words, got)
}

func TestContaining(t *testing.T) {
	p := prefix.New()
	words := []string{"afoo", "bfoo", "ccfoo", "ddfoodd", "fooe", "fbar", "barg"}
	for _, w := range words {
		p.Upsert(w)
	}

	got := p.Containing("foo").Strings().Slice(nil)
	sort.Strings(got)

	expected := []string{"afoo", "bfoo", "ccfoo", "ddfoo", "foo"}
	assert.Equal(t, expected, got)

	assert.Nil(t, p.Containing(""))

	p.Upsert("test")
	got = p.Containing("t").Strings().Slice(nil)
	lt.Sort(got)
	expected = []string{"t", "test"}
	assert.Equal(t, expected, got)
}

func TestSuggest(t *testing.T) {
	p := prefix.New()
	p.Upsert("abcd")
	p.Upsert("abcde")
	p.Upsert("abce")
	p.Upsert("ace")
	p.Upsert("adgf")
	p.Upsert("adg")
	p.Upsert("abc")

	n := p.Find("a")
	s := n.Suggest(2)
	assert.Len(t, s, 2)

	expected := slice.New([]string{"abcde", "abcd", "abc"})
	assert.Equal(t, expected, s[0].Words("a"))

	expected = slice.New([]string{"adgf", "adg"})
	assert.Equal(t, expected, s[1].Words("a"))
}
