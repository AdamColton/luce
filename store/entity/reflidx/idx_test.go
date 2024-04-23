package reflidx_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/serial/wrap/json"
	"github.com/adamcolton/luce/store/entity/reflidx"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/stretchr/testify/assert"
)

type foo struct {
	key string
	S   string
	F   float64
	I   int
}

func (f *foo) EntKey() []byte {
	return []byte(f.key)
}
func (f *foo) SetEntKey(key []byte) {
	f.key = string(key)
}

var r = reflidx.Round{
	Serializer:   json.NewSerializer("", ""),
	Deserializer: json.Deserializer{},
}

func TestIdx(t *testing.T) {
	fac := ephemeral.Factory(bytebtree.New, 1)
	es := reflidx.New[*foo](r, fac)

	f := &foo{
		key: "key1",
		S:   "just a string",
		F:   3.1415,
		I:   42,
	}
	_, err := es.Put(f, nil)
	assert.NoError(t, err)

	f2 := &foo{
		key: "key1",
	}
	err = es.Load(f2)
	assert.NoError(t, err)
	assert.Equal(t, f, f2)
	found, f3, err := es.Get(f.EntKey())
	assert.True(t, found)
	assert.NoError(t, err)
	assert.Equal(t, f, f3)
}
