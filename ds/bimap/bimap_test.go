package bimap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/bimap"
	"github.com/stretchr/testify/assert"
)

func TestBimap(t *testing.T) {
	bi := bimap.New[string, int](20)
	d := bi.Add("123", 123)
	assert.False(t, d.A.Deleted)
	assert.False(t, d.B.Deleted)

	a, found := bi.B(123)
	assert.True(t, found)
	assert.Equal(t, "123", a)

	b, found := bi.A("123")
	assert.True(t, found)
	assert.Equal(t, 123, b)

	d = bi.Add("123", 456)
	assert.False(t, d.A.Deleted)
	assert.True(t, d.B.Deleted)
	assert.Equal(t, 123, d.B.Value)

	a, found = bi.B(456)
	assert.True(t, found)
	assert.Equal(t, "123", a)

	_, found = bi.B(123)
	assert.False(t, found)

	_, found = bi.A("456")
	assert.False(t, found)

	d = bi.Add("456", 456)
	assert.True(t, d.A.Deleted)
	assert.Equal(t, "123", d.A.Value)
	assert.False(t, d.B.Deleted)
}
