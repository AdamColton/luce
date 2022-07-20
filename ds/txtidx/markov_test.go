package txtidx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChildCount(t *testing.T) {
	m := newMarkov()
	m.upsert("abc")
	_, n := m.find("a")
	i := findChild('b', n.children)
	assert.Equal(t, uint16(1), n.children[i].count)

	m.upsert("abd")
	_, n = m.find("a")
	i = findChild('b', n.children)
	assert.Equal(t, uint16(2), n.children[i].count)

	m.upsert("abe")
	_, n = m.find("a")
	i = findChild('b', n.children)
	assert.Equal(t, uint16(3), n.children[i].count)
	_, n = m.find("ab")
	i = findChild('e', n.children)
	assert.Equal(t, uint16(1), n.children[i].count)

	m.upsert("abel")
	_, n = m.find("a")
	i = findChild('b', n.children)
	assert.Equal(t, uint16(4), n.children[i].count)
	_, n = m.find("ab")
	i = findChild('e', n.children)
	assert.Equal(t, uint16(2), n.children[i].count)

	m.upsert("abc")
	_, n = m.find("a")
	i = findChild('b', n.children)
	assert.Equal(t, uint16(4), n.children[i].count)

	m.deleteWord("abcx")
	_, n = m.find("a")
	i = findChild('b', n.children)
	assert.Equal(t, uint16(4), n.children[i].count)

	m.deleteWord("abc")
	_, n = m.find("a")
	i = findChild('b', n.children)
	assert.Equal(t, uint16(3), n.children[i].count)
}

func TestSuggest(t *testing.T) {
	m := newMarkov()
	m.upsert("abcd")
	m.upsert("abcde")
	m.upsert("abcf")
	m.upsert("aghi")
	m.upsert("aghij")
	m.upsert("aklm")
	m.upsert("anop")

	s := m.suggest("a", -1)
	expected := []Suggestion{
		{Word: "bcde", Terminals: []int{4, 3}},
		{Word: "ghij", Terminals: []int{4, 3}},
		{Word: "klm", Terminals: []int{3}},
		{Word: "nop", Terminals: []int{3}},
	}
	assert.Equal(t, expected, s)
}
