package morph_test

import (
	"strconv"
	"testing"

	"github.com/adamcolton/luce/ds/morph"
	"github.com/stretchr/testify/assert"
)

func TestKeyValConversions(t *testing.T) {
	all := morph.NewKeyValAll(func(k string, v int) string { return k + strconv.Itoa(v) })
	out, include := all.ToKV()("a", 1)
	assert.Equal(t, "a1", out)
	assert.True(t, include)

	kv := morph.NewKeyVal(func(k string, v int) (int, bool) { return v * 2, v > 1 })
	out2, include := kv("a", 1)
	assert.Equal(t, 2, out2)
	assert.False(t, include)
	out2, include = kv("a", 2)
	assert.Equal(t, 4, out2)
	assert.True(t, include)
}

func TestValConversions(t *testing.T) {
	all := morph.NewValAll(func(v int) string { return strconv.Itoa(v) })
	out, include := all.ToV()(7)
	assert.Equal(t, "7", out)
	assert.True(t, include)

	v := morph.NewVal(func(v int) (int, bool) { return v * 2, v%2 == 0 })
	_, include = v(1)
	assert.False(t, include)
	out2, include := v(2)
	assert.Equal(t, 4, out2)
	assert.True(t, include)
}

func TestNull(t *testing.T) {
	n := morph.NewNull(func() int { return 7 })
	assert.Equal(t, 7, n())
}
