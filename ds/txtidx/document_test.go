package txtidx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentString(t *testing.T) {
	expected := "--- This is a Test "

	c := NewCorpus()
	c.Max.DocID = 12

	pp := c.newPP()
	pp.set(expected)
	d := pp.build()

	assert.Equal(t, expected, d.String(c))

	// 4 words but only 2 variants
	// "is " and "a " have the same variant
	// "This " and "Test " have the same variant
	assert.Len(t, c.Variants, 2)

	assert.True(t, c.Find("test").Has(d.DocID))
}
