package prefix

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestRemove(t *testing.T) {
	p := New()
	p.Upsert("abc")
	p.Upsert("ab")
	p.Upsert("a")

	p.Upsert("defga")
	p.Upsert("ef")

	p.Remove("ab")
	assert.True(t, p.Find("abc").IsWord())
	assert.False(t, p.Find("ab").IsWord())
	p.Remove("abc")
	assert.Nil(t, p.Find("abc"))
	assert.Nil(t, p.Find("ab"))

	p.Remove("ef")
	p.Remove("xa")

	before := len(p.starts)
	got := p.Containing("a").Strings().ToSlice(nil)
	slice.LT[string]().Sort(got)
	assert.Equal(t, []string{"a", "defga"}, got)
	p.Purge()
	assert.True(t, len(p.starts) < before)
}
