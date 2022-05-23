package txtidx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentString(t *testing.T) {
	expected := "--- This is a test."

	c := NewCorpus()
	c.Max.DocID = 12

	pp := c.newPP()
	pp.set(expected)
	d := pp.build()

	assert.Equal(t, expected, d.String(c))

	_, hasDoc := c.Roots["is"].Documents[d.DocID]
	assert.True(t, hasDoc)
}
