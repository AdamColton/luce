package prefix_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/prefix"
	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/entity/enttest"
	"github.com/stretchr/testify/assert"
)

func TestPrefixEntity(t *testing.T) {
	enttest.Setup()
	p := prefix.New()

	words := []string{"test", "testing", "tea", "adam"}

	for _, w := range words {
		p.Upsert(w)
	}

	n := p.Find("te")
	assert.NotNil(t, n)

	ref, err := p.Save()
	assert.NoError(t, err)
	k := ref.EntKey()
	ref.Clear(true)
	entity.ClearCache()

	r2 := entity.KeyRef[prefix.Prefix](k)
	p = r2.GetPtr()
	if !assert.NotNil(t, p) {
		return
	}
	n = p.Find("te")
	assert.NotNil(t, n)
	s := n.Suggest(2)
	assert.Len(t, s, 2)
	expected := []string{"sting", "st"}
	assert.Equal(t, expected, s[0].Words("").ToSlice(nil))

	ls := p.Find("").AllWords().Strings().ToSlice(nil)
	lt.Sort(ls)
	lt.Sort(words)
}
