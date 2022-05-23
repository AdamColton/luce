package txtidx

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentString(t *testing.T) {
	expected := "--- This is a Test "

	c := NewCorpus()
	c.Max.DocID = 12

	d := c.AddDoc(expected)

	assert.Equal(t, expected, d.String(c))

	// 4 words but only 2 variants
	// "is " and "a " have the same variant
	// "This " and "Test " have the same variant
	assert.Len(t, c.Variants, 2)

	assert.True(t, c.Find("test").Has(d.DocID))
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
}
