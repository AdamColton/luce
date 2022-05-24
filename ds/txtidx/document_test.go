package txtidx

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentString(t *testing.T) {
	expected := "--- This is a Test "

	c := NewCorpus()

	d := c.AddDoc(expected)

	assert.Equal(t, expected, d.String(c))

	// 4 words but only 2 variants
	// "is " and "a " have the same variant
	// "This " and "Test " have the same variant
	assert.Len(t, c.Variants, 2)

	assert.True(t, c.Find("test").Has(d.DocIDX))
}

func TestMultiDoc(t *testing.T) {
	c := NewCorpus()
	c.AddDoc("The sun was shining on the sea")
	c.AddDoc("Shining with all it's might")
	c.AddDoc("And it did the very best it could")
	c.AddDoc("To make the billows smooth and bright")
	c.AddDoc("And this was very odd because")
	c.AddDoc("It was the middle of the night")
	c.AddDoc("The moon was shining skulkily")
	c.AddDoc("Because she thought the sun")

	the := c.Find("the")
	assert.Equal(t, 6, the.Len())

	shining := c.Find("shining")
	assert.Equal(t, 3, shining.Len())

	both := c.Find("the", "shining").Slice(c)
	sort.Strings(both)
	expected := []string{
		"The moon was shining skulkily",
		"The sun was shining on the sea",
	}
	assert.Equal(t, expected, both)

	ill := c.Find("ill").Slice(c)
	expected = []string{
		"To make the billows smooth and bright",
	}
	assert.Equal(t, expected, ill)
}

func TestMarkov(t *testing.T) {
	m := NewMarkov()
	w := m.Upsert("test")
	w.Documents.add(123)
	w2 := m.Upsert("test")

	assert.True(t, w2.Documents.Has(123))

	w3 := m.Find("test")
	assert.True(t, w3.Documents.Has(123))
}

func TestDelete(t *testing.T) {
	c := NewCorpus()
	d0 := c.AddDoc("this is document 0 keyphrase")
	assert.Equal(t, DocIDX(0), d0.DocIDX)

	d1 := c.AddDoc("this is document 1")
	assert.Equal(t, DocIDX(1), d1.DocIDX)

	w := c.Roots.Find("keyphrase")
	assert.True(t, w.Documents.Has(d0.DocIDX))
	wid := w.WordIDX

	c.Delete(d0.DocIDX)
	assert.Nil(t, c.Docs[0])
	assert.Nil(t, c.Words[wid])

	w = c.Roots.Find("keyphrase")
	assert.Nil(t, w)

	d2 := c.AddDoc("this is document 2, reallocated DocIDX 0 keyphrase")
	assert.Equal(t, DocIDX(0), d2.DocIDX)
	assert.NotNil(t, c.Words[wid])

	w = c.Roots.Find("keyphrase")
	assert.True(t, w.Documents.Has(d2.DocIDX))
}
