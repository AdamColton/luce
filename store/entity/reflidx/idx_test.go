package reflidx_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/wrap/json"
	"github.com/adamcolton/luce/store/entity/reflidx"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/stretchr/testify/assert"
)

type foo struct {
	key string
	S   string `entityIndexSingle:""`
	F   float64
	I   int
}

func (f *foo) EntKey() []byte {
	return []byte(f.key)
}

func (f *foo) SetEntKey(key []byte) {
	f.key = string(key)
}

func (f *foo) Mod2() []byte {
	return []byte{byte(f.I % 2)}
}

var r = reflidx.Round{
	Serializer:   json.NewSerializer("", ""),
	Deserializer: json.Deserializer{},
}

func TestIdx(t *testing.T) {
	fac := ephemeral.Factory(bytebtree.New, 1)
	es := reflidx.New[*foo](r, fac)
	es.AddMethod("Mod2", true)

	f := &foo{
		key: "key1",
		S:   "just a string",
		F:   3.1415,
		I:   42,
	}
	_, err := es.Put(f, nil)
	assert.NoError(t, err)

	notFound := lerr.Str("not found")
	s := lerr.OK(es.Index("S"))(notFound)
	m2 := lerr.OK(es.Index("Mod2"))(notFound)

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

	ents := lerr.Must(s.LookupEnts([]byte(f.S)))
	assert.Len(t, ents, 1)
	assert.Equal(t, f, ents[0])

	ents = lerr.Must(m2.LookupEnts([]byte{0}))
	assert.Len(t, ents, 1)
	assert.Equal(t, f, ents[0])

}
