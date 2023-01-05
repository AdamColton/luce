package prefix

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	p := New()
	n, insert := p.Upsert("abc")
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
}

func TestContains(t *testing.T) {
	p := New()
	p.Upsert("afoo")
	p.Upsert("bfoo")
	p.Upsert("ccfoo")
	p.Upsert("ddfoodd")
	p.Upsert("fooe")
	p.Upsert("fbar")
	p.Upsert("barg")

	foos := p.Contains("foo")
	assert.Len(t, foos, 5)
	var got []string
	for _, foo := range foos {
		got = append(got, foo.Gram())
	}
	sort.Strings(got)
	expected := []string{"afoo", "bfoo", "ccfoo", "ddfoo", "foo"}
	assert.Equal(t, expected, got)
}

func TestSuggest(t *testing.T) {
	p := New()
	p.Upsert("abcd")
	p.Upsert("abcde")
	p.Upsert("abce")
	p.Upsert("ace")
	p.Upsert("adgf")
	p.Upsert("adg")

	n := p.Find("a")
	s := n.Suggest(2)
	assert.Equal(t, "bcde", s[0].Word)
	assert.Equal(t, "dgf", s[1].Word)
}
